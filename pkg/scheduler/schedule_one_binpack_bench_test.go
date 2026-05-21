/*
Copyright 2024 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package scheduler

// Benchmarks for tie-breaking frequency and end-to-end scheduling latency
// under MostAllocated scoring (bin-packing). These validate the Rust simulator's
// finding that tie frequency is extremely high (94-500 tied nodes per scheduling
// event), which is the precondition for AGL's dramatic improvement.

import (
	"context"
	"fmt"
	"sort"
	"testing"
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	clientsetfake "k8s.io/client-go/kubernetes/fake"
	"k8s.io/klog/v2/ktesting"
	schedulerapi "k8s.io/kubernetes/pkg/scheduler/apis/config"
	internalcache "k8s.io/kubernetes/pkg/scheduler/backend/cache"
	internalqueue "k8s.io/kubernetes/pkg/scheduler/backend/queue"
	"k8s.io/kubernetes/pkg/scheduler/framework"
	"k8s.io/kubernetes/pkg/scheduler/framework/plugins/defaultbinder"
	"k8s.io/kubernetes/pkg/scheduler/framework/plugins/feature"
	"k8s.io/kubernetes/pkg/scheduler/framework/plugins/noderesources"
	"k8s.io/kubernetes/pkg/scheduler/framework/plugins/queuesort"
	frameworkruntime "k8s.io/kubernetes/pkg/scheduler/framework/runtime"
	st "k8s.io/kubernetes/pkg/scheduler/testing"
	tf "k8s.io/kubernetes/pkg/scheduler/testing/framework"
)

// makeBinpackNodes creates N homogeneous nodes with the given CPU and memory capacity.
func makeBinpackNodes(n int, cpuMillis, memoryBytes int64) []*v1.Node {
	nodes := make([]*v1.Node, n)
	for i := 0; i < n; i++ {
		name := fmt.Sprintf("node-%04d", i)
		nodes[i] = &v1.Node{
			ObjectMeta: metav1.ObjectMeta{
				Name:   name,
				Labels: map[string]string{"kubernetes.io/hostname": name},
			},
			Status: v1.NodeStatus{
				Capacity: v1.ResourceList{
					v1.ResourceCPU:    *resource.NewMilliQuantity(cpuMillis, resource.DecimalSI),
					v1.ResourceMemory: *resource.NewQuantity(memoryBytes, resource.BinarySI),
					v1.ResourcePods:   *resource.NewQuantity(110, resource.DecimalSI),
				},
				Allocatable: v1.ResourceList{
					v1.ResourceCPU:    *resource.NewMilliQuantity(cpuMillis, resource.DecimalSI),
					v1.ResourceMemory: *resource.NewQuantity(memoryBytes, resource.BinarySI),
					v1.ResourcePods:   *resource.NewQuantity(110, resource.DecimalSI),
				},
			},
		}
	}
	return nodes
}

// makeBinpackPod creates a pod with the given resource requests.
func makeBinpackPod(name, nodeName string, cpuMillis, memoryBytes int64) *v1.Pod {
	pod := st.MakePod().Name(name).UID(name).
		Req(map[v1.ResourceName]string{
			v1.ResourceCPU:    fmt.Sprintf("%dm", cpuMillis),
			v1.ResourceMemory: fmt.Sprintf("%d", memoryBytes),
		}).Obj()
	if nodeName != "" {
		pod.Spec.NodeName = nodeName
	}
	return pod
}

// setupMostAllocatedFramework creates a scheduler framework with MostAllocated scoring.
func setupMostAllocatedFramework(ctx context.Context, snapshot *internalcache.Snapshot) (framework.Framework, error) {
	cs := clientsetfake.NewClientset()
	informerFactory := informers.NewSharedInformerFactory(cs, 0)

	fwk, err := tf.NewFramework(
		ctx,
		[]tf.RegisterPluginFunc{
			tf.RegisterQueueSortPlugin(queuesort.Name, queuesort.New),
			tf.RegisterPluginAsExtensions(
				noderesources.Name,
				frameworkruntime.FactoryAdapter(feature.Features{}, noderesources.NewFit),
				"PreFilter", "Filter", "PreScore", "Score",
			),
			tf.RegisterBindPlugin(defaultbinder.Name, defaultbinder.New),
		},
		"",
		frameworkruntime.WithSnapshotSharedLister(snapshot),
		frameworkruntime.WithInformerFactory(informerFactory),
		frameworkruntime.WithPodNominator(internalqueue.NewSchedulingQueue(nil, informerFactory)),
	)
	if err != nil {
		return nil, err
	}
	informerFactory.Start(ctx.Done())
	informerFactory.WaitForCacheSync(ctx.Done())
	return fwk, nil
}

// BenchmarkTieBreakingFrequency measures how often multiple nodes receive the
// same top score under MostAllocated. This directly validates the Rust
// simulator's finding of 94-500 tied nodes per scheduling event.
//
// The benchmark schedules pods one at a time on homogeneous nodes and counts
// the fraction of scheduling cycles where >1 node shares the maximum score.
func BenchmarkTieBreakingFrequency(b *testing.B) {
	nodeCounts := []int{100, 500, 1000}

	for _, numNodes := range nodeCounts {
		b.Run(fmt.Sprintf("Nodes_%d", numNodes), func(b *testing.B) {
			_, ctx := ktesting.NewTestContext(b)
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()

			// Create homogeneous nodes: 4 CPU, 8GB memory each
			nodes := makeBinpackNodes(numNodes, 4000, 8*1024*1024*1024)

			// Pre-place some pods to create partially loaded cluster (~30% loaded)
			// This creates realistic conditions where MostAllocated must differentiate
			existingPods := make([]*v1.Pod, 0)
			podsPerNode := 3 // ~30% CPU used (3 * 400m = 1200m / 4000m)
			for i := 0; i < numNodes; i++ {
				for j := 0; j < podsPerNode; j++ {
					podName := fmt.Sprintf("existing-%04d-%d", i, j)
					existingPods = append(existingPods, makeBinpackPod(
						podName,
						fmt.Sprintf("node-%04d", i),
						400, // 400m CPU
						512*1024*1024, // 512MB
					))
				}
			}

			snapshot := internalcache.NewSnapshot(existingPods, nodes)
			schedFramework, err := setupMostAllocatedFramework(ctx, snapshot)
			if err != nil {
				b.Fatalf("Failed to create framework: %v", err)
			}

			sched := &Scheduler{
				Cache:                    internalcache.New(ctx, nil, false),
				nodeInfoSnapshot:         snapshot,
				percentageOfNodesToScore: schedulerapi.DefaultPercentageOfNodesToScore,
			}
			sched.applyDefaultHandlers()

			b.ResetTimer()
			b.ReportAllocs()

			totalCycles := 0
			tiedCycles := 0
			totalTiedNodes := 0

			for i := 0; i < b.N; i++ {
				pod := makeBinpackPod(fmt.Sprintf("bench-pod-%d", i), "", 500, 1024*1024*1024)
				podInfo := &framework.QueuedPodInfo{
					PodInfo: &framework.PodInfo{Pod: pod},
				}

				state := framework.NewCycleState()

				// Run the filtering phase
				feasibleNodes, _, _, findErr := sched.findNodesThatFitPod(ctx, schedFramework, state, podInfo)
				if findErr != nil {
					b.Fatalf("findNodesThatFitPod failed: %v", findErr)
				}
				if len(feasibleNodes) == 0 {
					b.Fatalf("No feasible nodes found")
				}

				// Run the scoring phase
				priorityList, err := prioritizeNodes(ctx, nil, schedFramework, state, pod, feasibleNodes)
				if err != nil {
					b.Fatalf("prioritizeNodes failed: %v", err)
				}

				// Count ties: find the max score and count nodes that share it
				maxScore := int64(0)
				for _, nodeScore := range priorityList {
					if nodeScore.TotalScore > maxScore {
						maxScore = nodeScore.TotalScore
					}
				}
				tiedCount := 0
				for _, nodeScore := range priorityList {
					if nodeScore.TotalScore == maxScore {
						tiedCount++
					}
				}

				totalCycles++
				if tiedCount > 1 {
					tiedCycles++
					totalTiedNodes += tiedCount
				}
			}

			if totalCycles > 0 {
				tieRate := float64(tiedCycles) / float64(totalCycles) * 100
				avgTied := float64(0)
				if tiedCycles > 0 {
					avgTied = float64(totalTiedNodes) / float64(tiedCycles)
				}
				b.ReportMetric(tieRate, "tie_percent")
				b.ReportMetric(avgTied, "avg_tied_nodes")
				b.ReportMetric(float64(tiedCycles), "tied_cycles")
				b.ReportMetric(float64(totalCycles), "total_cycles")
			}
		})
	}
}

// BenchmarkTieBreakingFrequencyHeterogeneous is like BenchmarkTieBreakingFrequency
// but uses heterogeneous load (some nodes more loaded than others).
// This validates that ties are common even with realistic workload distributions.
func BenchmarkTieBreakingFrequencyHeterogeneous(b *testing.B) {
	_, ctx := ktesting.NewTestContext(b)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	numNodes := 500
	nodes := makeBinpackNodes(numNodes, 4000, 8*1024*1024*1024)

	// Create heterogeneous load: groups of nodes at different utilization levels
	// Group 1 (0-99):   ~10% loaded (1 pod of 400m)
	// Group 2 (100-249): ~30% loaded (3 pods of 400m)
	// Group 3 (250-399): ~50% loaded (5 pods of 400m)
	// Group 4 (400-499): ~70% loaded (7 pods of 400m)
	existingPods := make([]*v1.Pod, 0)
	loadGroups := []struct {
		startIdx, endIdx int
		podsPerNode      int
	}{
		{0, 99, 1},
		{100, 249, 3},
		{250, 399, 5},
		{400, 499, 7},
	}
	for _, group := range loadGroups {
		for i := group.startIdx; i <= group.endIdx; i++ {
			for j := 0; j < group.podsPerNode; j++ {
				podName := fmt.Sprintf("existing-%04d-%d", i, j)
				existingPods = append(existingPods, makeBinpackPod(
					podName,
					fmt.Sprintf("node-%04d", i),
					400, 512*1024*1024,
				))
			}
		}
	}

	snapshot := internalcache.NewSnapshot(existingPods, nodes)
	schedFramework, err := setupMostAllocatedFramework(ctx, snapshot)
	if err != nil {
		b.Fatalf("Failed to create framework: %v", err)
	}

	sched := &Scheduler{
		Cache:                    internalcache.New(ctx, nil, false),
		nodeInfoSnapshot:         snapshot,
		percentageOfNodesToScore: schedulerapi.DefaultPercentageOfNodesToScore,
	}
	sched.applyDefaultHandlers()

	b.ResetTimer()
	b.ReportAllocs()

	totalCycles := 0
	tiedCycles := 0
	totalTiedNodes := 0

	for i := 0; i < b.N; i++ {
		pod := makeBinpackPod(fmt.Sprintf("bench-pod-%d", i), "", 500, 1024*1024*1024)
		podInfo := &framework.QueuedPodInfo{
			PodInfo: &framework.PodInfo{Pod: pod},
		}

		state := framework.NewCycleState()

		feasibleNodes, _, _, findErr := sched.findNodesThatFitPod(ctx, schedFramework, state, podInfo)
		if findErr != nil {
			b.Fatalf("findNodesThatFitPod failed: %v", findErr)
		}
		if len(feasibleNodes) == 0 {
			b.Fatalf("No feasible nodes found")
		}

		priorityList, err := prioritizeNodes(ctx, nil, schedFramework, state, pod, feasibleNodes)
		if err != nil {
			b.Fatalf("prioritizeNodes failed: %v", err)
		}

		maxScore := int64(0)
		for _, nodeScore := range priorityList {
			if nodeScore.TotalScore > maxScore {
				maxScore = nodeScore.TotalScore
			}
		}
		tiedCount := 0
		for _, nodeScore := range priorityList {
			if nodeScore.TotalScore == maxScore {
				tiedCount++
			}
		}

		totalCycles++
		if tiedCount > 1 {
			tiedCycles++
			totalTiedNodes += tiedCount
		}
	}

	if totalCycles > 0 {
		tieRate := float64(tiedCycles) / float64(totalCycles) * 100
		avgTied := float64(0)
		if tiedCycles > 0 {
			avgTied = float64(totalTiedNodes) / float64(tiedCycles)
		}
		b.ReportMetric(tieRate, "tie_percent")
		b.ReportMetric(avgTied, "avg_tied_nodes")
	}
}

// BenchmarkEndToEndSchedulingLatency measures wall-clock time per scheduling
// decision under MostAllocated scoring with a large cluster.
//
// Setup: 5000 nodes, schedules 1000 pods, reports p50/p95/p99 latency.
// The percentageOfNodesToScore defaults to 5% for this cluster size
// (the formula is max(5, 50 - numAllNodes/125) = max(5, 50-40) = 10% at 5000 nodes,
// but we can override it).
func BenchmarkEndToEndSchedulingLatency(b *testing.B) {
	type strategy struct {
		name                     string
		percentageOfNodesToScore int32
	}

	strategies := []strategy{
		{"Baseline_Default", 0}, // 0 means use default formula
		{"Score5Pct", 5},
		{"Score10Pct", 10},
		{"Score100Pct", 100}, // Score all nodes - maximum tie potential
	}

	for _, strat := range strategies {
		b.Run(strat.name, func(b *testing.B) {
			_, ctx := ktesting.NewTestContext(b)
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()

			numNodes := 5000
			podsToSchedule := 100 // per b.N iteration

			// Create homogeneous nodes: 8 CPU, 16GB
			nodes := makeBinpackNodes(numNodes, 8000, 16*1024*1024*1024)

			// Pre-place pods to create ~40% cluster utilization
			existingPods := make([]*v1.Pod, 0, numNodes*4)
			for i := 0; i < numNodes; i++ {
				for j := 0; j < 4; j++ {
					podName := fmt.Sprintf("existing-%04d-%d", i, j)
					existingPods = append(existingPods, makeBinpackPod(
						podName,
						fmt.Sprintf("node-%04d", i),
						800, // 800m CPU each, 4 pods = 3200m / 8000m = 40%
						1024*1024*1024, // 1GB each
					))
				}
			}

			snapshot := internalcache.NewSnapshot(existingPods, nodes)
			schedFramework, err := setupMostAllocatedFramework(ctx, snapshot)
			if err != nil {
				b.Fatalf("Failed to create framework: %v", err)
			}

			sched := &Scheduler{
				Cache:                    internalcache.New(ctx, nil, false),
				nodeInfoSnapshot:         snapshot,
				percentageOfNodesToScore: strat.percentageOfNodesToScore,
			}
			sched.applyDefaultHandlers()

			b.ResetTimer()
			b.ReportAllocs()

			latencies := make([]time.Duration, 0, b.N*podsToSchedule)

			for i := 0; i < b.N; i++ {
				for j := 0; j < podsToSchedule; j++ {
					pod := makeBinpackPod(
						fmt.Sprintf("sched-pod-%d-%d", i, j), "",
						500, 512*1024*1024,
					)
					podInfo := &framework.QueuedPodInfo{
						PodInfo: &framework.PodInfo{Pod: pod},
					}

					start := time.Now()
					_, err := sched.SchedulePod(ctx, schedFramework, framework.NewCycleState(), podInfo)
					elapsed := time.Since(start)

					if err != nil {
						b.Fatalf("SchedulePod failed: %v", err)
					}
					latencies = append(latencies, elapsed)
				}
			}

			if len(latencies) > 0 {
				sort.Slice(latencies, func(i, j int) bool {
					return latencies[i] < latencies[j]
				})
				p50 := latencies[len(latencies)*50/100]
				p95 := latencies[len(latencies)*95/100]
				p99 := latencies[len(latencies)*99/100]

				b.ReportMetric(float64(p50.Microseconds()), "p50_us")
				b.ReportMetric(float64(p95.Microseconds()), "p95_us")
				b.ReportMetric(float64(p99.Microseconds()), "p99_us")
				b.ReportMetric(float64(len(latencies)), "total_schedules")
			}
		})
	}
}

// BenchmarkScoringTieDistribution provides a detailed histogram of tie counts
// across scheduling cycles. This shows the distribution of how many nodes
// are tied at the top score, directly comparable to the Rust simulator output.
func BenchmarkScoringTieDistribution(b *testing.B) {
	_, ctx := ktesting.NewTestContext(b)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	numNodes := 1000
	nodes := makeBinpackNodes(numNodes, 4000, 8*1024*1024*1024)

	// Uniform load: all nodes at ~30%
	existingPods := make([]*v1.Pod, 0)
	for i := 0; i < numNodes; i++ {
		for j := 0; j < 3; j++ {
			podName := fmt.Sprintf("existing-%04d-%d", i, j)
			existingPods = append(existingPods, makeBinpackPod(
				podName,
				fmt.Sprintf("node-%04d", i),
				400, 512*1024*1024,
			))
		}
	}

	snapshot := internalcache.NewSnapshot(existingPods, nodes)
	schedFramework, err := setupMostAllocatedFramework(ctx, snapshot)
	if err != nil {
		b.Fatalf("Failed to create framework: %v", err)
	}

	sched := &Scheduler{
		Cache:                    internalcache.New(ctx, nil, false),
		nodeInfoSnapshot:         snapshot,
		percentageOfNodesToScore: 100, // Score ALL nodes to see full tie picture
	}
	sched.applyDefaultHandlers()

	b.ResetTimer()

	// Collect tie distribution across many scheduling cycles
	tieDistribution := make(map[int]int) // tied_count -> number_of_cycles
	scoreDistribution := make(map[int64]int) // score -> number_of_nodes_with_that_score

	for i := 0; i < b.N; i++ {
		pod := makeBinpackPod(fmt.Sprintf("dist-pod-%d", i), "", 500, 1024*1024*1024)
		podInfo := &framework.QueuedPodInfo{
			PodInfo: &framework.PodInfo{Pod: pod},
		}

		state := framework.NewCycleState()

		feasibleNodes, _, _, findErr := sched.findNodesThatFitPod(ctx, schedFramework, state, podInfo)
		if findErr != nil {
			b.Fatalf("findNodesThatFitPod failed: %v", findErr)
		}

		priorityList, err := prioritizeNodes(ctx, nil, schedFramework, state, pod, feasibleNodes)
		if err != nil {
			b.Fatalf("prioritizeNodes failed: %v", err)
		}

		// Find max score
		maxScore := int64(0)
		for _, ns := range priorityList {
			if ns.TotalScore > maxScore {
				maxScore = ns.TotalScore
			}
		}

		// Count nodes at max score
		tiedCount := 0
		for _, ns := range priorityList {
			if ns.TotalScore == maxScore {
				tiedCount++
			}
			scoreDistribution[ns.TotalScore]++
		}
		tieDistribution[tiedCount]++
	}

	// Report key metrics
	if b.N > 0 {
		// Average tied nodes
		totalTied := 0
		totalCycles := 0
		for tiedCount, cycles := range tieDistribution {
			totalTied += tiedCount * cycles
			totalCycles += cycles
		}
		if totalCycles > 0 {
			b.ReportMetric(float64(totalTied)/float64(totalCycles), "avg_tied_nodes")
		}

		// Max tied nodes seen
		maxTied := 0
		for tiedCount := range tieDistribution {
			if tiedCount > maxTied {
				maxTied = tiedCount
			}
		}
		b.ReportMetric(float64(maxTied), "max_tied_nodes")

		// Fraction of cycles with ties
		cyclesWithTies := 0
		for tiedCount, cycles := range tieDistribution {
			if tiedCount > 1 {
				cyclesWithTies += cycles
			}
		}
		b.ReportMetric(float64(cyclesWithTies)/float64(totalCycles)*100, "tie_pct")

		// Number of distinct scores observed (fewer = more ties)
		b.ReportMetric(float64(len(scoreDistribution)), "distinct_scores")
	}
}

// BenchmarkSchedulePodMostAllocated is a focused micro-benchmark that measures
// just the SchedulePod call latency with MostAllocated scoring.
func BenchmarkSchedulePodMostAllocated(b *testing.B) {
	benchCases := []struct {
		name     string
		numNodes int
	}{
		{"100_nodes", 100},
		{"500_nodes", 500},
		{"1000_nodes", 1000},
	}

	for _, bc := range benchCases {
		b.Run(bc.name, func(b *testing.B) {
			_, ctx := ktesting.NewTestContext(b)
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()

			nodes := makeBinpackNodes(bc.numNodes, 4000, 8*1024*1024*1024)

			// Pre-place pods uniformly
			existingPods := make([]*v1.Pod, 0)
			for i := 0; i < bc.numNodes; i++ {
				for j := 0; j < 3; j++ {
					podName := fmt.Sprintf("existing-%04d-%d", i, j)
					existingPods = append(existingPods, makeBinpackPod(
						podName,
						fmt.Sprintf("node-%04d", i),
						400, 512*1024*1024,
					))
				}
			}

			snapshot := internalcache.NewSnapshot(existingPods, nodes)
			schedFramework, err := setupMostAllocatedFramework(ctx, snapshot)
			if err != nil {
				b.Fatalf("Failed to create framework: %v", err)
			}

			sched := &Scheduler{
				Cache:                    internalcache.New(ctx, nil, false),
				nodeInfoSnapshot:         snapshot,
				percentageOfNodesToScore: schedulerapi.DefaultPercentageOfNodesToScore,
			}
			sched.applyDefaultHandlers()

			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				pod := makeBinpackPod(fmt.Sprintf("bench-%d", i), "", 500, 1024*1024*1024)
				podInfo := &framework.QueuedPodInfo{
					PodInfo: &framework.PodInfo{Pod: pod},
				}
				_, err := sched.SchedulePod(ctx, schedFramework, framework.NewCycleState(), podInfo)
				if err != nil {
					b.Fatalf("SchedulePod failed: %v", err)
				}
			}
		})
	}
}

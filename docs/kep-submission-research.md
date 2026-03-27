# KEP Submission Research: ReplicaSet Consolidation Scale-In Strategy

**Date:** 2026-03-27
**Bead:** k8-6x8

---

## 1. KEP Submission Process

### What is a KEP?

A Kubernetes Enhancement Proposal (KEP) is a structured design document for proposing, communicating, and coordinating new efforts in the Kubernetes project. KEPs combine feature tracking, product requirements, and design documentation into a single artifact managed in version control.

### Where KEPs Live

- Repository: `kubernetes/enhancements`
- Path: `keps/sig-apps/NNNN-your-title/`
- Each KEP directory contains:
  - `README.md` — the KEP document (from the template)
  - `kep.yaml` — structured metadata

### Submission Process

1. **Socialize the idea** with the sponsoring SIG (sig-apps). Send to their mailing list or add to a meeting agenda. Get agreement that the work is worth pursuing.
2. **File an enhancement tracking issue** in `kubernetes/enhancements`.
3. **Create the KEP** by opening a PR to `kubernetes/enhancements` adding a new directory under `keps/sig-apps/` with `README.md` (from the template) and `kep.yaml`.
4. **Get SIG approval** — SIG chairs/TLs are the approvers who move the KEP to `implementable`.

### Required KEP Sections (from template)

- Summary
- Motivation (Goals / Non-Goals)
- Proposal (User Stories, Risks and Mitigations)
- Design Details (Test Plan, Graduation Criteria, Upgrade/Downgrade Strategy, Version Skew Strategy)
- Production Readiness Review Questionnaire (Feature Enablement, Rollout Planning, Monitoring, Dependencies, Scalability, Troubleshooting)
- Implementation History
- Drawbacks
- Alternatives

### kep.yaml Metadata

```yaml
title: ReplicaSet Consolidation-Aware Scale-In Strategy
kep-number: NNNN  # assigned when enhancement issue is filed
authors:
  - "@your-github-handle"
owning-sig: sig-apps
status: provisional
creation-date: 2026-03-27
reviewers:
  - TBD
approvers:
  - TBD
stage: alpha
latest-milestone: "v1.37"
milestone:
  alpha: "v1.37"
feature-gates:
  - name: ConsolidatingScaleDown
    components:
      - kube-controller-manager
    disable-supported: true
```

### KEP Lifecycle

| Status | Meaning |
|--------|---------|
| `provisional` | Proposed and actively being defined. SIG has accepted the work direction. |
| `implementable` | Approvers have approved for implementation. |
| `implemented` | Feature is implemented and no longer actively changed. |
| `deferred` | Proposed but not actively worked on. |
| `rejected` | Approvers decided not to proceed. |
| `withdrawn` | Authors withdrew the KEP. |
| `replaced` | Superseded by another KEP. |

### Feature Graduation Path

- **Alpha** — behind a feature gate (disabled by default). Requires: KEP at `implementable`, basic tests, PRR approval.
- **Beta** — feature gate enabled by default. Requires: feedback from alpha, comprehensive tests, monitoring metrics, PRR re-review.
- **GA (Stable)** — feature gate locked on. Requires: conformance tests, 2-week flake-free window, all known bugs fixed.

### Approvals Needed

1. **SIG Approval** — sig-apps chairs/TLs must approve the KEP status transitions.
2. **PRR (Production Readiness Review)** — required before targeting a milestone. PRR reviewers assess operational readiness.
3. **API Review** — if the KEP introduces or modifies API fields (likely applicable since the consolidation-strategy branch adds a feature gate and modifies controller behavior, though it doesn't add new API fields currently).

---

## 2. SIG Ownership

### Owning SIG: sig-apps

The ReplicaSet controller is owned by sig-apps under the **workloads-api** subproject. Relevant code paths:
- `kubernetes/kubernetes/pkg/controller/replicaset/`
- `kubernetes/kubernetes/pkg/controller/controller_utils.go`

### Current Leadership

**Chairs (also serve as Technical Leads):**
- Janet Kuo (@janetkuo) — Google
- Kenneth Owens (@kow3ns) — Snowflake
- Maciej Szulik (@soltysh) — Defense Unicorns

**Steering Committee Liaison:** Antonio Ojea (@aojea)

### Meeting Schedule

- **Regular SIG Meeting:** Mondays at 9:00 AM PT, biweekly
  - Zoom: https://zoom.us/j/739385290?pwd=ekVmNGRjT214MGJkY1JUUUpPMVlJUT09
  - Notes/Agenda: https://docs.google.com/document/d/1LZLBGW2wRDwAfdBNHJjFfk9CFoyZPcIYGWU7R1PQ3ng/edit
  - Recordings: https://www.youtube.com/playlist?list=PL69nYSiGNLP2LMq7vznITnpd2Fk1YIZF3

### Next Meeting from 2026-03-27

The meetings are biweekly on Mondays at 9:00 AM PT. The next meeting would be **Monday, March 30, 2026** or **Monday, April 6, 2026** (depending on the biweekly cadence — check the calendar link to confirm which week is active). Note: KubeCon EU is March 23-26, so the March 30 meeting may be affected.

- Calendar: https://calendar.google.com/calendar/embed?src=phfni1v25vnmi4q06m851230so%40group.calendar.google.com

### Getting on the Agenda

Add your topic to the shared Google Doc agenda (linked above) before the meeting. Include:
- Brief description of the proposal
- Link to the KEP PR or draft
- What you're asking for (feedback, sponsorship, approval)

### Communication Channels

- **Slack:** `#sig-apps` on kubernetes.slack.com
- **Mailing list:** https://groups.google.com/a/kubernetes.io/g/sig-apps
- **GitHub teams:** @kubernetes/sig-apps-proposals (for design proposals)

---

## 3. Related Work

### KEP-2255: ReplicaSet Pod Deletion Cost (BETA since K8s 1.22)

- **Location:** `keps/sig-apps/2255-pod-cost/`
- **Feature gate:** `PodDeletionCost`
- **What it does:** Adds `controller.kubernetes.io/pod-deletion-cost` annotation. Pods with lower cost are deleted first during scale-down.
- **Limitations:** Annotation must be set BEFORE scale-down. Doesn't integrate with HPA. Continuously updating annotations is an anti-pattern (API server load).
- **Status:** Beta, not yet GA.
- **Relationship to our work:** Our consolidation strategy is complementary — it adds a node-topology-aware heuristic that doesn't require annotation management. We should reference KEP-2255 and explain how our approach addresses its limitations for the consolidation use case.

### Issue #107598: Configurable Down-Scaling Behaviour (Closed)

- **Author:** @thesuperzapper (2022)
- **Proposal:** Add a `scaleConfig` field to ReplicaSetSpec defining a sequence of sorting methods for scale-down pod selection (cost annotations, node replicas, pod age, external API, heuristic probes).
- **Status:** Closed as not planned (went stale/rotten). Superseded by #123541.
- **Relevance:** Very ambitious scope. Our proposal is narrower and more focused — a single consolidation heuristic behind a feature gate, not a pluggable framework.

### Issue #123541: Extend Scale Subresource for pod-deletion-cost (Closed)

- **Author:** @thesuperzapper (2024)
- **Proposal:** Extend the `Scale` v1 subresource API to read/write pod-deletion-cost annotations, enabling atomic scale+cost updates.
- **Status:** Closed as not planned (went stale/rotten).
- **Relevance:** Different approach (API-level). Our work modifies the controller's internal sorting logic rather than exposing new API surface.

### Other Related Issues

- **#45509** (2017): Original request for customizable scale-down behavior
- **#4301** (2015): Early discussion about pod deletion ordering
- **kubernetes/enhancements#2255**: Enhancement tracking issue for PodDeletionCost
- **kubernetes/enhancements#3189**: Deployment/ReplicaSet Downscale Pod Picker

### Key Observation

All previous proposals (#107598, #123541) were closed as stale/rotten without being implemented beyond KEP-2255. This suggests:
1. The community recognizes the need but hasn't found the right scope/approach
2. Our focused, feature-gated approach (single heuristic, no API changes) may be more palatable
3. We need strong SIG sponsorship to avoid the same fate

---

## 4. What Our Change Does (consolidation-strategy branch)

The branch adds a `ConsolidatingScaleDown` feature gate that modifies ReplicaSet scale-down pod selection:

1. **Node pod density ranking:** Instead of preferring to delete pods on nodes with MORE colocated replicas (spreading heuristic), it prefers deleting pods on nodes with FEWER total active pods (consolidation heuristic). This helps Karpenter/cluster-autoscaler consolidate workloads onto fewer nodes.

2. **Do-not-disrupt awareness:** Checks for `karpenter.sh/do-not-disrupt: "true"` annotation on pods. Pods on nodes with do-not-disrupt pods are deprioritized for deletion.

3. **Node informer integration:** Adds optional node informer to the ReplicaSet controller for future resource-based scoring.

**Files modified:**
- `pkg/features/kube_features.go` — new feature gate
- `pkg/controller/controller_utils.go` — modified `ActivePodsWithRanks.Less()` sorting
- `pkg/controller/replicaset/replica_set.go` — new ranking function, node informer wiring
- `cmd/kube-controller-manager/app/apps.go` — node informer passed to controller

---

## 5. Draft KEP Outline

### KEP-NNNN: ReplicaSet Consolidation-Aware Scale-In Strategy

**Summary:** Add an opt-in consolidation heuristic to the ReplicaSet controller's scale-down pod selection algorithm, enabling workload consolidation onto fewer nodes when scaling down.

**Motivation:**
- Current scale-down prefers deleting pods on nodes with more colocated replicas (spreading). This works against node consolidation.
- Cluster autoscalers (Karpenter, cluster-autoscaler) can reclaim empty/underutilized nodes, but only if workloads consolidate during scale-down.
- KEP-2255 (PodDeletionCost) requires external annotation management, which is complex and doesn't integrate with HPA.

**Goals:**
- Provide a feature-gated consolidation-aware pod deletion heuristic
- Respect do-not-disrupt signals from node autoscalers
- Maintain backward compatibility (feature gate off = current behavior)

**Non-Goals:**
- Replace or deprecate KEP-2255
- Add new API fields to ReplicaSet/Deployment specs
- Implement a general-purpose pluggable scale-down framework

**Design Details:**
- Feature gate: `ConsolidatingScaleDown` (alpha, disabled by default)
- Modified sorting in `ActivePodsWithRanks.Less()`:
  - Step 4.5: Deprioritize pods on do-not-disrupt nodes
  - Step 5: Invert rank comparison (prefer deleting from nodes with fewer pods)
- Node informer added to ReplicaSet controller for pod-per-node counting
- No API changes required

**Graduation Criteria:**
- Alpha: Feature gate, unit tests, integration tests
- Beta: Metrics, feedback from alpha users, e2e tests
- GA: Conformance tests, 2-week flake-free window

**Risks:**
- Karpenter-specific annotation (`karpenter.sh/do-not-disrupt`) couples core K8s to a specific autoscaler. May need a more generic mechanism.
- Consolidation heuristic may conflict with pod topology spread constraints.

---

## 6. Concrete Next Steps

### Immediate (This Week)

1. **Join #sig-apps Slack** and introduce the proposal briefly
2. **Add to sig-apps meeting agenda** for the next available meeting (check calendar for March 30 or April 6)
3. **File enhancement tracking issue** in `kubernetes/enhancements` with title "ReplicaSet Consolidation-Aware Scale-In Strategy"

### Short-Term (Next 2 Weeks)

4. **Write the full KEP** using the template, addressing all required sections including PRR questionnaire
5. **Open KEP PR** to `kubernetes/enhancements` under `keps/sig-apps/NNNN-consolidation-scale-in/`
6. **Present at sig-apps meeting** — request sponsorship from chairs/TLs

### Medium-Term (v1.37 Targeting)

7. **Address feedback** from SIG review and PRR
8. **Get KEP to `implementable`** status
9. **Open implementation PR** against `kubernetes/kubernetes` (based on consolidation-strategy branch, cleaned up)
10. **Target v1.37 alpha** (v1.36 enhancements freeze has already passed as of Feb 11, 2026)

### Key Contacts

| Role | Person | GitHub |
|------|--------|--------|
| sig-apps Chair | Janet Kuo | @janetkuo |
| sig-apps Chair | Kenneth Owens | @kow3ns |
| sig-apps Chair | Maciej Szulik | @soltysh |
| Steering Liaison | Antonio Ojea | @aojea |

### Timeline Note

- **v1.36** is in code freeze (March 18) → release April 22, 2026. Too late for v1.36.
- **v1.37** cycle will start ~May 2026. Target enhancements freeze for v1.37.
- KubeCon EU (March 23-26, 2026) is happening NOW — excellent opportunity for hallway conversations with sig-apps members if attending.

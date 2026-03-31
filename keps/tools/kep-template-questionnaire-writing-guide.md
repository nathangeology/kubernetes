# KEP Template, Questionnaire, and LLM Writing Guide

## Changelog

- 2026-03-31: Improvements based on replicaset-consolidation-scale-in KEP authoring experience
  - Added User Stories questions and template to Proposal section (Section 5)
  - Added Release Signoff Checklist template (Section 0)
  - Added Drawbacks section (new Section 8, renumbered subsequent sections)
  - Added Feature Gate design detail questions to Design Details (Section 6)
  - Added Version Skew Strategy question to Design Details (Section 6)
  - Added RBAC/permission change questions to Proposal (Section 5)
  - Added before/after behavior comparison table prompt to Proposal (Section 5)
  - Restructured PRR questionnaire (Section 12) with specific sub-questions matching upstream KEP format
  - Consolidated redundant questions between Design Details and PRR sections
  - Added guidance on structured comparison tables and user story format to writing guide
  - Moved Risks and Mitigations into Proposal section to match actual KEP structure
  - Added pitfall notes from real authoring experience

---

## SECTION 0: YAML Metadata Block [REQUIRED]

### Questionnaire

1. What is the KEP number (assigned by SIG leads or leave as NNNN)?
2. What is the title of your enhancement?
3. Who are the authors (GitHub handles)?
4. Which SIG owns this KEP (e.g., sig-autoscaling, sig-node)?
5. What is the current status? (provisional / implementable / implemented / deferred / rejected / withdrawn / replaced)
6. What stage is this targeting? (alpha / beta / stable)
7. What Kubernetes release are you targeting (e.g., v1.36)?
8. What is the feature gate name?
9. Who are the reviewers and approvers (GitHub handles)?
10. Creation date (YYYY-MM-DD)?

### Template

```yaml
---
kep-number: NNNN
title: "[Your Title]"
authors:
  - "@github-handle"
owning-sig: sig-XXXXX
status: provisional
creation-date: YYYY-MM-DD
reviewers:
  - "@reviewer1"
approvers:
  - "@approver1"
stage: alpha
latest-milestone: "vX.XX"
feature-gates:
  - name: YourFeatureGateName
milestone:
  alpha: "vX.XX"
  beta: "vX.XX"
  stable: "vX.XX"
---
```

### Release Signoff Checklist

Include this checklist after the YAML metadata and table of contents:

```markdown
## Release Signoff Checklist

Items marked with (R) are required *prior to targeting to a milestone / release*.

- [ ] (R) Enhancement issue in release milestone, which links to KEP dir in [kubernetes/enhancements]
- [ ] (R) KEP approvers have approved the KEP status as `implementable`
- [ ] (R) Design details are appropriately documented
- [ ] (R) Test plan is in place, giving consideration to SIG Architecture and SIG Testing input
  - [ ] e2e Tests for all Beta API Operations (CLI, key scenarios, etc.)
  - [ ] (R) Ensure GA e2e tests meet requirements for [Coverage Onboarding Checklist]
- [ ] (R) Graduation criteria is in place
  - [ ] (R) [all GA Coverage Onboarding Checklist items] are met
- [ ] (R) Production readiness review completed
- [ ] (R) Production readiness review approved
- [ ] "Implementation History" section is up-to-date for milestone
- [ ] User-facing documentation has been created in [kubernetes/website]
- [ ] Supporting documentation (additional design documents, mailing list discussions, relevant PRs/issues, release notes)
```

### LLM Writing Instructions

→ Populate directly from questionnaire answers. Use exact GitHub handles with @ prefix. Leave fields as TBD only if the user explicitly says they don't know yet. Use ISO date format.
→ Always include the Release Signoff Checklist — it is required by the KEP process and reviewers will ask for it if missing.

## SECTION 1: Summary [REQUIRED]

### Questionnaire

1. In 2-4 sentences, what does this KEP propose?
2. What is the single most important outcome if this KEP is implemented?

### Template

## Summary

[2-4 sentence summary of the proposal and its key outcome]

### LLM Writing Instructions

→ Write exactly 2-4 sentences. First sentence: what the proposal does. Second sentence: why it matters. Optional third/fourth: key mechanism or scope boundary.
→ Trade-off — Brevity vs. Completeness: The summary must be self-contained enough that a reader can decide whether to read further, but short enough to scan in 10 seconds. Ask the user: "Should the summary emphasize the technical mechanism or the user-facing benefit?" Default to user-facing benefit if unclear.
→ Do NOT include implementation details, formulas, or code. Use present tense ("This KEP introduces...").

## SECTION 2: Motivation [REQUIRED]

### Questionnaire

1. What is the current behavior or limitation that this KEP addresses?
2. What specific problems do users experience today? (Include links to GitHub issues if available)
3. Why can't this be solved with existing features or configuration?
4. How widespread is this problem? (Rough estimate: number of issues, affected users, frequency)

### Template

## Motivation

[Description of the current state and its limitations]

[Description of user pain points with issue links]

[Why existing solutions are insufficient]

### Goals

- [Goal 1]
- [Goal 2]
- [Goal 3]

### Non-Goals

- [Non-goal 1]
- [Non-goal 2]

### LLM Writing Instructions

→ Start with the current state, then transition to pain points. Link GitHub issues inline using [#NNNN](url) format.
→ Trade-off — Depth vs. Accessibility: Deep technical motivation helps expert reviewers but alienates newcomers. Ask the user: "Is your primary audience SIG members who already know the system, or a broader Kubernetes community audience?" For expert audiences, lead with technical specifics. For broader audiences, lead with user-visible symptoms and add technical details after.
→ Use concrete examples over abstract descriptions. "Nodes at 95% CPU are disrupted" is better than "highly utilized nodes may be unnecessarily disrupted."
→ Do NOT propose solutions in this section.

## SECTION 3: Goals [REQUIRED]

### Questionnaire

1. List 3-5 specific, measurable outcomes this KEP aims to achieve.
2. For each goal, how would you verify it was achieved?

### LLM Writing Instructions

→ Each goal should be a single bullet point, starting with a verb. Goals should be testable/verifiable. Avoid vague goals like "improve performance" — prefer "reduce unnecessary node consolidation events by X% in clusters with mixed workload priorities."
→ Keep to 3-5 goals. If the user provides more, help them consolidate or move extras to Non-Goals or Future Work.

## SECTION 4: Non-Goals [REQUIRED]

### Questionnaire

1. What related problems does this KEP explicitly NOT try to solve?
2. Are there adjacent features that readers might expect this KEP to cover but shouldn't?
3. Are any of these non-goals planned for future KEPs?

### LLM Writing Instructions

→ Frame each non-goal as a clear boundary statement. Use "This KEP does not..." or "Out of scope:..." format.
→ If a non-goal is planned for future work, note it: "Out of scope for this KEP but planned as a follow-up: ..."
→ Non-goals prevent scope creep during review. Include anything a reasonable reviewer might ask "why doesn't this also do X?"

## SECTION 5: Proposal [REQUIRED]

### Questionnaire

1. What is the proposed solution at a high level?
2. What is the proposed API change or new API? (Provide YAML/Go snippets if applicable)
3. How does the new behavior differ from the current behavior? Describe the current step-by-step algorithm/pipeline and show exactly which steps change, using a before/after comparison table if applicable.
4. What are the key design decisions and why were they made?
5. Are there any new CRD fields, flags, annotations, or labels?
6. Are there any RBAC or permission changes required? (e.g., new ClusterRole verbs, service account bindings)
7. Walk through 2-3 concrete user stories showing who benefits and how. Format: "As a [role], I want [action] so that [outcome]."
8. Walk through 2-3 concrete examples showing the proposal in action with specific numbers.

### Template

## Proposal

### User Stories

#### Story 1: [Title]

As a [role], I [want/have] [action/situation] so that [outcome], [without requiring workaround].

#### Story 2: [Title]

As a [role], I [want/have] [action/situation] so that [outcome].

[High-level description of the solution]

### Proposed API/Spec

[YAML or Go code showing the API surface, or "No API changes" if none]

### RBAC Changes

[Description of any RBAC/permission changes, or "No RBAC changes required"]

### How It Works

[Detailed explanation of the mechanism]

#### Before/After Comparison

| Step | Current Behavior | New Behavior (gate enabled) | Gate-Dependent? |
|------|-----------------|---------------------------|-----------------|
| 1 | [current] | [new or unchanged] | No |
| N | [current] | [CHANGED: new behavior] | Yes |

### Risks and Mitigations

| Risk | Severity | Mitigation |
|------|----------|------------|
| [Risk 1] | [High/Medium/Low] | [Mitigation] |

### Examples

[2-3 worked examples with concrete numbers]

### LLM Writing Instructions

→ Trade-off — Precision vs. Readability: Mathematical formulas and pseudocode are precise but can lose readers. Narrative explanations are accessible but can be ambiguous. Ask the user: "Do you want to lead with the formula/algorithm and then explain it, or lead with intuition and then formalize?" The recommended pattern is leading with intuition then formalizing.
→ Start with the high-level mechanism in 2-3 sentences, then drill into specifics.
→ For API changes, show a complete YAML example that a user could copy-paste.
→ User stories should be concrete and specific to real personas (platform engineer, cost-conscious operator, SRE), not generic "as a user" stories.
→ The before/after comparison table is critical for proposals that modify existing algorithms or pipelines. Show every step, mark which ones change, and indicate gate-dependency. Reviewers rely on this to understand the blast radius.
→ Worked examples should cover: (a) a case where the proposal approves an action, (b) a case where it rejects, and (c) an edge case.
→ Risks and Mitigations belong here under Proposal, not in a separate top-level section. Reviewers expect to see risks immediately after understanding the proposal. Include at least: correctness risks, performance risks, operational risks, and any coupling risks (e.g., dependency on external project annotations).
→ Feature gates are always a mitigation — mention them.
→ Use consistent terminology throughout. Define terms on first use.
→ Common pitfall: If your proposal reads an annotation or label from another project (e.g., Karpenter), explicitly address the coupling risk and describe the migration path to a Kubernetes-native equivalent.

## SECTION 6: Design Details [REQUIRED]

### Questionnaire

1. What are the edge cases and how are they handled? For each, state: the scenario, the expected behavior, and why that behavior is correct.
2. How does this interact with existing features? (List each feature and the interaction)
3. Are there any new metrics, logs, or events? (If none in alpha, state what is planned for beta)
4. What happens during upgrade/downgrade? Is any data migration required?
5. Can the feature be disabled after enabling? What happens to already-created resources?
6. Feature gate details: What is the gate name, which component(s) depend on it, what is the default value, and is disable supported?
7. What is the version skew strategy? How does this feature behave when control plane components are at different versions?

### Template

## Design Details

### Feature Gate

- Name: `YourFeatureGateName`
- Component: `component-name`
- Default: `false` (alpha)
- Disable-supported: `true`

When disabled, [describe exact behavior — must be identical to current upstream behavior].

### Edge Cases

- **[Edge case name]**: [Scenario]. [Expected behavior]. [Why this is correct].

### Interaction with Existing Features

- **[Feature X]**: [How this KEP interacts with Feature X].
- **[Feature Y]**: No interaction expected.

### Observability

[New metrics, logs, events — or "No new observability in alpha. Planned for beta: ..."]

### Test Plan

[x] I/we understand the owners of the involved components may require updates to existing tests to make this code solid enough prior to committing the changes necessary to implement this enhancement.

#### Unit Tests

- [Test area and coverage target]

#### Integration Tests

- [Test scenario]

#### End-to-End Tests

- [Test scenario, or "Planned for beta"]

### Graduation Criteria

#### Alpha

- [Criterion 1]

#### Beta

- [Criterion 1]

#### GA

- [Criterion 1]

### Upgrade / Downgrade Strategy

[What happens on upgrade. What happens on downgrade. Is data migration required?]

### Version Skew Strategy

[How does this feature behave when control plane components are at different versions? Is there a version skew concern between API server, controller manager, scheduler, kubelet?]

### LLM Writing Instructions

→ Be exhaustive on edge cases — reviewers will probe these. For each edge case, state the scenario, the expected behavior, and why that behavior is correct.
→ For feature interactions, use a consistent format: "Feature X: [How this KEP interacts with Feature X]." Cover the most important interactions in detail and briefly note others as "no interaction expected."
→ The Feature Gate subsection is required. Reviewers need to know the exact gate name, component, default, and what happens when disabled. The "when disabled" description must guarantee identical behavior to upstream.
→ Version Skew Strategy is required for any feature that spans multiple components. Even if the feature is in a single component, state that explicitly ("no version skew concern because...").
→ Test Plan, Graduation Criteria, Upgrade/Downgrade, and Version Skew all belong under Design Details in the final KEP structure, not as separate top-level sections.
→ Common pitfall: Don't forget to describe what happens when the node/pod informer cache hasn't synced yet during controller startup. This is a common edge case for controller-based features.

## SECTION 7: Alternatives Considered [REQUIRED]

### Questionnaire

1. What alternative approaches did you consider? (List at least 2-3, give each a descriptive name)
2. For each alternative, what are the pros and cons?
3. For each alternative, provide a concrete scenario where it produces the wrong outcome (if applicable).
4. Are there any related PRs, proposals, or prior art? List each with a brief description of what it proposed and why it didn't ship.

### Template

## Alternatives

### [Alternative 1 Name]

[Description, pros/cons, and why it was rejected]

### [Alternative 2 Name]

[Description, pros/cons, and why it was rejected]

## Drawbacks

[Honest assessment of the downsides of this proposal, even if mitigated. Include complexity cost, resource overhead, and any tradeoffs users should be aware of.]

### LLM Writing Instructions

→ For each alternative, use the pattern: describe the alternative, then present a concrete scenario where it fails. This is more persuasive than abstract pros/cons lists.
→ Include at least 2-3 alternatives. "Do nothing" is a valid alternative.
→ Trade-off — Advocacy vs. Objectivity: The author naturally favors their proposal, but reviewers trust documents that honestly present alternatives. Ask the user: "Are there any alternatives that are genuinely close calls? Where is the strongest argument against your proposal?" Present those fairly — it builds credibility.
→ Link to related PRs and prior art with brief descriptions. This shows awareness of the ecosystem and prevents reviewers from asking "did you consider X?" For each prior proposal, note: what it proposed, when, and why it didn't ship.
→ The Drawbacks section is separate from Risks. Risks are things that could go wrong; Drawbacks are known downsides that exist even when everything works correctly (e.g., increased complexity, memory overhead, potential for uneven distribution). Reviewers expect honest self-assessment here.

## SECTION 8: Production Readiness Review (PRR) Questionnaire [REQUIRED]

### Questionnaire

Answer ALL of the following sub-questions. These match the upstream KEP PRR format exactly. Do not write "N/A" — explain why something doesn't apply.

#### Feature Enablement and Rollback

1. How can this feature be enabled / disabled in a live cluster? (Feature gate name, component, HA rolling restart support?)
2. Does enabling the feature change any default behavior?
3. Can the feature be disabled once it has been enabled (i.e. can we roll back the enablement)?
4. What happens if we reenable the feature if it was previously rolled back?
5. Are there any tests for feature enablement/disablement?

#### Rollout, Upgrade and Rollback Planning

6. How can a rollout or rollback fail? Can it impact already running workloads?
7. What specific metrics should inform a rollback?
8. Were upgrade and rollback tested? Was the upgrade→downgrade→upgrade path tested?
9. Is the rollout accompanied by any deprecations and/or removals of features, APIs, machines, permissions, or service account bindings?

#### Monitoring Requirements

10. How can an operator determine if the feature is in use by workloads?
11. How can someone using this feature know that it is working for their instance?
12. What are the SLIs (Service Level Indicators) an operator can use to determine the health of the service?
13. Are there any missing metrics that would be useful to have in this category?
14. What are the reasonable SLOs (Service Level Objectives) for the enhancement?
15. What are the failure modes for the enhancement? (Use a table: Failure Mode | Impact | Detection | Mitigation)

#### Dependencies

16. Does this feature depend on any specific services running in the cluster?

#### Scalability

17. Will enabling / using this feature result in any new API calls?
18. Will enabling / using this feature result in introducing new API types?
19. Will enabling / using this feature result in any new calls to the cloud provider?
20. Will enabling / using this feature result in increasing size or count of the existing API objects?
21. Will enabling / using this feature result in increasing time taken by any operations?
22. Will enabling / using this feature result in non-negligible increase of resource usage (CPU, memory)?
23. Can enabling / using this feature result in resource exhaustion of some node resources?
24. Will enabling / using this feature result in any new resource claim usage?

#### Troubleshooting

25. How does this feature react if the API server and/or etcd is unavailable?
26. What are other known failure modes?
27. What steps should be taken if SLOs are not being met to determine the problem?

### Template

## Production Readiness Review Questionnaire

### Feature Enablement and Rollback

###### How can this feature be enabled / disabled in a live cluster?

- [x] Feature gate
  - Feature gate name: `YourFeatureGateName`
  - Components depending on the feature gate: `component-name`

###### Does enabling the feature change any default behavior?

[Answer]

###### Can the feature be disabled once it has been enabled (i.e. can we roll back the enablement)?

[Answer]

###### What happens if we reenable the feature if it was previously rolled back?

[Answer]

###### Are there any tests for feature enablement/disablement?

[Answer]

### Rollout, Upgrade and Rollback Planning

###### How can a rollout or rollback fail? Can it impact already running workloads?

[Answer]

###### What specific metrics should inform a rollback?

[Answer]

###### Were upgrade and rollback tested? Was the upgrade→downgrade→upgrade path tested?

[Answer]

###### Is the rollout accompanied by any deprecations and/or removals of features, APIs, machines, permissions, or service account bindings?

[Answer]

### Monitoring Requirements

###### How can an operator determine if the feature is in use by workloads?

[Answer]

###### How can someone using this feature know that it is working for their instance?

[Answer]

###### What are the SLIs (Service Level Indicators) an operator can use to determine the health of the service?

- [x] Metrics
  - Metric name: [name]
  - Components exposing the metric: [component]

###### Are there any missing metrics that would be useful to have in this category?

[Answer]

###### What are the reasonable SLOs (Service Level Objectives) for the enhancement?

[Answer]

###### What are the failure modes for the enhancement?

| Failure Mode | Impact | Detection | Mitigation |
|---|---|---|---|
| [Mode 1] | [Impact] | [Detection] | [Mitigation] |

###### What steps should be taken if SLOs are not being met to determine the problem?

[Answer]

### Dependencies

###### Does this feature depend on any specific services running in the cluster?

[Answer]

### Scalability

###### Will enabling / using this feature result in any new API calls?

[Answer]

###### Will enabling / using this feature result in introducing new API types?

[Answer]

###### Will enabling / using this feature result in any new calls to the cloud provider?

[Answer]

###### Will enabling / using this feature result in increasing size or count of the existing API objects?

[Answer]

###### Will enabling / using this feature result in increasing time taken by any operations?

[Answer]

###### Will enabling / using this feature result in non-negligible increase of resource usage?

[Answer — include specific estimates for memory and CPU]

###### Can enabling / using this feature result in resource exhaustion of some node resources?

[Answer]

###### Will enabling / using this feature result in any new resource claim usage?

[Answer]

### Troubleshooting

###### How does this feature react if the API server and/or etcd is unavailable?

[Answer]

###### What are other known failure modes?

[Answer]

###### What steps should be taken if SLOs are not being met to determine the problem?

[Answer]

### LLM Writing Instructions

→ This section is reviewed by the PRR team and is a gate for inclusion in a release. Be thorough.
→ Answer every question directly. Do not write "N/A" — explain why something doesn't apply.
→ Use the exact ###### heading format shown above — this matches the upstream KEP template and PRR reviewers expect it.
→ The failure modes table is critical. Include at least: cache sync failures, informer failures, and annotation parsing edge cases.
→ For scalability answers, provide specific estimates (e.g., "approximately 5-10 MB additional memory for a 5,000-node cluster") rather than vague statements.
→ Common pitfall: Don't say "will be tested during alpha" for the upgrade/rollback question. Describe the test plan concretely, even if it hasn't been executed yet.

## SECTION 9: Backward Compatibility [REQUIRED]

### Questionnaire

1. Does this change break any existing API contracts?
2. Do existing configurations continue to work unchanged?
3. Is there a migration path for users on the old behavior?

### Template

## Backward Compatibility

[Compatibility statement and migration details]

### LLM Writing Instructions

→ State explicitly: "Existing [resource/config] continues to work unchanged." Be specific about what is unchanged.
→ If there is a migration path, describe it step by step.
→ Address unconditional changes (like RBAC additions) even when the feature gate is off.

## SECTION 10: Future Work [OPTIONAL]

### Questionnaire

1. What follow-up enhancements are planned but out of scope for this KEP?
2. What would you build if this KEP is successful?

### Template

## Future Work

### [Future Item 1]

[Brief description]

### LLM Writing Instructions

→ Keep each item to 2-4 sentences. This section sets expectations without committing to scope.
→ Link back to Non-Goals where applicable.

## SECTION 11: Implementation History [OPTIONAL — becomes REQUIRED at Beta]

### Template

## Implementation History

- YYYY-MM-DD: KEP created
- YYYY-MM-DD: Alpha implementation merged (vX.XX)

### LLM Writing Instructions

→ Update this section as milestones are reached. Use ISO dates.

## SECTION 12: Infrastructure Needed [OPTIONAL]

### Questionnaire

1. Do you need new CI jobs, test clusters, or other infrastructure?
2. Do you need new test frameworks or tools?

### Template

## Infrastructure Needed

[Infrastructure requirements or "No new infrastructure needed."]

---

## Global Writing Trade-offs Questionnaire

Answer these questions to help the LLM calibrate the overall writing style:

1. **Technical depth vs. accessibility**: Is your primary audience deeply familiar with the codebase, or does this need to be understood by the broader Kubernetes community? (Scale: 1=broad community, 5=core maintainers only)
2. **Conciseness vs. thoroughness**: Do you prefer a shorter, punchier document or a comprehensive one that anticipates every reviewer question? (Scale: 1=minimal, 5=exhaustive)
3. **Formality vs. conversational tone**: Should this read like an academic paper or like a well-structured design discussion? (Scale: 1=conversational, 5=formal/academic)
4. **Formula-first vs. intuition-first**: When presenting algorithms or scoring mechanisms, should the document lead with the mathematical formulation or with intuitive examples? (Scale: 1=intuition first, 5=formula first)
5. **Advocacy vs. neutrality**: How strongly should the document advocate for the proposed solution vs. presenting options neutrally? (Scale: 1=neutral presentation, 5=strong advocacy with clear recommendation)

---

## LLM Master Style Guide

### Voice and Tone

→ Use active voice and present tense ("This KEP introduces..." not "This KEP will introduce...")
→ Write in third person for the proposal itself, second person for instructions to operators ("Operators can enable...")
→ Be direct and specific. Avoid hedging language ("might," "could potentially," "it is possible that")
→ Match the technical level to the audience (see trade-off question #1)

### Structure

→ Every section should be self-contained enough to be reviewed independently
→ Use headers liberally — reviewers often jump to specific sections
→ Lead with the conclusion/recommendation, then provide supporting detail
→ Use tables for comparisons, bullet lists for enumerations, prose for narratives
→ The final KEP structure nests Test Plan, Graduation Criteria, Upgrade/Downgrade, and Version Skew under Design Details — not as top-level sections

### Technical Writing

→ Define acronyms on first use
→ Use consistent terminology — pick one term and stick with it throughout
→ When presenting formulas, always follow with a worked example
→ When describing behavior, always cover: normal case, edge case, error case
→ Link to source code, issues, and PRs using full URLs
→ When describing changes to an existing algorithm or pipeline, always provide a before/after comparison table showing every step, which ones change, and whether the change is gate-dependent

### User Stories

→ User stories should use concrete personas (platform engineer, SRE, cost-conscious operator) not generic "as a user"
→ Each story should describe a real scenario with enough specificity that a reviewer can evaluate whether the proposal actually solves it
→ Include at least one story that demonstrates the do-nothing/workaround pain to motivate the change

### Common Pitfalls to Avoid

→ Don't mix motivation with proposal — keep them in separate sections
→ Don't present alternatives as straw men — give them fair treatment
→ Don't leave PRR questions as "TBD" — reviewers will block on this
→ Don't forget to address what happens when the feature is disabled after being enabled
→ Don't assume the reader has read related KEPs — provide enough context inline
→ Don't forget the Drawbacks section — it's separate from Risks and shows honest self-assessment
→ Don't forget Version Skew Strategy even for single-component features — state "no concern because..." explicitly
→ Don't forget RBAC changes — even read-only permissions need to be documented
→ Don't put Risks in a separate top-level section — they belong under Proposal where reviewers expect them
→ When your proposal reads annotations from external projects (Karpenter, Istio, etc.), always document the coupling risk and the migration path to a Kubernetes-native equivalent

### Formatting

→ Use GitHub-flavored Markdown
→ Code blocks should specify the language (yaml, go, etc.)
→ Keep line length reasonable for diff review (~80-100 chars in code blocks)
→ Use monospace for field names, API objects, and CLI flags
→ Use `######` (h6) for PRR sub-questions to match upstream KEP format

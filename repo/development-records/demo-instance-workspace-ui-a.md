# DEMO-INSTANCE-WORKSPACE-UI-A: Production-Oriented Instance Workspace

## Status

Completed.

## Purpose

Redesign the staged instance demo as a production-oriented console reference.
The page is still backed by demo APIs, but the information architecture,
operation grouping, and lifecycle coverage are intended to be reusable by the
future production console.

## Scope

- Replaces the simple demo table with an instance workspace:
  - resource-domain navigation for instance, image, network, storage, GPU pool,
    console sessions, tasks, and audit
  - create blueprint panel
  - instance list with filters
  - selected instance detail panel
  - detail tabs for overview, metrics, network, storage, security, events,
    audit, snapshots, and backups
  - lifecycle action panel
  - ops/action output panel
  - create timeline panel
- Covers all current demo instance classes:
  - VM
  - container
  - GPU container
- Exposes full staged lifecycle actions:
  - start
  - stop
  - restart
  - resize
  - delete
- Exposes staged operations:
  - VM console in a separate page
  - container logs/terminal/exec/events
  - GPU container logs/terminal/exec/metrics
- Adds clickable demo surfaces outside the instance lifecycle:
  - top-level dashboard, model, knowledge-base, usage, and settings pages
  - resource-domain pages for image, network, storage, GPU pools, console
    sessions, tasks, and audit
  - mock tables, metrics, detail panels, and action feedback for presentation
    continuity before their production APIs are implemented
- Reframes the console shell with a cloud-console layout:
  - dark global top bar
  - service search affordance
  - left service navigation
  - breadcrumb and dense workspace content area
  - Evergreen/KuberCloud brand icon loaded from `assets/brand`

## Guardrails

- The UI calls Gateway demo APIs only; it does not bypass `WorkloadInstanceService`.
- VM console opens a separate `/demo-console` page from the workspace, matching
  public-cloud console navigation rather than embedding shell state in the list.
- Non-instance resource domains are represented as production information
  architecture mock surfaces until their M1/M2/M3 backend APIs are promoted.
- Mock pages must not bypass the real instance API. VM/container/GPU container
  creation, lifecycle, ops, and VM console paths remain wired to Gateway demo
  APIs.
- Current shell execution is a demo backend workspace; production must replace
  it with noVNC/WebSocket or cloud-provider console proxy while preserving the
  same UI flow.

## Validation

```bash
npm run build
make validate-demo-instances
git diff --check
```

## Release Impact

`MINOR`: frontend demo UX enhancement with no core provider contract change.

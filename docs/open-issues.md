# GoMasterForge — Open Issues

## Decisions That Block Implementation

### OI-001 — Events vs. diffs (blocks Phase 2)
**Source:** concept § 2
Should state broadcast emit typed domain events or whole-state diffs?
Typed events are the preferred direction: lower SSH bandwidth, explicit change context, and composable with the event bus in `internal/events/`.
Must be resolved before `GameState` design begins.

### OI-002 — Guest players (blocks Phase 4)
**Source:** concept § 0
Should unauthenticated players (no profile) be supported in v1?
Directly affects whether the connect-time profile gate is mandatory for all connections.

### OI-003 — GM-only main menu items
**Source:** concept § 1.2
Should **Start Campaign** and **Mods** be restricted to GM-role profiles only, or visible to all connected users?
A player profile probably should not be able to start a new campaign.

### OI-004 — Session reconnection policy
**Source:** review.md
No rule exists for what happens when a player SSH session drops mid-campaign.
- Does their token stay on the map?
- Do they re-join and get state restored?
- Is there a timeout before they are considered disconnected?
Needs a rule or at minimum an explicit open question in concept § 2.

---

## Decisions With Implicit Answers (Close These)

### OI-005 — Debounce and checkpoint defaults
**Source:** concept § 8
`data/settings/server.json` already defines `debounce = 2s`, `checkpoint = 5m`, `snapshot_retention = 10`.
Action: write these as the documented defaults in concept § 8 and remove the open question.

### OI-007 — Sidebar width breakpoint
**Source:** concept § 4
"Exact width breakpoint for sidebar-to-tab switch?" — 120 columns is a reasonable default.
Action: pick 120, document it, close the question.

---

## Lower Priority / Post-v1

### OI-008 — Mod conflict resolution strategy
**Source:** concept § 3
Should mod conflicts fail hard at load time, or allow override precedence with warnings?

### OI-009 — OGL linter (duplicate — consolidate)
**Source:** concept § 3 and § 10 (same question appears twice)
Should a mod linter block load/publish when required OGL legal metadata is missing from the manifest?
Action: consolidate into one question and remove the duplicate.

### OI-010 — GM impersonation of players
**Source:** concept § 5
Should the GM be able to impersonate player actions for troubleshooting purposes?

### OI-011 — Macro rollback on partial failure
**Source:** concept § 6
Should macro execution support transactional rollback if a step fails mid-execution?

### OI-012 — Chat rate limiting
**Source:** concept § 6
Is per-user chat rate limiting needed in v1 to prevent spam or abuse?

### OI-013 — Translation key validation at startup
**Source:** concept § 7
Should missing required translation keys hard-fail at startup, or warn and continue?

### OI-014 — Campaign override storage strategy
**Source:** concept § 9
Should campaign-level overrides be stored as full copies of mod templates, or as patch deltas against the base mod version?

### OI-015 — Cross-domain ID ownership
**Source:** concept § 9
Which package owns cross-domain IDs and shared type definitions — `mods/core` schema contracts or an `internal/` model package?

---

## Structural Gaps

### OI-017 — `docs/adrs/` directory missing
Phase 0 requires ADRs to be written here. Directory does not exist yet.

### OI-018 — `schemas/manifest.json` manifest schema not written
Phase 5 depends on this. Should be the first artifact committed to `schemas/`.

---

## Closed

### OI-016 — `go.mod` declares `go 1.26`
Closed: Go `1.26` is valid (`go1.26.1` is current).

### OI-019 — `mods/core/manifest.json` and `mods/core-dnd35/manifest.json` are empty
Closed: both manifests now contain full metadata and entrypoints.

### OI-020 — Race detection deferred to Phase 9
Closed: `go test -race` is already included in Phase 2 done criteria in `docs/plan.md`.

### OI-006 — Startup time target wording is broken
Closed: replaced with a `1m` soft watchdog, progress/warning logging, and operator continue/abort prompt.

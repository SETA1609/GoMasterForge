# GoMasterForge — Concept & Plan Review

## Issues to Fix

### Concept § 1.1 — Startup watchdog and operator decision flow

Startup timing was clarified to a soft `1m` watchdog with operator prompt, per-step progress logging, and slow-step warnings. This removes the unrealistic per-mod timing language and keeps slow-start debugging actionable.

### Concept § 3 — Core mod owns registry logic (architecture violation)

> "The `core` mod owns the registry scan and validation logic"

This contradicts the stated boundary: *code behavior in `internal/`, rules in `mods/`*. The `core` mod is a data pack — it should not own Go loader logic. That belongs in `internal/mods/`. Move the ownership claim, or clarify what "owns" means here (data contracts only? Go code?).

### Concept § 6 — Chat ordering via microsecond timestamps is fragile

> "microsecond-resolution timestamps to ensure stable sorted order"

Since all connections are sessions on the *same server process*, a server-assigned monotonic sequence counter is simpler and more reliable than timestamps (which have clock skew risk if distributed components are ever added). Switch to a server-side sequence number for chat ordering.

### Concept § 8 — "Custom text format" for entities/campaign state is undefined

> "Entities and campaign/runtime state use a custom text format that matches the struct contract."

This needs to be either defined (what format?) or changed to JSON/TOML. A custom format implies a custom parser — significant maintenance cost with no stated benefit. If human-editability is the goal, TOML is better. If not, JSON is consistent with the rest of the persistence layer.

### Plan — Race detection deferred to Phase 9

Concurrency is introduced in Phase 2 (`GameState` + subscriptions). Running `go test -race` only in Phase 9 means data races will be found late. Add race testing to the Phase 2 done condition.

### Plan — SSH port never defined in concept

Phase 9 "Immediate Next Steps" mentions port `:2323` — this is the first mention anywhere. The concept has no canonical SSH port or configuration path for it. Add it to § 1.1 or § 2.

---

## Things to Improve

### Concept § 2 — Events vs. diffs open question must be resolved in Phase 0

This directly affects Phase 2 design. Typed domain events are the better choice: lower bandwidth over SSH, give subscribers explicit change context, and compose naturally with the event bus in `internal/events/`. Settle this before Phase 2 starts.

### Concept — GM identity and auth is undefined

There is no rule about how the GM role is assigned. First connection? Config file with a key? This is a security boundary — the entire permission model depends on it. Add a rule to § 5 or a dedicated auth section. The open question about guest players in § 0 is directly related.

### Concept — Session reconnection policy is missing

If a player SSH session drops mid-campaign, what happens? Do they re-join and get state restored? Is their token removed from the map? This needs at minimum an open question entry and ideally a defined rule by v1.

### Plan — Phase 4 references Mods feature that only lands in Phase 5

The "Mods" option in the main menu will be a stub in Phase 4 and only wired in Phase 5. That is fine, but the plan should explicitly note it to avoid confusion during implementation.

### Plan — Phase 8 is too large

All gameplay tabs in one phase is a significant risk. Consider splitting:
- Map + Initiative (core gameplay loop) as its own phase.
- GM-only utility tabs (Notes, Tables, Template Browser) as a follow-on.

Map + Initiative is the highest-value and most complex surface; isolating it gives a natural checkpoint.

### Architecture doc — Package structure does not match the DDD mention

The doc mentions future bounded contexts (`campaign`, `chat`, `mapinitiative`, `profile`) but the current `internal/` has `app`, `events`, `localization`, `mods`, `persistence`, `server`, `state`, `ui`. These are orthogonal organizations. Either commit to the DDD-style structure now (easier to refactor while packages are empty) or remove the DDD reference and call it a technical layers structure.

### Concept § 8 — Debounce/checkpoint open question should just be answered

The question asks "Debounce window and checkpoint interval defaults (e.g. 2s debounce + 5m checkpoint)?" — the suggested values are reasonable. Make them the defaults and document them rather than leaving them open.

---

## Quick Wins

- `go.mod` declares `go 1.26` — current stable is 1.24.x. Verify this is intentional or update it.
- `docs/adrs/` does not exist yet; Phase 0 requires it. Create a stub directory.
- `schemas/` only has a README — the `manifest.json` manifest schema should be the first artifact committed there, before Phase 5 can properly begin.
- `mods/core/manifest.json` and `mods/core-dnd35/manifest.json` are empty. These are referenced as real artifacts in the plan but contain no content.

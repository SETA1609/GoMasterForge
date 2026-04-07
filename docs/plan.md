# GoMasterForge Project Plan

This plan is execution-focused: each phase has clear outputs and a done condition.

Flow contract: users start in a pre-campaign main menu; tabs and split-screen gameplay UI are available only after loading or starting a campaign.

## Phase 0 - Foundations and Decisions

Goal: lock critical decisions before feature growth.

Work:
- Resolve open questions from `docs/concept-v2-0.md` (auth scope, events vs diffs, mod conflicts, save intervals).
- Write ADRs in `docs/adrs/` for each decision.
- Freeze repo conventions: module naming, file naming, package boundaries.
- Freeze container packaging contract: image includes `cmd/`, `internal/`, `mods/`, `schemas/`; runtime state is mounted from `data/`.

Done when:
- All P0 architecture decisions are documented.
- No blocking open questions remain for multiplayer core.

## Phase 1 - Runtime Boot Skeleton

Goal: reliable app boot and clean shutdown.

Work:
- Keep `cmd/gomasterforge/main.go` as canonical entrypoint.
- Implement server bootstrap in `internal/server/` with cancellation support.
- Add structured logging and startup configuration loading.
- Render a minimal pre-campaign main menu shell as the default startup screen.

Done when:
- `go run ./cmd/gomasterforge` starts and exits cleanly.
- Startup logs include version/build info and active config.
- App lands in main menu state before any campaign runtime UI is shown.

## Phase 2 - Shared State and Session Model

Goal: one authoritative multiplayer state.

Work:
- Implement `GameState` with thread-safe mutation APIs in `internal/state/`.
- Add session registration/subscription lifecycle.
- Add broadcast mechanism for state updates.

Done when:
- Two concurrent sessions observe the same state transitions.
- Race checks pass for state and subscription paths.

## Phase 3 - Chat Vertical Slice

Goal: first complete multiplayer feature.

Work:
- Implement shared chat model (append, order, broadcast).
- Add minimal chat UI surface in Bubble Tea flow.
- Add first integration test for cross-session chat propagation.

Done when:
- Message sent from session A appears in session B in order.
- Integration test is green in CI.

## Phase 4 - Main Menu, Profiles, and Permissions

Goal: deliver pre-campaign entry flow and enforce GM/player security boundaries on the server.

Work:
- Implement main menu options: Start Campaign, Continue Campaign, Profile Selection, Mods, Settings.
- Implement profile selection workflows: load, create, export profile.
- Implement pre-campaign settings workflow with persistence to selected profile.
- Implement campaign entry transition: Start/Continue moves to in-campaign runtime.
- Add profile model and role assignment.
- Enforce permissions per action (not only by hidden UI).
- Add authorization tests for forbidden actions.

Done when:
- Main menu features are navigable end-to-end.
- Tabs/split-screen gameplay UI appears only after Start Campaign or Continue Campaign.
- Player attempts to run GM-only actions are rejected server-side.
- Profile data persists and reloads correctly.

## Phase 5 - Mod System MVP

Goal: rules and content are truly data-driven.

Work:
- Define and freeze `manifest.json` schema in `schemas/`.
- Implement loader + dependency validation in `internal/mods/`.
- Load `mods/core` first, then `mods/core-dnd35` with semver checks.

Done when:
- Invalid manifests fail fast with actionable errors.
- Enabled mod set is deterministic and logged at startup.

## Phase 6 - Localization MVP

Goal: i18n is first-class from early stages.

Work:
- Implement i18n loader and fallback chain in `internal/localization/`.
- Wire profile language preference into string resolution.
- Add missing-key warnings without crashing.

Done when:
- Switching profile locale changes visible UI strings.
- Missing keys fallback to `EN` and emit warnings.

## Phase 7 - Persistence and Recovery

Goal: all long-lived data survives restarts safely.

Work:
- Implement file stores in `internal/persistence/` for campaigns/profiles/settings.
- Add debounced save + periodic checkpoint + forced save triggers.
- Add corruption detection and clear load errors.

Done when:
- Restart roundtrip keeps expected campaign and profile state.
- Save/load behavior is covered by automated tests.

## Phase 8 - Gameplay Tabs (v1 Scope)

Goal: deliver core v1 gameplay workflows.

Work:
- Implement Profiles, Settings, Campaigns tabs.
- Implement Map + Initiative shared tab with turn restrictions.
- Implement GM tabs: Entity Manager, Template Browser, Notes, Tables.
- Implement player tabs: Create Character, My Character export.
- Ensure tabs and split layout are in-campaign only (never shown in pre-campaign main menu).
- In in-campaign Settings, implement Network section (connected clients, secrets, related controls) with GM-only server-side permissions.
- In in-campaign Settings, implement Return to Main Menu confirmation with Exit Without Saving and Save and Exit.

Done when:
- All v1 tabs from concept are reachable and permission-safe.
- GM and player smoke tests pass in narrow and wide terminals.
- Save and Exit forces a campaign save before returning to main menu.
- Exit Without Saving returns to main menu without an extra save and behaves per persistence policy.

## Phase 9 - Hardening and Release

Goal: ship a stable v1.

Work:
- Add load tests for target concurrency.
- Run `go test ./...`, `go test -race ./...`, and `go vet ./...` in CI.
- Validate Docker image + mounted `data/` volumes.
- Validate `docker compose up --build` as the default local deployment path.
- Verify persistence roundtrip through compose lifecycle (`up` -> gameplay writes -> `down` -> `up`).
- Final legal pass for OGL and About notice coverage.

Done when:
- CI is green, critical bugs triaged, and release notes prepared.
- Compose startup is reproducible and documented, with data persisting under `./data`.
- First tagged release is published.

## Milestone Order

1. Multiplayer proof: phases 1-3.
2. Security + mod foundation: phases 4-6.
3. Durability + gameplay: phases 7-8.
4. Production readiness: phase 9.

## Current Immediate Next Steps

1. Implement pre-campaign main menu state and routing shell in `internal/app/`.
2. Implement graceful shutdown and config loading in `internal/server/`.
3. Build `GameState` mutation API and subscription hub in `internal/state/`.
4. Add a smoke test checklist for `docker compose up --build` and SSH connectivity on `:2323`.

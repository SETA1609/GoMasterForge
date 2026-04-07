# App Runtime (`internal/`)

This directory contains the Go application runtime for GoMasterForge.

## Layout

- `internal/server/` - SSH server bootstrap and session wiring
- `internal/app/` - top-level Bubble Tea model and routing
- `internal/state/` - shared `GameState`, locking, subscriptions
- `internal/events/` - typed event bus and cross-context dispatch
- `internal/ui/` - layout and tab rendering composition
- `internal/mods/` - mod loader, manifest validation, dependency graph
- `internal/localization/` - language loading and key resolution
- `internal/persistence/` - file-based save/load and migrations

## Architectural Rule

`internal/` owns runtime behavior. Rules, templates, and system content stay in `mods/`. Persistent campaign and profile data stays in `data/`.

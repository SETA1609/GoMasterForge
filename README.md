# GoMasterForge

Pure-terminal VTT for GM-hosted sessions over SSH. Players join with plain `ssh`, share one synchronized game state, and use data-driven mods with localization.

## Repository Layout

- `cmd/gomasterforge/` - executable entrypoint
- `internal/` - application runtime code (server, state, UI, localization, persistence)
- `mods/` - game/rules content packs (source of truth for rules/templates/macros)
- `data/` - runtime persistence (campaigns, profiles, settings, templates, exports)
- `docs/` - concept, implementation plan, architecture notes
- `schemas/` - shared JSON schemas used for validation

## Quick Start

```bash
go run ./cmd/gomasterforge
```

## Docker Compose

Build and run with persisted runtime data mounted from `./data`:

```bash
docker compose up --build
```

Run in detached mode:

```bash
docker compose up --build -d
```

Stop:

```bash
docker compose down
```

## Project Docs

- `docs/concept-v2-0.md`
- `docs/plan.md`
- `docs/architecture.md`
- `docs/txt-format-spec.md`
- `docs/open-issues.md`

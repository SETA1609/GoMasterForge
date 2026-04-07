# GoMasterForge Architecture

This document defines the repository architecture for the data-driven design.

## Top-Level Structure

- `cmd/` - executable entrypoints
- `internal/` - runtime application code
- `mods/` - rulesets/templates/macros/locales as data packs
- `data/` - persistent runtime data (host-mounted in containers)
- `docs/` - concept, plan, ADRs, contributor docs
- `schemas/` - shared JSON schema contracts

## Why This Structure

1. **Clear ownership boundaries**
   - Code behavior lives in `internal/`
   - Rules/content live in `mods/`
   - Runtime persistence lives in `data/`

2. **Data-driven by default**
   - Gameplay systems are configured by mod data, not hardcoded per ruleset.

3. **Operationally friendly**
   - `data/` is easy to mount, back up, inspect, and restore.

4. **DDD-friendly growth**
   - `internal/*` can evolve by bounded context (`campaign`, `chat`, `mapinitiative`, `profile`) while preserving top-level boundaries.

## Runtime Rule

Runtime code in `internal/` must not become the source of truth for game system rules. Rules are loaded from `mods/` and merged with campaign overrides from `data/`.

## Contracts and Core Mod

- Go interfaces and reusable services belong in `internal/` (runtime code).
- `mods/core` is a shared data-contract pack, not a service layer.
- `mods/core` should contain shared schemas, translation keys, canonical template categories, and other cross-mod data definitions.
- Registry scanning, dependency validation, and merge logic belong in `internal/mods/`.

## Current Entry Point

- `cmd/gomasterforge/main.go`

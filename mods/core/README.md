# `mods/core`

Shared base data contracts for all game mods.

## Purpose

`core` is not executable runtime logic. It provides cross-mod definitions that the runtime in `internal/` can validate and consume.

## Should Contain

- shared JSON schemas used by multiple mods
- shared template categories/taxonomy
- shared translation key namespaces
- compatibility metadata and legal metadata templates
- shared effect contracts in `data/effects/` for reuse by ruleset and expansion mods

## Should Not Contain

- Go interfaces
- Go services or registry/loader code
- ruleset-specific gameplay content

## Ownership Boundary

- Runtime behavior ownership: `internal/mods/`, `internal/state/`, and other `internal/*` packages
- Data contract ownership: `mods/core`

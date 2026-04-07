# `mods/core-dnd35`

SRD/OGL-safe D&D 3.5 content pack that depends on `mods/core`.

## Purpose

Provide templates, macros, tables, and localization data for the D&D 3.5 system using shared contracts defined by `mods/core`.

## Data Layout

- `data/ogl/` - OGL license text, copyright notice, and source provenance
- `data/indexes/` - cross-domain content index for loader/bootstrap
- `data/*` domains - SRD-derived records (races, classes, feats, spells, monsters, items, rules, tables)

## Dependency

- requires `core`
- validated by runtime loader in `internal/mods/`

## Boundary

This mod contains data only. Runtime behavior is implemented in `internal/*` packages.

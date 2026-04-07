# Mod Manifest Templates

Use these templates when creating new mod packages.

- Layout reference: `mod-layout.txt`

## Ruleset Mod

- File: `ruleset-manifest.template.json`
- Purpose: core rules package that provides runtime rules behavior plus base data.
- Should include: `entrypoints.data = ["data/"]`.
- Must include: `runtime.rules_engine`.

## Expansion Mod

- File: `expansion-manifest.template.json`
- Purpose: data-only content for one ruleset (races, entities, items, tables, lore).
- Should include: `entrypoints.data = ["data/"]`.
- Must include: `ruleset`.
- Must not include: `runtime.rules_engine`.

## Type Contract

- `type = ruleset`: implementation-backed core ruleset.
- `type = expansion`: data-only content pack attached to a ruleset.

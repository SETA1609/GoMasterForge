# Core Effects Catalog

This directory defines reusable effect contracts shared across ruleset and expansion mods.

- Effects are metadata contracts, not executable Go code.
- Ruleset engines map `effect_id` values to implementation handlers.
- Content packs (spells, feats, items, abilities) should reference these IDs instead of duplicating effect logic.

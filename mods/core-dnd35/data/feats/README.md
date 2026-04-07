# Feat Model Notes

This mod uses parameterized feat definitions to avoid one feat record per weapon or school.

## Parameterized Feats

- `dnd35.feat.weapon_focus` accepts `weapon_ref`.
- `dnd35.feat.spell_focus` accepts `school`.

Consumers attach chosen parameters at selection time:

```json
{
  "id": "dnd35.feat.weapon_focus",
  "params": {
    "weapon_ref": "dnd35.item.longsword"
  }
}
```

Runtime must validate parameters against `parameter_schema` before applying effects.

The current catalog includes all feats from the SRD index (`110` records). Many feats still use placeholder `prerequisites`/`effects` and should be enriched from SRD feat descriptions.

Current enrichment status:

- all feats include `srd_anchor` and `srd_benefit_summary`
- prerequisite parsing is normalized where possible (ability, BAB, class level, feat references)
- generic and parameterized feat families are modeled (`weapon_*`, `spell_*`, `skill_focus`, metamagic)
- feat-to-feat prerequisite references are used for dependency chains (for example `spirited_charge` -> `mounted_combat` + `ride_by_attack`, `two_weapon_defense` -> `two_weapon_fighting`)

Update:

- no feat records have empty `effects`
- non-numeric/action feats are represented through `core.grant_tag` capability tags
- no remaining `raw` prerequisite tokens

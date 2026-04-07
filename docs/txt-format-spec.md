# GoMasterForge TXT Format Spec (v1)

This document defines the custom `.txt` format used by GoMasterForge for non-localization data.

Localization remains JSON only.

## Scope

Applies to:

- `data/campaigns/*.txt`
- `data/profiles/*.txt`
- `data/settings/*.txt`
- `mods/*/templates/*.txt`

`manifests.json` is JSON and out of scope for this TXT format spec.

## File Contract

- File extension is `.txt`.
- First metadata block must include `format` and `version`.
- Parser must fail fast if either key is missing.

Example header:

```txt
format = gmf-campaign
version = 1
```

## Syntax

- Comments start with `#` and run to end-of-line.
- Empty lines are ignored.
- Key-value line: `key = value`.
- Section: `[section-name]`.
- Repeated section: `[[section-name]]`.
- Keys are lowercase snake_case.

## Value Rules

- Strings: raw text after `=` (trim outer spaces).
- Bool: `true` or `false`.
- Int: `-?[0-9]+`.
- Float: `-?[0-9]+(\.[0-9]+)?`.
- Duration: `Ns`, `Nm`, or `Nh` (examples: `2s`, `5m`, `1h`).
- CSV list: `a,b,c` (split by comma, trim each item).

## Supported Formats

### `format = gmf-campaign`

Required top-level keys:

- `campaign_id`
- `campaign_name`
- `active_profile`
- `active_mods` (CSV `mod@version`)

Optional sections:

- `[initiative]`
- `[map]`
- `[[token]]`

### `format = gmf-profile`

Required top-level keys:

- `profile_id`
- `display_name`
- `role` (`gm` or `player`)
- `locale`

Optional sections:

- `[settings]`
- `[network]`

### `format = gmf-settings`

Required top-level keys:

- `ssh_host`
- `ssh_port`
- `default_locale`

Optional sections:

- `[security]`
- `[save_policy]`

### `format = gmf-dndtxt`

Required top-level keys:

- `kind` (`entities` or `items`)

Required repeated sections by `kind`:

- `kind = entities` -> `[[entity]]`
- `kind = items` -> `[[item]]`

## Versioning

- Current supported version: `version = 1`.
- Any other version must fail with a clear migration/error message.

## Parser Behavior

- Parse order is top-to-bottom.
- Later duplicate keys in the same scope overwrite earlier keys.
- Repeated sections append records in file order.
- Unknown sections or keys produce warnings by default.
- Strict mode may treat unknown keys/sections as errors.

## Validation Errors (Recommended Codes)

- `E100` unknown format
- `E101` missing required metadata key (`format` or `version`)
- `E102` unsupported version
- `E110` malformed section header
- `E111` invalid key syntax
- `E120` invalid value type
- `E130` missing required domain key
- `E140` invalid enum value
- `E150` invalid dependency/version token in CSV fields

Example error output:

```txt
E130 data/campaigns/example-campaign.txt:7 missing required key: active_mods
```

## Non-Goals (v1)

- Nested objects beyond one section level.
- Multi-line string literals.
- Include/import directives.

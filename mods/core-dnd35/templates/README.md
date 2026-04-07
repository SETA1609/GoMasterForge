# Templates (`.txt` with format flag)

This directory contains D&D 3.5 template data in plain `.txt` files using the custom GoMasterForge format.

Structured SRD/OGL datasets now live in `../data/` as JSON domain files.

## Files

- `entities.txt` - entity templates (PC/NPC/monster)
- `items.txt` - item templates

## Format Notes

- Header fields define `format`, `version`, and `kind`.
- Repeated blocks use `[[entity]]` or `[[item]]`.
- Values are key-value pairs and should stay ASCII.

Parsers must validate the header `format` value and fail fast if it is missing or unexpected.

Example parser strategy: split into sections by `[[...]]`, parse each section into maps, then validate with `mods/core` schema contracts.

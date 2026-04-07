# Schemas

Shared JSON schema contracts that are not tied to a single mod package.

Use this directory for:

- `manifest.json` manifest schema
- shared data contract schemas used by app validation
- versioned schema migrations and compatibility notes

Mod-specific schemas can still live in each mod under `mods/<mod-id>/schemas/` when ownership is local to that mod.

Custom `.txt` parser and validation rules are specified in `docs/txt-format-spec.md`.

## Related Example Files

- root `manifests.json` (manifest index)
- `data/campaigns/*.txt` (campaign state, with `format` header)
- `data/profiles/*.txt` (profile state, with `format` header)
- `data/settings/*.txt` (server/runtime settings, with `format` header)

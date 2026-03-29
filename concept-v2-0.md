# dnd-tui — Concept v2.0

**Project name:** `dnd-tui`
**Tech stack:** Go, Bubble Tea (Elm architecture), Wish (SSH server), lipgloss, bubbles, file-based persistence (v1), in-memory registries/indexes for runtime lookups, i18next-style JSON localization.

---

## 0) Product Vision

Build a pure-terminal virtual tabletop where one GM hosts a server and players join with plain `ssh`. No client install is ever required. The system is modular, data-driven, localization-first, and designed to support multiple tabletop rulesets as drop-in mods.

**Rules:**
- GM runs one process; players connect via standard `ssh` from any terminal.
- Core code is rules-agnostic; all game logic and content lives in mods.
- Full UI and mod-text localization.
- All campaign and profile data survive restarts.

**Open questions:**
- Should v1 target only LAN/self-host, or include public internet hardening by default?
- Should guest (unauthenticated) players be supported in v1?

---

## 1) v1 Scope

**In scope for v1:**
- SSH server, shared game state, chat, Profiles tab, Settings tab, Map+Initiative tab.
- Mod loader with `core` mod and one game mod (`core-dnd35`, SRD/OGL-safe content only).
- Role-based tab and action permissions enforced server-side.
- i18next-style localization with per-profile language preference.
- Persistence for campaigns, profiles, settings, templates, notes, exports.

**Out of scope for v1** (post-v1 extras):
- Advanced fog of war.
- PDF export.
- Optional sound hooks.
- Searchable rules reference.
- GM mute/ghost player tools.
- Deep roll automation beyond basic macros.

**Rules:**
- Any extra feature must not block core multiplayer stability.

---

## 1.1) Operational Targets

**Intent:** Define baseline resource and startup expectations for deployment planning.

**Rules:**
- Target runtime memory budget: app and container should operate within `512 MB` RAM.
- Startup duration is variable based on enabled mod count and mod size.
- Startup estimate target: `30s` minimum to `60s` maximum per enabled mod in the mod list.
- Startup should log per-mod load timing to help diagnose slow mods.

**Open questions:**
- Should startup time limits be soft warnings or hard-fail thresholds?

---

## 2) Runtime Architecture

**Intent:** One shared authoritative game state with per-session UI models.

**Rules:**
- One process hosts Wish SSH and all active sessions.
- One shared `GameState` is the single authoritative source for live campaign state.
- Each SSH connection gets its own Bubble Tea session model (strict Elm update loop).
- Session models mutate shared state only through explicit state/service APIs — never directly.
- State changes are broadcast to all subscribed sessions.
- Concurrency safety is mandatory: lock discipline and race-tested paths.

**Open questions:**
- Should broadcast emit whole-state diffs or typed domain events only?

---

## 3) Mod System

**Intent:** Treat mods as the sole source of rules, templates, macros, and system-specific behavior.

**Rules:**
- Mods live in `./mods/<mod-name>/` and are auto-discovered at startup.
- The `core` mod loads first; dependent mods load only after dependency validation.
- Each mod provides a `mod.json` manifest with `id`, `version`, `dependencies`, and compatibility metadata.
- Templates (entities, items, spells, feats, etc.) are data files — not hardcoded in app logic.
- At bootstrap, loaded templates are compiled into in-memory registries (maps/indexes) for fast lookups.
- Runtime queries for gameplay/macros/UI read from in-memory registries, not by reparsing files.
- Campaign-level custom templates and overrides live in `data/templates/` and are merged after base mods.
- Invalid manifests or failed dependency checks fail fast with clear startup errors.
- The `core` mod owns the registry scan and validation logic; `internal/mods/` is the loader that invokes it.
- GM can manage an enable/disable mod list from the app.
- GM can export/import mod selection profiles for reuse across campaigns.
- A mod cannot be enabled if any declared dependency mod is missing or disabled.

**Open questions:**
- Strict semver compatibility checks or a custom compatibility window for mod versions?
- Should mod conflicts fail hard or allow override precedence with warnings?
- Do we need a mod linter that blocks load when OGL legal metadata is missing?

---

## 4) UI Layout

**Intent:** Adaptive terminal UI that works on desktop and narrow/mobile clients.

**Rules:**
- Default layout: main tab area takes 4/5 of terminal width; chat side panel takes 1/5 (right).
- On narrow terminals, chat collapses into a dedicated tab instead of the side panel.
- Layout mode switches deterministically based on a terminal-width threshold.
- All top-level tabs remain keyboard-navigable in both layout modes.
- No gameplay-critical action may be accessible only in one layout mode.

**Open questions:**
- Exact width breakpoint for sidebar-to-tab switch (e.g. 120 columns)?

---

## 5) Tabs and Permissions

**Intent:** Strict GM/player permissions enforced at the server, not just in the UI.

**Tab inventory:**

| Tab | Visible to | Notes |
|-----|-----------|-------|
| Profiles | Everyone | Players see only their own; GM sees and manages all |
| Settings | Everyone | Connection status, localization, export options, About/legal |
| Entity Manager | GM only | Create/edit live entities; search and grant items to any entity |
| Bestiary / Template Browser | GM only | Filters by game, version, category, tags, CR; preview pane |
| Campaigns | GM only | Create, load, save, list, and manage campaigns |
| Notes (Journal) | GM only | Timestamped append-only journal; line-level delete via vi-style single-key commands |
| Tables | GM only | Loot tables, random encounters, weather, and other mod-defined tables |
| Map + Initiative | Everyone | GM controls all tokens and can lock turns; players move only their own token on their turn |
| Create Character | Player only | Select template, customize, and save to profile |
| My Character | Player only | Equipment, level-ups, spellbook; export character (Markdown / plain text / JSON) |

**Rules:**
- Permissions are enforced per action, not only by tab visibility.
- Server-side permission checks are mandatory for all state mutations.
- Character export is accessible from both **My Character** and the **Settings** tab.

**Open questions:**
- Should the GM be able to impersonate player actions for troubleshooting?

---

## 6) Chat and Macros

**Intent:** Chat is both the communication surface and the command surface, powered by mod-defined macros and a central event bus.

**Rules:**
- Chat is shared live across all connected sessions.
- Chat supports plain messages and slash-style macros (e.g. `/attack goblin-sword`, `/grant-item playername sword`).
- Macros are defined entirely by active mod data and executed by server-side handlers.
- Macro execution passes permission checks and input validation before running.
- Macro outcomes emit typed events consumed by the relevant subsystems (initiative, map, inventory).
- All results are broadcast to every session with timestamp and actor context.
- Chat viewport supports scrolling, timestamps, and per-user color coding.
- Chat is not part of durable campaign persistence.
- If chat history is persisted, it is written only as a session log file for audit/recap purposes.

**Open questions:**
- Should macro execution support transactional rollback on partial failures?
- Is per-user rate limiting needed in v1 to prevent spam/abuse?

---

## 7) Localization

**Intent:** Every user-facing string is translatable with reliable fallback behavior.

**Rules:**
- Translation files use i18next-style JSON format.
- Default language is English (`en`).
- Resolution order: profile language → server default → `en` fallback.
- Missing keys render a safe fallback string and emit a log warning (no crash).
- Both app UI strings and mod content use translation keys.
- Language preference persists per profile.

**Open questions:**
- Should translation key validation run at startup and hard-fail on missing required keys?

---

## 8) Persistence and Storage

**Intent:** All long-lived data is durable, inspectable, and container-volume friendly.

**Rules:**
- All persistence roots live under `./data/` subdirectories:
  `campaigns/`, `profiles/`, `settings/`, `templates/`, `translations/`, `mod-cache/`, `notes/`, `exports/`.
- v1 persistence is file-based (JSON/Markdown/plain text as appropriate) — no SQLite dependency.
- Campaign state includes entities, initiative order, map data, and all session-independent state.
- Chat history is excluded from campaign state persistence.
- Profiles, settings, custom templates, notes, and exports are stored separately.
- Save/load is deterministic and versioned for future migrations.
- Data integrity checks run on load; corrupted files fail with actionable errors.
- Container deployment uses host-mounted volumes for all `data/` directories.
- `mod-cache/` is derived data and can be rebuilt — it is not a source of truth.
- Authoritative runtime state is in memory; disk writes follow a hybrid save policy.
- Campaign save metadata must include required mod ids and compatible version constraints.
- A campaign cannot be loaded when required mods are missing, disabled, or incompatible.

**Save policy:** Debounced writes after mutations (short window) + periodic checkpoint autosave + forced save on explicit save, campaign load/switch, and clean shutdown. Avoid write-on-every-mutation to reduce disk churn and contention.

**Open questions:**
- Debounce window and checkpoint interval defaults (e.g. 2s debounce + 5m checkpoint)?
- Snapshot retention policy and rotation count for v1?

---

## 9) Canonical Data Ownership

One source of truth per domain — no ambiguity about where to read or write:

| Domain | Authoritative source |
|--------|---------------------|
| Live gameplay state | `GameState` (in-memory, broadcast on change) |
| Base rules and templates | Loaded mods (compiled into in-memory registries at startup) |
| Campaign-specific overrides | `data/campaigns/` and `data/templates/` |
| Profile preferences and locale | `data/profiles/` |
| Localization strings | App and mod translation catalogs |
| Derived caches | `data/mod-cache/` (rebuildable, not authoritative) |

**Open questions:**
- Should campaign overrides be stored as full copies or as patch deltas against mod templates?
- Which package owns cross-domain IDs and shared type definitions — `core` mod schema or a framework internal model?

---

## 10) Legal and Licensing

**Rules:**
- App source code is MIT-licensed.
- OGL-based game mods include `LICENSE-OGL.txt` and clearly separate Open Game Content from Product Identity.
- `core` mod contains only generic interfaces and shared structures.
- `core-dnd35` includes SRD/Open Game Content only — no Product Identity, trademarks, or non-SRD copyrighted material.
- Each game mod's `mod.json` documents how OGL compliance is handled.
- README and in-app **Settings → About** include a legal notice.
- Do not claim official compatibility or use Wizards of the Coast / Paizo trademarks.

**Risk level:** Extremely low for a non-commercial, open-source hobby project that follows these rules. Thousands of SRD-based tools exist without issue. Risk increases if you ship full Player's Handbook content, monetize the app, or distribute non-SRD material.

**Open questions:**
- Should a mod linter block load/publish when required OGL legal metadata is missing from the manifest?

---

## 11) Post-v1 Extras (Backlog)

Desirable but explicitly out of scope for v1:

- Fog of war on the map (text-based unexplored-area hiding).
- Campaign export to Markdown or PDF for post-session recaps.
- Searchable rules reference pulled from the active mod.
- Optional terminal sound hooks (critical hits, important events).
- GM tools to mute or ghost individual players.
- Built-in dice roller accessible from any tab or chat (may be promoted to v1 if trivial to add).

---

## 12) Initial Execution Checklist

First concrete actions to move from concept to code:

1. Finalize `mod.json` manifest schema and dependency rules.
2. Implement server boot path and per-session Bubble Tea model wiring.
3. Implement shared `GameState` + broadcaster + basic chat.
4. Implement role model and server-side action guards.
5. Implement localization loader with fallback chain.
6. Implement persistence bootstrap for profiles, settings, and campaigns.
7. Write first integration test: two SSH sessions receive one shared chat update.

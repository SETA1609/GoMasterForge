# GoMasterForge — Concept v2.0

**Project name:** `GoMasterForge`
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
- Main menu flow before campaign load: Start Campaign, Continue Campaign, Profile Selection, Mods, Settings.
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
- `512 MB` RAM is the low-end sizing target/estimate, not a hard runtime cap.
- Startup duration is variable based on enabled mod count and mod size.
- Startup should log per-mod load timing to help diagnose slow mods.
- Startup has a soft watchdog at `1m` total elapsed time.
- If startup reaches `1m`, the server pauses for operator input, shows the current startup step (for example mod registration or data scan), and asks whether to continue or abort.
- If an individual mod registration or data load step exceeds expected duration, emit a warning with step name and elapsed time to aid debugging.
- Slow startup does not hard-fail automatically; operator may continue when the delay is expected (for example very large mod data sets).
- Soft concurrency target: GM + up to 8 concurrent players.
- Connections above 8 players are allowed on a best-effort basis (warn, do not hard-block).

---

## 1.2) Main Menu and Session Entry Flow

**Intent:** Ensure startup UX is predictable before entering live campaign state.

**Connect-time profile gate (before main menu):**
- On SSH connect, the session is immediately presented with a profile screen — not the main menu.
- Profile screen options: **Load Profile**, **Import Profile**, and **Create Profile**.
- **Load Profile** lists profiles already present in `data/profiles/` (the host-mounted volume); user selects one by name.
- **Import Profile** accepts a profile exported from another instance (e.g. a `.txt` file the user pastes or transfers); it is saved into `data/profiles/` on import.
- **Create Profile** collects: display name, role (`gm` or `player`), and locale preference; saved to `data/profiles/` on creation.
- A profile must be active before the main menu is shown.
- If `data/profiles/` is empty, **Load Profile** is hidden and only **Import Profile** and **Create Profile** are offered.

**Rules:**
- Main menu is shown only after a profile is active for the session.
- Main menu options are: **Start Campaign**, **Continue Campaign**, **Mods**, **Settings**.
- **Profile Selection** is removed from the main menu; profile is handled at connect time. A **Switch Profile** option may be offered from **Settings** to re-enter the profile gate.
- **Mods** supports mod list view, mod-list export, and enable/disable actions.
- **Settings** in main menu controls pre-campaign global preferences and writes them to the selected profile.
- Entering **Start Campaign** or **Continue Campaign** transitions from pre-campaign mode to in-campaign mode.
- Tabs and split-screen layout are part of in-campaign runtime only; they are not shown in the main menu flow.
- After a campaign is loaded/started, **Settings** must include a **Network** section for GM operations (view connected clients, manage connection secrets, and related session controls).
- In-campaign **Settings** must include **Return to Main Menu** with a confirmation prompt offering **Exit Without Saving** and **Save and Exit**.

**Open questions:**
- Should **Start Campaign** and **Mods** be restricted to GM-role profiles only, or visible to all?

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
- Each mod provides a `manifest.json` manifest with `id`, `version`, `dependencies`, and compatibility metadata.
- Templates (entities, items, spells, feats, etc.) are data files — not hardcoded in app logic.
- At bootstrap, loaded templates are compiled into in-memory registries (maps/indexes) for fast lookups.
- Runtime queries for gameplay/macros/UI read from in-memory registries, not by reparsing files.
- Campaign-level custom templates and overrides live in `data/templates/` and are merged after base mods.
- Invalid manifests or failed dependency checks fail fast with clear startup errors.
- Registry scan, dependency resolution, and validation logic are owned by `internal/mods/` runtime code.
- The `core` mod provides shared data contracts (schemas, translation keys, base template taxonomy) consumed by the runtime and dependent mods.
- Mods are data packs only and do not own executable Go services.
- GM can manage an enable/disable mod list from the app.
- GM can export/import mod selection profiles for reuse across campaigns.
- A mod cannot be enabled if any declared dependency mod is missing or disabled.
- Mod compatibility checks use semantic versioning (`semver`) from manifest versions.
- Downgrade compatibility rule: only major downgrade is blocked (`2.x.x` -> `1.x.x` blocked; `2.x.x` -> lower `2.y.z` allowed).

**Open questions:**
- Should mod conflicts fail hard or allow override precedence with warnings?
- Do we need a mod linter that blocks load when OGL legal metadata is missing?

---

## 4) UI Layout

**Intent:** Adaptive in-campaign terminal UI that works on desktop and narrow/mobile clients.

**Rules:**
- Layout rules in this section apply only after **Start Campaign** or **Continue Campaign**.
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

**Scope rule:** The tab system is part of campaign runtime and becomes available only after a campaign is started or loaded.

**Tab inventory:**

| Tab | Visible to | Notes |
|-----|-----------|-------|
| Profiles | Everyone | Players see only their own; GM sees and manages all |
| Settings | Everyone | Pre-campaign global settings (saved in profile); in-campaign includes connection status/localization/export/About, GM network controls, and Return to Main Menu flow |
| Entity Manager | GM only | Create/edit live entities; search and grant items to any entity |
| Bestiary / Template Browser | GM only | Filters by game, version, category, tags, CR; preview pane |
| Campaigns | GM only | Create, load, save, list, and manage campaigns |
| Notes (Journal) | GM only | Timestamped append-only journal; line-level delete via vi-style single-key commands |
| Tables | GM only | Loot tables, random encounters, weather, and other mod-defined tables |
| Map + Initiative | Everyone | GM controls all tokens and can lock turns; players move only their own token on their turn |
| Create Character | Player only | Select template, customize, and save to profile |
| My Character | Player only | Equipment, level-ups, spellbook; export character (Markdown / plain text / JSON) |

**Rules:**
- GM role is granted to the session that starts or continues a campaign using a profile with `role = gm`.
- Only one GM session is active per campaign at a time.
- Players join with profiles where `role = player`; their role is read from their profile on connect.
- Permissions are enforced per action, not only by tab visibility.
- Server-side permission checks are mandatory for all state mutations.
- Character export is accessible from both **My Character** and the **Settings** tab.
- Network controls in **Settings** are server-authoritative and restricted to GM-level permissions.
- In-campaign **Return to Main Menu** is server-authoritative and restricted to GM-level permissions.

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
- Outgoing chat messages are appended to one shared chat queue/data structure.
- Chat ordering is determined by microsecond-resolution timestamps to ensure stable sorted order.

**Open questions:**
- Should macro execution support transactional rollback on partial failures?
- Is per-user rate limiting needed in v1 to prevent spam/abuse?

---

## 7) Localization

**Intent:** Every user-facing string is translatable with reliable fallback behavior.

**Rules:**
- Translation files use i18next-style JSON format.
- Locale keys use language-only uppercase codes (for example `EN`, `ES`, `DE`); region variants such as `en-US` are out of scope for v1.
- Default language is English (`EN`).
- Resolution order: profile language → server default → `EN` fallback.
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
- v1 persistence is file-based with no SQLite dependency.
- Mod manifests and localization catalogs use JSON.
- Entities and campaign/runtime state use a custom text format that matches the struct contract.
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
- On campaign load mismatch, the app must return actionable errors listing missing dependencies and version mismatches.

**Save policy:** Debounced writes after mutations (short window) + periodic checkpoint autosave + forced save on explicit save, campaign load/switch, and clean shutdown. Avoid write-on-every-mutation to reduce disk churn and contention.

**Campaign exit policy (from in-campaign Settings):** **Save and Exit** performs a forced campaign save before returning to main menu. **Exit Without Saving** returns to main menu without an additional save and may discard in-memory changes since the last successful save/checkpoint.

**Migration policy:** If loaded campaign metadata requires migration for the active mod set, migration is applied as an override write after the first autosave checkpoint.

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
- Which package owns cross-domain IDs and shared type definitions — `core` mod schema or an app internal model?

---

## 10) Legal and Licensing

**Rules:**
- App source code is MIT-licensed.
- OGL-based game mods include `LICENSE-OGL.txt` and clearly separate Open Game Content from Product Identity.
- `core` mod contains only generic data contracts and shared definitions.
- `core-dnd35` includes SRD/Open Game Content only — no Product Identity, trademarks, or non-SRD copyrighted material.
- Each game mod's `manifest.json` documents how OGL compliance is handled.
- README and in-app **Settings → About** include a legal notice.
- README must clearly warn that automatic migration updates campaign data and that major-version mod downgrades are not supported.
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

1. Finalize `manifest.json` manifest schema and dependency rules.
2. Implement server boot path and per-session Bubble Tea model wiring.
3. Implement shared `GameState` + broadcaster + basic chat.
4. Implement role model and server-side action guards.
5. Implement localization loader with fallback chain.
6. Implement persistence bootstrap for profiles, settings, and campaigns.
7. Write first integration test: two SSH sessions receive one shared chat update.

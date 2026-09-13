# subtext-verify

## 0.3.0

### Minor Changes

- 285ac83: subtext-verify-shared: upload the `.sightmap/` corpus with the `sightmap` CLI
  (`sightmap export --url <url>`) instead of the bundled Python collector, and remove
  `collect_and_upload_sightmap.py`.

  Matches the Subtext change: `sightmap export` routes the upload through the Go
  loader and POSTs the whole canonical wire (components, views/routes, requests,
  messages, memory, tags), dropping the Python 3 / PyYAML dependency. Applies to the
  `live-connect` / `live-tunnel` upload URLs; the `subtext-tunnel` setup note is
  updated to require the `sightmap` binary rather than PyYAML.

- 9f81179: Reframe Subtext Verify as a companion plugin to Subtext and remove skill-name
  collisions and duplicated guidance. Verify stays a plugin (manifests, marketplace
  install, SessionStart proof-trigger hook), but now explicitly requires Subtext
  installed first and no longer re-ships anything Subtext already owns.

  - Namespaced the skills that collided with Subtext: `subtext-shared` →
    `subtext-verify-shared` (slimmed to verify-only deltas; general conventions
    defer to Subtext's `subtext-shared`), and `subtext-using-subtext` →
    `subtext-verify-using` (the SessionStart hook now injects this).
  - Removed the `subtext-sightmap` schema fork (it duplicated Subtext's canonical
    `sightmap-authoring`); `subtext-recipe-sightmap-setup` and `subtext-onboard`
    now reference `sightmap-authoring`.
  - Removed the separate `subtext-setup-plugin` skill; its tunnel-server
    connectivity and PyYAML-dependency notes are folded into `subtext-tunnel`.
  - MCP wiring is now tunnel-only: Verify declares just the `subtext-tunnel` stdio
    server and relies on base Subtext for the shared `subtext` HTTP server.
  - Absorbed the sightmap request-enrichment note into `subtext-live`.

  Every skill name is now unique across a combined Subtext + Verify install, so the
  two compose without collisions.

## 0.2.0

### Minor Changes

- Enrich the marketplace manifests with Subtext branding. The Codex manifest gains an `interface` block — display name, short/long descriptions, brand color (#F5447B), example prompts, capabilities, legal links, and bundled composer/logo icons. Homepage, repository, and keywords are added across the Claude, Codex, and Cursor manifests, and author identity is standardized to Subtext (subtext@fullstory.com, https://subtext.fullstory.com).
- 06f22e7: Prefix every skill folder with `subtext-` (e.g. `subtext-proof`, `subtext-live`, `subtext-comments`, `subtext-docs`, `subtext-tunnel`, `subtext-sightmap`, `subtext-onboard`, `subtext-first-session`, `subtext-recipe-sightmap-setup`, `subtext-shared`, `subtext-using-subtext`, `subtext-setup-plugin`).

  The namespace now lives in the skill folder name itself, so skills stay collision-free across the harnesses that don't namespace plugins (Cursor, `.agents/skills`, `npx openskills`). Skill invocation names change accordingly. Cross-references in skill bodies and the README, the `SessionStart` hook's skill path, and the sightmap-upload helper path were all updated to match.

  Note: the renamed `subtext-tunnel` skill now shares its name with the existing `subtext-tunnel` MCP server. They live in separate namespaces (skill vs. MCP server) so there's no functional clash, but it's worth being aware of.

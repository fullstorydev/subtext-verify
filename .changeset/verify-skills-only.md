---
"subtext-verify": minor
---

Reframe Subtext Verify as a companion plugin to Subtext and remove skill-name
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

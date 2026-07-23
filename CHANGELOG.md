# subtext-verify

## 0.2.0

### Minor Changes

- Enrich the marketplace manifests with Subtext branding. The Codex manifest gains an `interface` block — display name, short/long descriptions, brand color (#F5447B), example prompts, capabilities, legal links, and bundled composer/logo icons. Homepage, repository, and keywords are added across the Claude, Codex, and Cursor manifests, and author identity is standardized to Subtext (subtext@fullstory.com, https://subtext.fullstory.com).
- 06f22e7: Prefix every skill folder with `subtext-` (e.g. `subtext-proof`, `subtext-live`, `subtext-comments`, `subtext-docs`, `subtext-tunnel`, `subtext-sightmap`, `subtext-onboard`, `subtext-first-session`, `subtext-recipe-sightmap-setup`, `subtext-shared`, `subtext-using-subtext`, `subtext-setup-plugin`).

  The namespace now lives in the skill folder name itself, so skills stay collision-free across the harnesses that don't namespace plugins (Cursor, `.agents/skills`, `npx openskills`). Skill invocation names change accordingly. Cross-references in skill bodies and the README, the `SessionStart` hook's skill path, and the sightmap-upload helper path were all updated to match.

  Note: the renamed `subtext-tunnel` skill now shares its name with the existing `subtext-tunnel` MCP server. They live in separate namespaces (skill vs. MCP server) so there's no functional clash, but it's worth being aware of.

# Subtext Verify

Drive a hosted browser, capture before/after proof of code changes, and leave
reviewer-facing evidence (screenshots, viewer links, comments) against your
running app — for coding agents.

This plugin bundles:

- **Skills** — `subtext-proof` (before/after visual evidence, the inner loop), `subtext-onboard`
  and `subtext-first-session` (first-run flows), plus the tool catalogs `subtext-live`,
  `subtext-comments`, `subtext-docs`, `subtext-tunnel`, `subtext-sightmap`, the `subtext-recipe-sightmap-setup` recipe,
  and `subtext-shared` / `subtext-using-subtext` / `subtext-setup-plugin`.
- **MCP servers**
  - `subtext` (HTTP) — `https://api.fullstory.com/mcp/subtext` (EU1 mirror:
    `https://api.eu1.fullstory.com/mcp/subtext`). Serves the `live-*`,
    `comment-*`, `doc-*`, `artifact-*`, and `clip-*` tools.
  - `subtext-tunnel` (stdio) — the local reverse-tunnel client, run via
    `npx -y @subtextdev/subtext-cli@latest tunnel mcp`, published from the
    [`fullstorydev/subtext-cli`](https://github.com/fullstorydev/subtext-cli) repo.
- **Hooks** — a `SessionStart` hook injects the `subtext-using-subtext` router so the
  agent reaches for `subtext-proof` before editing UI code.

## Install

**Claude Code**
```
/plugin marketplace add fullstorydev/subtext-review
/plugin install subtext-verify@subtext-marketplace
```

**Cursor** — install from the Marketplace panel (or a Team Marketplace that
imports this repo).

**Codex** — open `/plugins`, install **subtext-verify** from the repo marketplace.

**Gemini CLI**
```
gemini extensions install https://github.com/fullstorydev/subtext-verify
```

## Notes

- The `subtext` HTTP tool surface is gated server-side by org/feature flag; this
  plugin assumes the verification flag is enabled for your org.
- Read-only review of completed sessions and privacy-rule management live in the
  separate **Subtext Review** plugin.

# Subtext Verify

Drive a hosted browser, capture before/after proof of code changes, and leave
reviewer-facing evidence (screenshots, viewer links, comments) against your
running app — for coding agents.

This plugin bundles:

- **Skills** — `proof` (before/after visual evidence, the inner loop), `onboard`
  and `first-session` (first-run flows), plus the tool catalogs `live`,
  `comments`, `docs`, `tunnel`, `sightmap`, the `recipe-sightmap-setup` recipe,
  and `shared` / `using-subtext` / `setup-plugin`.
- **MCP servers**
  - `subtext` (HTTP) — `https://api.fullstory.com/mcp/subtext` (EU1 mirror:
    `https://api.eu1.fullstory.com/mcp/subtext`). Serves the `live-*`,
    `comment-*`, `doc-*`, `artifact-*`, and `clip-*` tools.
  - `subtext-tunnel` (stdio) — the local reverse-tunnel client, run via
    `npx -y @subtextdev/subtext-cli@latest tunnel mcp`. Source lives in
    [`cli/`](cli/).
- **Hooks** — a `SessionStart` hook injects the `using-subtext` router so the
  agent reaches for `proof` before editing UI code.

## Install

**Claude Code**
```
/plugin marketplace add fullstorydev/subtext-verify
/plugin install subtext-verify@subtext-verify-marketplace
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

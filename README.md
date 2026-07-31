# Subtext Verify

Drive a hosted browser, capture before/after proof of code changes, and leave
reviewer-facing evidence (screenshots, viewer links, comments) against your
running app — for coding agents.

Subtext Verify is a **companion plugin** to [Subtext](https://github.com/fullstorydev/subtext).
It layers a live-browser proof workflow on top of a normal Subtext install — it
is not a standalone product.

## Prerequisite: install Subtext first

Install [Subtext](https://github.com/fullstorydev/subtext) the usual way before
adding Verify. Subtext provides the pieces Verify builds on:

- the **`subtext` MCP server** (HTTP) that also serves Verify's `live-*` /
  `comment-*` / `doc-*` / `artifact-*` / `clip-*` tools (gated server-side by org
  flag),
- the general MCP conventions and security rules in Subtext's `subtext-shared`,
- the canonical `sightmap-authoring` / `sightmap-browser` skills for the
  `.sightmap/` schema.

Verify assumes all of that is present; it only adds what's needed to *use the
verify tools*.

## What this plugin adds

- **Skills** — `subtext-proof` (before/after visual evidence, the inner loop),
  `subtext-onboard` and `subtext-first-session` (first-run flows), the tool
  catalogs `subtext-live` / `subtext-comments` / `subtext-docs` /
  `subtext-tunnel`, the `subtext-recipe-sightmap-setup` recipe, and
  `subtext-verify-shared` (verify-specific conventions). All skill names are
  `subtext-verify-*`-namespaced or unique, so nothing collides with Subtext.
- **MCP server** — only the **`subtext-tunnel`** stdio client (reverse tunnel for
  localhost dev servers), run via
  `npx -y @subtextdev/subtext-cli@latest tunnel mcp`. The shared `subtext` HTTP
  server is already configured by base Subtext, so Verify does not redeclare it.
- **Hook** — a `SessionStart` hook injects the `subtext-verify-using` router so the
  agent reaches for `subtext-proof` before editing UI code.

## Install

**Claude Code**
```
/plugin marketplace add fullstorydev/subtext-verify
/plugin install subtext-verify@subtext-marketplace
```

**Cursor** — install from the Marketplace panel (or a Team Marketplace that
imports this repo).

**Codex** — open `/plugins`, install **subtext-verify** from the repo marketplace.

**Gemini CLI**
```
gemini extensions install https://github.com/fullstorydev/subtext-verify
```

**Sightmap upload dependency (optional).** The sightmap collector script needs
**Python 3.9+ and PyYAML** (`pip install pyyaml`). Only needed if the project has
a `.sightmap/` corpus you want reflected in snapshots and network traces.

## Getting started

A good first run is one real change proven end-to-end — ask your agent something
like *"prove a small UI change on my local dev server."* Under the hood it will:

1. **Connect to your running app.** For a local server, follow the tunnel-first
   flow (`live-tunnel` → `tunnel-connect` → `live-view-new`) and share the
   `trace_url` so you can watch live; the same link replays afterward. See
   `subtext-tunnel`.
2. **Capture before/after proof.** Snapshot the affected area, make the change,
   re-capture, and package screenshots + trace into a proof document. See
   `subtext-proof`.
3. **Bootstrap a sightmap (optional).** Record what was learned about the touched
   components into `.sightmap/` so the next agent starts ahead. See
   `subtext-recipe-sightmap-setup`.

`subtext-onboard` walks a new user through exactly this on their first run.

## Notes

- The `subtext` tool surface is gated server-side by org/feature flag; Verify
  assumes the verification flag is enabled for your org.
- Read-only review of completed sessions and privacy-rule management live in
  [Subtext](https://github.com/fullstorydev/subtext) itself, not here.

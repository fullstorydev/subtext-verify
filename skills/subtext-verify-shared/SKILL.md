---
name: subtext-verify-shared
description: Foundation for the Subtext Verify skills — verify-specific MCP tool prefixes, sightmap upload, and the tunnel server. Read this when a Verify skill lists it in PREREQUISITE.
---

# Verify — Shared

Foundation for the Subtext **Verify** skills (`subtext-live`, `subtext-tunnel`,
`subtext-comments`, `subtext-docs`, `subtext-proof`, `subtext-first-session`, and
the `subtext-recipe-sightmap-setup` recipe). Read this when one lists it in
PREREQUISITE.

> **Layered on Subtext.** The Verify skills assume you already have Subtext
> installed the normal way. General MCP conventions, discovery, and security
> rules live in Subtext's `subtext-shared` — this skill only adds the
> verify-specific pieces.

## MCP servers

Verify's tools are served by the same **subtext** MCP server (HTTP) that the base
Subtext install already configures — the `live-*` / `comment-*` / `doc-*` /
`artifact-*` / `clip-*` tools are gated server-side by org flag, so no extra
wiring is needed for them. A **subtext-eu1** variant exists for EU1 sessions
(`app.eu1.fullstory.com`).

The one server Subtext does **not** configure is the reverse-tunnel client,
**subtext-tunnel** (stdio, its own `tunnel-connect` / `tunnel-disconnect` /
`tunnel-status` namespace), used only for localhost dev servers. It runs via
`npx -y @subtextdev/subtext-cli@latest tunnel mcp`; see the repo's `mcp.json` for
the exact snippet to add to your MCP config if you need localhost tunneling.

## Sightmap upload

Two live tools return a single-use sightmap upload URL:

| Tool | Field | Format |
|------|-------|--------|
| `live-connect` | `sightmap_upload_url` | text line in response |
| `live-tunnel` | `sightmapUploadUrl` | JSON field in response |

If the project has `.sightmap/` definitions, upload them with the bundled
collector after getting the URL and **before** `live-view-new` (tunnel-first
flow) or before interacting with the page:

```bash
# run from the project root (where .sightmap/ lives):
python3 <this skill's directory>/collect_and_upload_sightmap.py --url <sightmap_upload_url>
```

`collect_and_upload_sightmap.py` sits **beside this SKILL.md** — reference it at
that path. It walks `.sightmap/**/*.yaml`, flattens hierarchical components into
compound selectors, collects top-level `memory`, and POSTs using the single-use
token in the URL (no extra auth). Requires **Python 3.9+ and PyYAML**
(`pip install pyyaml`). Do NOT also pass a `sightmap` parameter to `live-connect`.

For the `.sightmap/` schema and how to author a corpus, see the `sightmap-authoring`
skill from Subtext.

## Tool name prefixes

| Prefix | Tools |
|--------|-------|
| `live-` | Browser automation: `live-connect`, `live-disconnect`, `live-view-*`, `live-act-*`, `live-log-*`, `live-net-*`, `live-tunnel`, `live-emulate`, `live-eval-script` |
| `comment-` | Comments: `comment-add`, `comment-list`, `comment-reply`, `comment-resolve` |
| `doc-` | Proof documents: `doc-create`, `doc-update`, `doc-attach`, `doc-close`, `doc-read`, `doc-diff`, `doc-list` |
| `artifact-` | Artifacts: `artifact-upload`, `artifact-url` |
| `clip-` | Replay clips: `clip-create` |

## Discovering tool parameters

Each MCP tool is self-describing — the schema is available at call time. Consult
the atomic skill (`subtext-live`, `subtext-comments`, `subtext-docs`) for which
tools exist, then let the schema guide parameters.

## Security

General rules are in Subtext's `subtext-shared`. Verify adds one: **confirm with
the user before any write operation that modifies production data.**

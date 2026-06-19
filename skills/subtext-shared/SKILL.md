---
name: subtext-shared
description: Foundation skill for the Subtext Verify plugin. MCP tool conventions, sightmap upload, and security rules.
---

# Shared

Foundation for all Subtext Verify skills. Read this when any workflow or recipe lists it in PREREQUISITE.

## MCP Servers

Tools are served from the **subtext** MCP server (HTTP). A **subtext-eu1** variant exists for EU1 data center sessions (`app.eu1.fullstory.com`). The reverse-tunnel client is a separate **subtext-tunnel** stdio MCP server with its own tool namespace (`tunnel-connect`, `tunnel-disconnect`, `tunnel-status`). The agent framework resolves tool prefixes automatically based on the configured MCP servers — you do not need to hardcode prefixes.

## Sightmap Upload

Two tools return a sightmap upload URL:

| Tool | Field | Format |
|------|-------|--------|
| `live-connect` | `sightmap_upload_url:` | text line in response |
| `live-tunnel` | `sightmapUploadUrl` | JSON field in response |

If the project has `.sightmap/` definitions, upload them via the side-band script after getting the URL and **before** `live-view-new` (tunnel-first flow) or before interacting with the page:

```bash
python3 ${CLAUDE_PLUGIN_ROOT}/skills/subtext-shared/collect_and_upload_sightmap.py --url <sightmap_upload_url>
```

The upload uses a single-use token embedded in the URL — no additional auth is needed. Do NOT pass the `sightmap` parameter directly to `live-connect`.

## Tool Name Prefixes

Tools within the subtext server are grouped by prefix:

| Prefix | Tools |
|--------|-------|
| `live-` | Browser automation: `live-connect`, `live-disconnect`, `live-view-*`, `live-act-*`, `live-log-*`, `live-net-*`, `live-tunnel`, `live-emulate`, `live-eval-script` |
| `comment-` | Comments: `comment-add`, `comment-list`, `comment-reply`, `comment-resolve` |
| `doc-` | Proof documents: `doc-create`, `doc-update`, `doc-attach`, `doc-close`, `doc-read`, `doc-diff`, `doc-list` |
| `artifact-` | Artifacts: `artifact-upload`, `artifact-url` |
| `clip-` | Replay clips: `clip-create` |

## Discovering MCP Tool Parameters

Each MCP tool is self-describing. If you're unsure about parameters, the tool's schema is available at call time. Don't memorize parameter lists — consult the atomic skill (`subtext-live`, `subtext-comments`, or `subtext-docs`) for which tools exist, then let the schema guide parameter usage.

## Security Rules

- Never expose API tokens, session tokens, or credentials in output.
- Confirm with the user before any write operation that modifies production data.
- Session URLs may contain sensitive user data — don't log or repeat them unnecessarily.

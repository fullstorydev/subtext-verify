# subtext-verify CLI

A minimal Go CLI whose only job is the **reverse tunnel** that lets the Subtext
hosted browser reach a local dev server.

## Commands

```
subtext tunnel mcp          # MCP stdio server (tunnel-connect / -disconnect / -status)
subtext tunnel connect      # connect a local dev server to the relay (foreground or --detach)
subtext tunnel disconnect   # stop a detached tunnel (--tunnel-id or --all)
subtext tunnel status       # list running tunnel daemons
subtext version
```

`subtext tunnel mcp` is the stdio server referenced by the `subtext-tunnel`
entry in `.mcp.json`:

```json
"subtext-tunnel": { "command": "npx", "args": ["-y", "@subtextdev/subtext-cli@latest", "tunnel", "mcp"] }
```

No API key is required for `tunnel mcp` — the relay URL supplied to
`tunnel-connect` already carries an embedded auth token (from `live-tunnel`).

## Build

```
go build ./cmd/subtext
```

## Provenance

This is the tunnel-only subset of the Subtext Go CLI (the `subtext tunnel mcp`
server from fullstorydev/subtext PR #98), vendored here so the Verify plugin
owns its own tunnel client. The MCP-passthrough namespaces (live, comment, doc,
review, privacy, …) from the original CLI are intentionally omitted — those
tools are reached over HTTP via the `subtext` MCP server, not the CLI.

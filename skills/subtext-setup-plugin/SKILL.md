---
name: subtext-setup-plugin
description: Install the Subtext Verify plugin and configure MCP servers. Authenticates via OAuth or API Key.
---

# Setup Plugin

Install and verify the Subtext Verify plugin/extension. Works for Claude Code, Cursor, Codex, and Gemini CLI.

## Pre-check

Verify the plugin is working by testing actual connectivity — do NOT read config files or plugin cache directories.

**Step 1: Test MCP connectivity**

Try calling a lightweight MCP tool to verify the server is reachable — for example, list the available tools on the `subtext` MCP server. If the call succeeds, the plugin is installed and connected.

Check these servers:
- `subtext` — required (live, comments, docs, artifacts, clips)
- `subtext-tunnel` — optional (local tunnel client, for localhost dev servers)

If MCP tools are available, the plugin is working. Report which servers connected and move on.

**Step 2: Verify local dependencies**

The sightmap upload script needs Python 3 + PyYAML:

```bash
python3 --version 2>&1 && python3 -c "import yaml; print('PyYAML OK')" 2>&1
```

If missing:
- No Python 3: suggest `brew install python@3.12` (macOS)
- No PyYAML: suggest `python3 -m pip install pyyaml`

**If all pass:** report "Plugin is set up — MCP servers connected, dependencies OK." and exit.

## Install

If MCP tools are not available, install the plugin. The command depends on the platform.

**Claude Code:**

```
/plugin marketplace add fullstorydev/subtext-review
/plugin install subtext-verify@subtext-marketplace
```

**Gemini CLI:**

```
gemini extensions install https://github.com/fullstorydev/subtext-verify
```

Note: Slash commands can't be executed by the agent — the user must run them directly.

## MCP connectivity failed

If the MCP connectivity test fails:

1. The subtext MCP server did not respond.
2. Double-check authentication settings for the MCP servers in the tool configuration.
3. Authenticate via the OAuth flow provided by your tool, or configure an API key header for the MCP server.

Re-run the connectivity check after authenticating.

## Explain

After setup, explain what was installed:

- **Skills** — `subtext-proof` (before/after evidence for UI changes), `subtext-onboard` (first-run walkthrough), `subtext-first-session` (agent-driven exploration), plus the underlying tool catalogs (`subtext-live`, `subtext-comments`, `subtext-docs`, `subtext-tunnel`, `subtext-sightmap`).
- **MCP servers** — `subtext` (live, comments, docs, artifacts, clips) and `subtext-tunnel` (local reverse-tunnel client for localhost dev servers).
- **Sightmap** — semantic component mapping, uploaded automatically after each connection. See `subtext-sightmap`.

> Read-only review of completed sessions and privacy-rule management live in the separate **Subtext Review** plugin.

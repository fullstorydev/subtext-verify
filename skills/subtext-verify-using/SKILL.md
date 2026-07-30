---
name: subtext-verify-using
description: Use when starting any conversation that may involve rendered UI, driving a running app, or producing reviewer-facing evidence (screenshots, viewer links, code diffs). Establishes how the Subtext Verify skills compose and when to invoke them before any response or action.
---

<EXTREMELY-IMPORTANT>
If the task touches rendered UI, a running app, or producing
proof-of-work evidence, you MUST invoke the relevant Subtext Verify
skill before responding.
</EXTREMELY-IMPORTANT>

Subtext Verify is a **companion** to Subtext — it assumes Subtext is already
installed (its `subtext` MCP server, shared conventions, and `sightmap-authoring`
skill). This router covers the Verify surface: driving a live browser and
capturing proof.

## Where this skill applies

Subtext Verify runs *where the work happens*. Unlike many process skills,
this includes subagent contexts.

- **Subagent doing UI/UX work or producing reviewer-facing evidence:**
  MUST invoke. Your orchestrator depends on you to surface evidence —
  screenshots, viewer URLs, comments — back up the chain.
- **Subagent doing purely backend / non-visual work:** trigger surface
  doesn't apply, skip.
- **Orchestrator running directly:** same rule, you invoke the relevant
  skill yourself.

## How to Access Skills

- **Claude Code & Cursor:** use the `Skill` tool.
- **Codex:** Skills load natively from `~/.agents/skills/`. Read the relevant SKILL.md directly when its description matches your task.
- **Gemini CLI:** Skills activate via the `activate_skill` tool.

## When to Reach for Subtext Verify

| Signal | Reach for |
|--------|-----------|
| Making UI/visual changes | `subtext-proof` |
| Need to drive a hosted browser | `subtext-live` |
| Connecting to a localhost dev server | `subtext-tunnel` |
| First run / walking a new user through Verify | `subtext-onboard` |
| Naming components / runtime model | `sightmap-authoring` (from Subtext) |

> Reviewing a *completed* session from a URL (read-only summary, privacy rules)
> is **Subtext**'s job, not Verify's.

## The Rule

Invoke the relevant Subtext Verify skill BEFORE any response or action that
touches the trigger surface. Even a 1% chance counts.

## Red Flags

These thoughts mean STOP — you're rationalizing:

| Thought | Reality |
|---------|---------|
| "I'll just check the diff" | Visual changes need visual proof. |
| "Tests passed, that's enough" | Tests verify code, not UX. |
| "I don't need a session for this small change" | Small UI changes regress silently. |
| "I'll describe what changed" | Screenshots > prose. |
| "Let me explore the app first" | `subtext-proof` tells you HOW to explore. |
| "I remember how proof works" | Skills evolve. Read current version. |
| "I got `Control transferred to human viewer`, let me retry" | Operator state is enforced server-side. Don't retry — poll `live-signal` until `operator=agent`. |

## Composition

- **Atomics** (`subtext-verify-shared`, `subtext-live`, `subtext-tunnel`, `subtext-comments`, `subtext-docs`) — tool catalogs.
- **Workflows** (`subtext-proof`) — orchestration. `subtext-proof` is the inner loop, captured as a recorded session.
- **Recipes** (`subtext-recipe-sightmap-setup`) — short step lists.
- **Onboarding** (`subtext-onboard`, `subtext-first-session`) — first-run flows.

```
proof ──▶ session recorded ──▶ (optional) review in Subtext
```

## Skill Types

- **Rigid** (`subtext-proof`): follow exactly.
- **Flexible** (atomics): adapt to context.

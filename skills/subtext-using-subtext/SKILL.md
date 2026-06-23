---
name: subtext-using-subtext
description: Use when starting any conversation that may involve rendered UI, driving a running app, or producing reviewer-facing evidence (screenshots, viewer links, code diffs). Establishes how the Subtext Verify skills compose and when to invoke them before any response or action.
---

<EXTREMELY-IMPORTANT>
If the task touches rendered UI, a running app, or producing
proof-of-work evidence, you MUST invoke the relevant Subtext skill
before responding.
</EXTREMELY-IMPORTANT>

## Where this skill applies

Subtext runs *where the work happens*. Unlike many process skills,
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

## When to Reach for Subtext

| Signal | Reach for |
|--------|-----------|
| Making UI/visual changes | `subtext-proof` |
| Need to drive a hosted browser | `subtext-live` |
| Connecting to a localhost dev server | `subtext-tunnel` |
| Setting up a new project | `subtext-onboard` |
| Naming components / runtime model | `subtext-sightmap` |

> Reviewing a *completed* session from a URL (read-only summary, privacy rules) is the separate **Subtext Review** plugin's job, not this one.

## The Rule

Invoke the relevant Subtext skill BEFORE any response or action that
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

- **Atomics** (`subtext-shared`, `subtext-live`, `subtext-sightmap`, `subtext-tunnel`, `subtext-comments`, `subtext-docs`) — tool catalogs.
- **Workflows** (`subtext-proof`) — orchestration. `subtext-proof` is the inner loop, captured as a recorded session.
- **Recipes** (`subtext-recipe-sightmap-setup`) — short step lists.
- **Onboarding** (`subtext-onboard`, `subtext-setup-plugin`, `subtext-first-session`) — first-time user setup.

```
proof ──▶ session recorded ──▶ (optional) review in the Subtext Review plugin
```

## Skill Types

- **Rigid** (`subtext-proof`): follow exactly.
- **Flexible** (atomics): adapt to context.

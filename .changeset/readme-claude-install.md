---
"subtext-verify": patch
---

Fix the Claude Code install instructions. Verify is listed in the Subtext
marketplace (`fullstorydev/subtext`), not a marketplace in this repo, so
`/plugin marketplace add fullstorydev/subtext-verify` failed with "Marketplace
file not found". The README now adds `fullstorydev/subtext`, installs
`subtext-verify@subtext-marketplace`, and enables it (the marketplace lists
Verify as disabled by default).

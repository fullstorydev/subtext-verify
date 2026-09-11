---
"subtext-verify": minor
---

subtext-verify-shared: upload the `.sightmap/` corpus with the `sightmap` CLI
(`sightmap export --url <url>`) instead of the bundled Python collector, and remove
`collect_and_upload_sightmap.py`.

Matches the Subtext change: `sightmap export` routes the upload through the Go
loader and POSTs the whole canonical wire (components, views/routes, requests,
messages, memory, tags), dropping the Python 3 / PyYAML dependency. Applies to the
`live-connect` / `live-tunnel` upload URLs; the `subtext-tunnel` setup note is
updated to require the `sightmap` binary rather than PyYAML.

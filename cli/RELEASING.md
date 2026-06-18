# Releasing the Subtext Verify CLI

The CLI ships two ways from one build: a GitHub Release with prebuilt binaries
(under a `cli/v*` tag) and the `@subtextdev/subtext-cli` npm package, whose
`postinstall` downloads the matching binary from that release.

## Prerequisites

- Write access to `fullstorydev/subtext-verify`.
- (First release only) `@subtextdev` npm scope access and a Trusted Publisher
  configured on npmjs.com for **this** repo (see first-release checklist).

## Release process

1. **Bump the npm wrapper version** in `cli/npm/package.json` to match the tag
   you're about to push (without the `cli/v` prefix). The release workflow fails
   if the tag and `package.json` disagree.

   ```bash
   git add cli/npm/package.json
   git commit -m "cli: bump npm version to 1.2.0"
   git push
   ```

2. **Tag and push** — triggers `.github/workflows/release-cli.yml`:

   ```bash
   git tag cli/v1.2.0
   git push origin cli/v1.2.0
   ```

3. The workflow runs `go test ./...`, builds darwin/linux/windows × amd64/arm64
   via GoReleaser, creates the GitHub Release `CLI v1.2.0` under `cli/v1.2.0`,
   then publishes `@subtextdev/subtext-cli` via OIDC (Trusted Publisher — no
   `NPM_TOKEN` secret).

4. **Smoke test:**

   ```bash
   npx @subtextdev/subtext-cli@1.2.0 tunnel mcp   # starts the stdio MCP server
   go install github.com/fullstorydev/subtext-verify/cli/cmd/subtext@v1.2.0
   ```

## GoReleaser and the cli/v* tag prefix

GoReleaser OSS doesn't support the `prefix/vX.Y.Z` monorepo tag format (a Pro
feature), and `go install` needs the `cli/v` prefix to resolve the module. The
workflow strips the prefix into `GORELEASER_CURRENT_TAG`, passes
`--skip=publish,validate`, and creates the Release manually so the tag stays
`cli/vX.Y.Z`.

## Snapshot builds (local)

```bash
cd cli
goreleaser release --snapshot --clean --skip=publish
./dist/subtext_darwin_arm64_v8.0/subtext tunnel --help
```

## Package name / version continuity

This repo publishes `@subtextdev/subtext-cli`, the same package previously
released from `fullstorydev/subtext` (last at `1.1.0`, per PR #98). Versions must
strictly increase across the move — start at `1.2.0` or higher. Confirm the
npmjs Trusted Publisher points at `fullstorydev/subtext-verify` /
`release-cli.yml` before the first publish from here.

## First release checklist

- [ ] `@subtextdev` npm scope exists.
- [ ] Trusted Publisher on npmjs.com: GitHub Actions, org `fullstorydev`, repo
      `subtext-verify`, workflow `release-cli.yml`, allow `npm publish`.
- [ ] First published version > last version published from the old repo.

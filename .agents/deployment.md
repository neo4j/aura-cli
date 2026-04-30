# Deployment

## Strategy: GitHub Releases via GoReleaser

Releases are triggered automatically when a push to `main` modifies `CHANGELOG.md`.

## Release Flow

1. `changie` CI workflow runs on push to `main`, batches new changelog entries, merges them, and opens a PR titled `Release <version>`
2. When that release PR is merged (updating `CHANGELOG.md`), the `release` workflow fires
3. GoReleaser builds multi-platform binaries and publishes a GitHub Release with release notes from `.changes/<version>.md`

## Platforms

| OS | Arch | Format |
|----|------|--------|
| Linux | amd64, arm64 | tar.gz |
| macOS | amd64, arm64 | tar.gz (signed + notarized) |
| Windows | amd64 | zip |

## macOS Code Signing

macOS binaries are signed with a `.p12` certificate and notarized via Apple's App Store Connect. Credentials are stored as GitHub secrets:
- `MACOS_SIGN_P12`, `MACOS_SIGN_PASSWORD`
- `MACOS_NOTARY_ISSUER_ID`, `MACOS_NOTARY_KEY_ID`, `MACOS_NOTARY_KEY`

## Version Injection

The version string is injected at link time: `-X "main.Version={{.Env.GORELEASER_CURRENT_TAG}}"`. For local builds, pass any value via `GORELEASER_CURRENT_TAG`.

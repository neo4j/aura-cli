# Deployment

## Strategy: GitHub Releases via GoReleaser

Releases are triggered automatically when a push to `main` modifies `CHANGELOG.md`.

## Dual-Binary Releases

Each release produces two separate binaries:

| Binary | Description |
|--------|-------------|
| `aura-cli` | Standalone Aura CLI (entrypoint: `./neo4j-cli/aura/cmd`) |
| `neo4j-cli` | Super-CLI that wraps `aura-cli` as `neo4j-cli aura <subcommand>` (entrypoint: `./neo4j-cli`) |

Each binary gets its own archive per platform/arch combination.

## Cascade Versioning Policy

`aura-cli` and `neo4j-cli` have **separate version numbers and separate `CHANGELOG.md` files**. The rule:

> Every `aura-cli` release that is bundled into a `neo4j-cli` release must also have a corresponding `neo4j-cli` changelog entry.

Run `changie new` for **both** binaries when a change affects the `aura` subcommand tree that ships in `neo4j-cli`. This keeps both changelogs accurate and ensures users of either binary see meaningful release notes.

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

macOS binaries (both `aura-cli` and `neo4j-cli`) are signed with a `.p12` certificate and notarized via Apple's App Store Connect. Both binary IDs are listed in `notarize.macos[].ids`. Credentials are stored as GitHub secrets:
- `MACOS_SIGN_P12`, `MACOS_SIGN_PASSWORD`
- `MACOS_NOTARY_ISSUER_ID`, `MACOS_NOTARY_KEY_ID`, `MACOS_NOTARY_KEY`

## Version Injection

The version string is injected at link time: `-X "main.Version={{.Env.GORELEASER_CURRENT_TAG}}"`. For local snapshot builds:

```bash
GORELEASER_CURRENT_TAG=dev goreleaser release --snapshot --clean
```

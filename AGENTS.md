<!-- BEGIN GENERATED: AGENTS-MD -->

# AGENTS.md

Learnings and patterns for future agents working on this project.

## Feedback Instructions

TEST COMMANDS: [`go test ./...`]
BUILD COMMANDS: [`go build -o bin/ ./...`, `go run neo4j-cli/main.go aura-cli`]
LINT COMMANDS: [`go vet ./...`, `staticcheck ./...`]
FORMAT COMMANDS: [`go fmt ./...`]
LICENSE CHECK: [`addlicense -f ./addlicense -check $(find . -name "*.go" -type f -print0 | xargs -0)`]

## Project Overview

PRIMARY LANGUAGES: [Go]

Neo4j Aura CLI (`aura-cli`) is a command-line tool for interacting with the [Neo4j Aura API](https://neo4j.com/docs/aura/platform/api/specification/). It allows users to manage Aura instances, credentials, tenants, deployments, data APIs, and more from the terminal.

## Build System

BUILD SYSTEMS: [Go toolchain, GoReleaser, changie]

See [`.agents/build.md`](.agents/build.md) for full details.

- Local build: `go build -o bin/ ./...`
- Local run (no build): `go run neo4j-cli/main.go aura-cli`
- Multi-platform snapshot: `GORELEASER_CURRENT_TAG=dev goreleaser release --snapshot --clean`
- All `.go` files must start with the Neo4j copyright header (enforced in CI via `addlicense`)
- PRs require a changelog entry added via `changie new`

## Testing Framework

TESTING FRAMEWORKS: [Go testing, testify, afero (in-memory FS)]

See [`.agents/testing.md`](.agents/testing.md) for full details.

- Tests are colocated with source as `*_test.go` files
- Run with `go test ./...`; CI runs on ubuntu, windows, and macos
- Mock HTTP server and filesystem helpers live in `neo4j-cli/aura/internal/test/testutils/`
- `neo4j-cli/` (the super-CLI package) has no test files; this is a pre-existing gap

## Architecture

ARCHITECTURE PATTERN: Cobra command tree — one file per leaf command, directory structure mirrors command hierarchy

See [`.agents/architecture.md`](.agents/architecture.md) for full details.

```
neo4j-cli/aura/
  cmd/main.go              # Binary entrypoint
  aura.go                  # Root command, registers subcommands
  internal/
    api/                   # HTTP client for Neo4j Aura REST API
    flags/                 # Custom reusable flag types
    output/                # JSON + table rendering
    subcommands/           # One directory per resource, one file per action
      instance/, tenant/, credential/, config/,
      deployment/, dataapi/graphql/, graphanalytics/,
      import/, customermanagedkey/
common/
  clicfg/                  # Config, credentials, project state (OS-specific paths)
  clierr/                  # Shared error types
```

Key CLI conventions (see `CONTRIBUTING.md`):
- Singular nouns for commands (`instance`, not `instances`)
- `<resource> <action>` form (`instance list`, not `list-instance`)
- One positional argument max; extras become flags
- `--output json|table` for all read commands
- `--await` flag for async operations

## Deployment

DEPLOYMENT STRATEGY: GitHub Releases via GoReleaser, triggered by CHANGELOG.md updates on `main`

See [`.agents/deployment.md`](.agents/deployment.md) for full details.

- `changie` batches changelog entries and opens release PRs automatically
- Merging a release PR triggers GoReleaser to publish binaries for linux/windows/darwin (amd64 + arm64)
- macOS binaries are code-signed and notarized

## Makefile Notes

- `license-check` target uses `$(GOPATH)/bin/addlicense` (not bare `addlicense`) — GOPATH/bin may not be on PATH
- `license-check` requires a Unix shell (`find` + `xargs`); won't work natively on Windows

## golangci-lint Notes

- Version installed: v2.11.4 (via Homebrew)
- golangci-lint v2 requires `version: "2"` at the top of `.golangci.yml`
- In v2, `gofmt` is a **formatter** (not a linter); put it under `formatters.enable`, not `linters.enable`
- Use `linters.default: none` to disable auto-enabled defaults (e.g. `ineffassign`) and run only explicitly listed linters
- Config lives at `.golangci.yml` in repo root

---

_This AGENTS.md was generated using agent-based project discovery._

<!-- END GENERATED: AGENTS-MD -->

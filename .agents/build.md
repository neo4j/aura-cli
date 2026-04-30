# Build System

## Tools

- **Go toolchain** — standard `go build`, `go run`, `go vet`, `go fmt`
- **GoReleaser** — multi-platform release builds
- **changie** — changelog management

## Commands

```bash
# Local development build
go build -o bin/ ./...

# Run without building (local dev)
go run neo4j-cli/main.go aura-cli

# Multi-platform release snapshot (local)
GORELEASER_CURRENT_TAG=dev goreleaser release --snapshot --clean

# Static analysis
go vet ./...
go fmt ./...
staticcheck ./...
```

## Release Process

Releases are triggered automatically in CI when `CHANGELOG.md` is updated on `main`. GoReleaser produces binaries for:
- `linux/amd64`, `linux/arm64`
- `windows/amd64`
- `darwin/amd64`, `darwin/arm64`

macOS binaries are code-signed and notarized using App Store Connect keys stored in GitHub secrets.

## Changelog

Uses `changie` to manage changelog entries. Before merging a PR, run:

```bash
changie new
```

Commit the generated file from `.changes/`. The `changie` CI workflow auto-batches entries and opens release PRs.

## Build Output

The project name is `aura-cli`. The binary entrypoint is at `neo4j-cli/aura/cmd/main.go`. The `neo4j-cli` directory also contains a top-level `main.go` that acts as a dispatcher.

# Neo4j CLI

## Build

```bash
go build -o bin/ ./...
```

## Run

```bash
go run cmd/aura/main.go
```

## Test

```bash
go test -v ./...
```

## Future notes

Build with something like:

```bash
go build -ldflags "-X main.Version `git describe --tags --abbrev=0`" aura
```

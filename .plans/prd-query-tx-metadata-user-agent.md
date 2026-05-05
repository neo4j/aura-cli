# PRD: User-Agent header + tx metadata for `query` and `:schema`

## Overview

`neo4j-cli query` (and `query :schema`) currently POST to `/db/<db>/query/v2` with no client identity: no `User-Agent` header, no `txMetadata` in the request body. Server-side ops cannot tell from `query.log` or `SHOW TRANSACTIONS` which client issued a query.

This feature adds both: a `User-Agent: neo4j-cli/v<version>` header on every HTTP request, and a `txMetadata` map in the request body tagging each transaction with `{app: "neo4j-cli", type: <type>}`. The `type` distinguishes user-issued queries (`"user-direct"`, mirroring upx) from `:schema` probe queries (`"schema"`). Tracks Linear [CLI-28](https://linear.app/neo4j/issue/CLI-28).

Scope is the `query` package only. Aura subcommand HTTP code paths are untouched.

## Goals

- Identify `neo4j-cli` to the database in `query.log` and `SHOW TRANSACTIONS`.
- Distinguish user queries from `:schema` probe queries in server-side logs.
- Match naming conventions of existing Neo4j clients (browser, upx) referenced in CLI-28.
- Single HTTP call site change (`runStatement`) covers both `query` and `:schema`.
- Zero new dependencies; reuse existing `cfg.Version` plumbing.

## Non-Goals

- OS / arch suffix on the User-Agent (`neo4j-cli/v1.2.3 (darwin; arm64)`) — kept bare.
- Aura subcommand HTTP requests (`neo4j-cli/aura/internal/api/`) — already send `Neo4jCLI/%s`; not reunified here.
- Configurability of the UA / tx metadata via flags or env — values are constants.
- Bolt user-agent (the v2 HTTP API is the only transport `query` uses).
- Additional tx metadata fields beyond `app` and `type`.
- Changes to the `:schema` output / probe query set.

## Requirements

### Functional Requirements

- **REQ-F-001**: Add a `userAgent string` field to the `conn` struct in `neo4j-cli/query/connect.go:49`.
- **REQ-F-002**: In `resolveConn` (`neo4j-cli/query/connect.go:84`), populate `conn.userAgent` from `cfg.Version` as `"neo4j-cli/v" + cfg.Version`. When `cfg.Version` is empty, fall back to `"neo4j-cli/vdev"` (matches the `Version = "dev"` default at `neo4j-cli/app/app.go:22`).
- **REQ-F-003**: Change `runStatement` signature (`neo4j-cli/query/connect.go:253`) to add a final `txType string` parameter. Callers always pass a non-empty type; do not provide a default.
- **REQ-F-004**: In `runStatement`, set `body["txMetadata"] = map[string]any{"app": "neo4j-cli", "type": txType}` before JSON-marshaling the request body.
- **REQ-F-005**: In `runStatement`, after `req.SetBasicAuth(...)`, set `req.Header.Set("User-Agent", c.userAgent)`. Skip the `Set` when `c.userAgent` is empty so Go's default UA is not clobbered with an empty string.
- **REQ-F-006**: `neo4j-cli/query/run.go:75` passes `"user-direct"` as the `txType`.
- **REQ-F-007**: All six `runStatement` call sites in `neo4j-cli/query/schema.go` (lines 158, 179, 210, 230, 244, 252) pass `"schema"` as the `txType`.
- **REQ-F-008**: Add a changelog entry: `changie new --projects neo4j-cli --kind Minor --body "query/:schema: send User-Agent header and txMetadata so server logs identify the CLI"`.

### Non-Functional Requirements

- **REQ-NF-001**: No new module dependencies. The implementation uses only the existing `net/http` + `encoding/json` imports already in `connect.go`.
- **REQ-NF-002**: No import of `neo4j-cli/app` from the `query` package — would create a cycle. Version flows via `cfg.Version`, which `query.NewCmd(cfg)` already receives.
- **REQ-NF-003**: All existing query-package tests stay green after touching `&conn{...}` literals and `runStatement` call sites; no behavioural regressions.
- **REQ-NF-004**: New unit tests cover (a) the resolved User-Agent string for a populated `cfg.Version`, (b) the fallback for empty `cfg.Version`, (c) the body shape for `runStatement(..., "user-direct")`, and (d) the body shape for `runStatement(..., "schema")`.
- **REQ-NF-005**: `make test`, `make fmt-check`, `make lint`, and `make license-check` all pass before merge (per AGENTS.md gate).
- **REQ-NF-006**: No file under `neo4j-cli/aura/` is modified by this change.

## Technical Considerations

**Single call site.** `runStatement` (`neo4j-cli/query/connect.go:253`) is the only place the `query` package issues HTTP requests. Patching it once covers both leaf commands. Verified via `grep` — the seven call sites are all in `run.go` and `schema.go`.

**Why a string parameter, not a struct.** A bare `txType string` keeps churn minimal across the seven call sites. A future need for per-call metadata richness can introduce a small struct then; YAGNI now.

**Why not a method on `conn`.** `runStatement` is already package-private and takes `c *conn`; adding a method would not improve clarity and would force splitting transport from request shaping. Keeping it as a function preserves the file's current style.

**Tx metadata naming.** Matches upx's `{ app: 'neo4j-import', type: 'user-direct' }`. `app: "neo4j-cli"` for both user and schema; `type: "user-direct" | "schema"` differentiates them. Server-side filtering via `SHOW TRANSACTIONS YIELD metaData` then becomes trivial (`WHERE metaData['app'] = 'neo4j-cli'`, optionally further by `type`).

**User-Agent format.** `neo4j-cli/v<version>` mirrors neo4j-browser's `neo4j-browser/v${version}` shape. Bare — no platform suffix. Distinct from Aura's `Neo4jCLI/%s` (used in `neo4j-cli/aura/internal/api/api.go:17`) which targets a different API; we deliberately do not reuse it so log filtering can distinguish DB-level from Aura-control-plane traffic.

**API support.** Confirmed via the Neo4j Query API docs: `txMetadata` is a documented body field on `POST /db/<name>/query/v2`. Logged in server `query.log` and exposed via `SHOW TRANSACTIONS YIELD metaData`. Supported on Neo4j 5.x.

**Version plumbing.** `cfg.Version` is already populated end-to-end (`common/clicfg/clicfg.go:34,72` ← `Version` in `neo4j-cli/main.go` and `neo4j-cli/aura/cmd/main.go`). `query.NewCmd(cfg)` already receives `cfg`. No new wiring required.

**Files touched.**
- `neo4j-cli/query/connect.go` — add `userAgent` field; populate in `resolveConn`; extend `runStatement` signature; set header + tx metadata.
- `neo4j-cli/query/run.go` — pass `"user-direct"` to `runStatement`.
- `neo4j-cli/query/schema.go` — pass `"schema"` to all six `runStatement` calls.
- `neo4j-cli/query/connect_test.go` — assert UA header + body shape; populate `userAgent` on every `&conn{...}` literal; new `TestResolveConn_UserAgent` (table-driven); new `TestRunStatement_SchemaTxType`.
- `neo4j-cli/query/schema_test.go` — populate `userAgent` on the single `&conn{...}` literal at line 371; pass `"schema"` to any direct `runStatement` calls.
- `.changes/unreleased/neo4j-cli-Minor-*.yaml` — changelog entry via `changie new` (or hand-authored per AGENTS.md if changie isn't installed).

**Files explicitly NOT touched.** `neo4j-cli/aura/**`, `common/clicfg/**`, `neo4j-cli/app/app.go`, the `:schema` probe Cypher itself, the `query_https_smoke_test.go` test (will be re-run as a check, not modified).

## Acceptance Criteria

- [ ] `bin/neo4j-cli query --uri ... "RETURN 1 AS n"` against a local Neo4j 5 — `query.log` shows the entry with `metaData={app=neo4j-cli, type=user-direct}` and the request's `User-Agent` is `neo4j-cli/v<version>` (verifiable via `tcpdump` / Neo4j HTTP access log if enabled).
- [ ] `bin/neo4j-cli query :schema --uri ...` — `query.log` entries for the probe queries show `metaData={app=neo4j-cli, type=schema}`.
- [ ] `SHOW TRANSACTIONS YIELD transactionId, metaData` while a long-running CLI statement is in flight — returns the corresponding `app/type` map.
- [ ] `TestRunStatement_HappyPath` asserts the `User-Agent` header equals `"neo4j-cli/vtest"` and the body contains `txMetadata.app == "neo4j-cli"`, `txMetadata.type == "user-direct"`.
- [ ] `TestRunStatement_SchemaTxType` asserts the body's `txMetadata.type == "schema"` when `runStatement` is called with `"schema"`.
- [ ] `TestResolveConn_UserAgent` (table-driven) asserts `c.userAgent == "neo4j-cli/v1.2.3"` for `cfg.Version == "1.2.3"` and `c.userAgent == "neo4j-cli/vdev"` for `cfg.Version == ""`.
- [ ] `make test` green.
- [ ] `make fmt-check` green (no gofmt drift).
- [ ] `make lint` green.
- [ ] `make license-check` green.
- [ ] `NEO4J_HTTPS_TEST=1 go test -run TestHTTPS_Smoke -v ./neo4j-cli/query/...` green — confirms a real Neo4j accepts the new body field.
- [ ] `.changes/unreleased/neo4j-cli-Minor-*.yaml` exists with the agreed body.
- [ ] No file under `neo4j-cli/aura/` is modified.

## Out of Scope

- OS / arch suffix on User-Agent.
- Aura subcommand HTTP code paths.
- Bolt-protocol user-agent (CLI doesn't use Bolt).
- Configurable UA / tx metadata via flag or env.
- Additional `txMetadata` fields (e.g. `user`, `host`, `cwd`, request id).
- Changes to `:schema` probe query set or output formatting.
- Refactor of `runStatement` into a struct-based options API.

## Open Questions

(none — all clarifying questions resolved during planning.)

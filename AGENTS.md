<!-- BEGIN GENERATED: AGENTS-MD -->

# AGENTS.md

Learnings and patterns for future agents working on this project.

## Feedback Instructions

TEST COMMANDS: [`make test`]
BUILD COMMANDS: [`make build`, `make run-aura`, `make run-neo4j`]
LINT COMMANDS: [`make lint`]
FORMAT COMMANDS: [`make fmt-check`] — runs `gofmt -l .` and fails on any output. `make fmt` rewrites silently and is NOT a gate; use `make fmt-check` to verify. CI's golangci-lint v2 includes `gofmt` as a formatter and will fail the build on unformatted code.
LICENSE CHECK: [`make license-check`]

**Always run `make test` AND `make fmt-check` as final gates before marking any task or plan complete.** All tests must pass and no file may need gofmt — a build that compiles but has failing tests or unformatted code is not done. `make fmt-check` is the local equivalent of CI's gofmt linter, so drift fails before the push instead of after.

## Cobra Command Layout

The repo follows a strict one-file-per-leaf cobra layout. Every command tree under `neo4j-cli/aura/internal/subcommands/<resource>/` and `common/skill/` follows it. Mirror it for any new command tree:

- **Parent file `<resource>.go`** — defines `NewCmd(cfg, ...) *cobra.Command`, registers persistent flags, calls `cmd.AddCommand(newXxxCmd(cfg, ...))` for each leaf. Keep it small (≤80 lines).
- **One file per leaf action `<action>.go`** — defines a private constructor like `newInstallCmd(cfg, ...) *cobra.Command` containing the leaf's flags + `RunE`. No leaf bodies inlined into the parent.
- **Colocated tests `<action>_test.go`** — tests for each leaf live next to its source.
- **Shared test helpers** in `<resource>_helpers_test.go` (or similar) when needed.

Examples: `neo4j-cli/aura/internal/subcommands/instance/{instance.go, list.go, list_test.go, get.go, delete.go, ...}` and `common/skill/{skill.go, install.go, remove.go, list.go, check.go, helpers.go}`.

Don't inline multiple leaves in the parent. Don't name the parent `command.go` — name it after the resource so `grep <resource>.go` finds it. Adding a new leaf = new `<action>.go` + `<action>_test.go`, plus one `cmd.AddCommand(...)` line in the parent.

## Project Overview

PRIMARY LANGUAGES: [Go]

Neo4j CLI (`neo4j-cli`) is a command-line tool for interacting with Neo4j.

## Build System

BUILD SYSTEMS: [Go toolchain, Makefile, golangci-lint, GoReleaser, changie]

See [`.agents/build.md`](.agents/build.md) for full details.

- Local build: `make build` (produces `bin/aura-cli` and `bin/neo4j-cli`)
- Local run (no build): `make run-aura` / `make run-neo4j`
- Release build (current platform, ldflags baked in): `make snapshot` (uses goreleaser, outputs to `bin/`)
- npm publish dry-run (template + ordering check): `make npm-publish-dry`. Works against empty `dist/` because `publish.sh --dry-run` stubs missing platform binaries with a 1-byte placeholder; run `make snapshot` first if you want real binaries packed. Real-binary path (CI) still hard-errors on missing binaries — the stub is dry-run-only.
- All `.go` files must start with the Neo4j copyright header (enforced in CI via `addlicense`)
- PRs require a changelog entry via `make changelog` **only for user-facing changes** (new features, bug fixes, behaviour changes visible to CLI users). Internal changes (CI/CD workflow fixes, build scripts, code refactors with no visible effect) do not need changelog entries. Because `neo4j-cli` bundles all child CLIs, user-facing changes to a child require entries for both — use `changie new --projects <child> --projects neo4j-cli --kind <kind> --body <body>` for non-interactive use

## Testing Framework

TESTING FRAMEWORKS: [Go testing, testify, afero (in-memory FS)]

See [`.agents/testing.md`](.agents/testing.md) for full details.

- Tests are colocated with source as `*_test.go` files
- Run with `go test ./...`; CI runs on ubuntu, windows, and macos
- Mock HTTP server and filesystem helpers live in `neo4j-cli/aura/internal/test/testutils/`
- `neo4j-cli/` (the super-CLI package) has no test files; this is a pre-existing gap
- **Prefer table-driven tests** (`for _, tc := range []struct{...}{...}`) when writing new tests — they reduce duplication and make it easy to add cases later

## Architecture

ARCHITECTURE PATTERN: Cobra command tree — one file per leaf command, directory structure mirrors command hierarchy

See [`.agents/architecture.md`](.agents/architecture.md) for full details.

Two binaries are produced:
- **`neo4j-cli`** — super-CLI entrypoint (`neo4j-cli/main.go`); wraps `aura-cli` under the `aura` subcommand
- **`aura-cli`** — standalone Aura CLI (`neo4j-cli/aura/cmd/main.go`)

```
neo4j-cli/
  app/app.go               # neo4j-cli cobra tree builder (NewCmd, Version) — importable
  main.go                  # thin entrypoint; mounts aura subcommand as "aura"
  internal/skill/          # per-binary skill template (bundle, description.txt, additions.md, gen/)
  aura/
    cmd/main.go            # aura-cli standalone entrypoint
    aura.go                # Root command, registers subcommands
    internal/
      api/                 # HTTP client for Neo4j Aura REST API
      flags/               # Custom reusable flag types
      output/              # JSON + table rendering
      skill/               # per-binary skill template (mirrors neo4j-cli/internal/skill)
      subcommands/         # One directory per resource, one file per action
        instance/, tenant/, credential/, config/,
        deployment/, dataapi/graphql/, graphanalytics/,
        import/, customermanagedkey/
common/
  clicfg/                  # Config, credentials, project state (OS-specific paths)
  clierr/                  # Shared error types
  skill/                   # Shared agent-skill logic (catalog, render, installer, cobra wrapper)
```

Agent-skill subsystem: `common/skill/` holds the binary-agnostic logic (agent catalog, path expansion, bundle render, install/remove/list/check, cobra wrapper). Each binary has its own `<bin>/internal/skill/` template (`embed.go` + `description.txt` + `additions.md` + `gen/main.go` + committed `bundle/`). Adding a new standalone CLI = copy the template, edit `description.txt`/`additions.md`/`gen/main.go` import, mount `skill.NewCmd(cfg, binskill.Bundle, "<newcli>")`, run `go generate`. No edits to `common/skill/`. See `CONTRIBUTING.md` "Generated content" for the full workflow.

Key CLI conventions (see `CONTRIBUTING.md`):
- Singular nouns for commands (`instance`, not `instances`)
- `<resource> <action>` form (`instance list`, not `list-instance`)
- One positional argument max; extras become flags
- `--output json|table` for all read commands
- `--await` flag for async operations
- Follow CLI best practices from https://clig.dev/ — source at https://github.com/cli-guidelines/cli-guidelines/blob/main/content/_index.md (fetch the raw markdown for token-efficient reference)

## Deployment

DEPLOYMENT STRATEGY: GitHub Releases via GoReleaser, triggered by `CHANGELOG-neo4j.md` updates on `main`

See [`.agents/deployment.md`](.agents/deployment.md) for full details.

- `changie` batches changelog entries and opens release PRs automatically (dual-project: `aura-cli` + `neo4j-cli`)
- Merging a release PR triggers GoReleaser to publish binaries for linux/windows/darwin (amd64 + arm64)
- macOS binaries are code-signed and notarized
- Each binary gets its own version: `AURA_CLI_VERSION` for `aura-cli`, `GORELEASER_CURRENT_TAG` for `neo4j-cli`
- Combined `release-notes.md` is generated with `## Versions` and `## Changes` sections before GoReleaser runs

## Makefile Notes

- `license-check` target uses `$(GOPATH)/bin/addlicense` (not bare `addlicense`) — GOPATH/bin may not be on PATH
- `license-check` requires a Unix shell (`find` + `xargs`); won't work natively on Windows
- `make generate` runs `go generate ./...`; `make generate-check` runs generate then `git diff --exit-code` (CI gate). Wired in `.github/workflows/test.yml` between Build and Lint, runs on full OS matrix.
- Drift sim: editing a bundle file directly to test generate-check is futile — `go generate` overwrites it. Mutate a real cobra-tree input (e.g. a Short string in `app.go`) to simulate stale-bundle detection.

## Changie Multi-Project Notes

- `ProjectConfig` in changie does NOT support `changesDir` or `changelogPath` fields — only `label`, `key`, `changelog`, and `replacements`
- Version files live at `changesDir/<key>/v*.md` (e.g., `.changes/aura-cli/v1.7.0.md`) — changie appends the project key to `changesDir` automatically
- All projects share a single unreleased directory at `changesDir/unreleasedDir/` (e.g., `.changes/unreleased/`) — change files are tagged with `project:` field inside the YAML, not by directory
- `changie latest --project aura-cli` outputs `aura-cliv1.7.0` (project key prepended with no separator by default) — use `--remove-prefix` to strip `v` but key is always prepended; shell workflows must strip `aura-cli` prefix (e.g., `sed 's/^aura-cli//'`)
- `ProjectsVersionSeparator` in `.changie.yaml` can be set to `-` to get `aura-cli-v1.7.0` instead of `aura-cliv1.7.0`; leave unset (empty) for `aura-cliv1.7.0`
- `changie merge` (no flags) automatically iterates all `projects:` in config and writes each to its own `changelog:` path — confirmed from source (`cmd/merge.go`). Calling `changie merge --project` is not supported by changie's CLI.
- `changie new --projects <a> --projects <b>` creates entries for multiple projects in one call; the interactive prompt (`make changelog`) also supports multi-select
- This repo uses kind labels `Major`, `Minor`, `Patch` (not `added`/`feat`) — check `.changie.yaml` `kinds:` before using `--kind`
- If changie isn't installed locally, hand-author YAML files under `.changes/unreleased/` named `<project>-<Kind>-<YYYYMMDD>-<HHMMSS>.yaml` with fields `project / kind / body / time` (single-quoted body, RFC3339 time). Write one file per project for dual-project entries.

## Changie Workflow Notes

- To detect whether `.changes/unreleased/` contains entries for a given project, use `grep -rl 'project: <key>' .changes/unreleased/ 2>/dev/null | grep -q .` — the `2>/dev/null` handles an absent/empty directory and `grep -q .` converts the file list to a boolean exit code
- Write boolean outputs to `GITHUB_OUTPUT` with `echo "has_<project>=true" >> $GITHUB_OUTPUT` / `false` in an if/else so downstream steps can use `if: steps.<id>.outputs.has_<project> == 'true'`
- Always gate terminal steps (e.g. `create-pull-request`) on the same detection outputs — skipped upstream steps produce empty outputs, not skipped downstream steps; without a guard the terminal step runs with blank inputs and creates a malformed artifact
- For multiline GitHub Actions outputs, use the heredoc form: `{ echo "key<<EOF"; echo "${VALUE}"; echo "EOF"; } >> $GITHUB_OUTPUT` — this avoids issues with inline `|` syntax
- When building multiline strings in `run: |` shell blocks, use `printf '...\n...'` instead of multi-line string assignment — the YAML indentation level (e.g. 10 spaces) carries into continuation lines as literal whitespace

## Release Workflow Notes

- Release workflow triggers on `CHANGELOG-neo4j.md` changes (not `CHANGELOG.md`)
- `AURA_CLI_VERSION` env var set in an earlier step must be re-referenced as `${{ env.AURA_CLI_VERSION }}` in the GoReleaser action's `env:` block — GitHub Actions does not auto-forward env vars set by previous steps into action env blocks
- The neo4j-cli changelog body for a version lives at `.changes/neo4j-cli/<version>.md`; `tail -n +2` strips the `## vX.Y.Z - DATE` header line
- Avoid heredoc indentation issues in `run: |` blocks: use `{ printf ...; } > file` instead of `cat > file << EOF ... EOF` when shell lines are indented under YAML
- Job-level `outputs:` block surfaces step outputs to downstream `workflow_run` consumers. To expose a step output, the step needs an `id:` and must `echo "key=val" >> $GITHUB_OUTPUT` — then reference as `${{ steps.<id>.outputs.<key> }}` in the job `outputs:` block. Output is always populated (downstream gates on the value, not whether it was set).
- All actions in `.github/workflows/` are SHA-pinned with a `# v<major>` trailing comment (e.g. `actions/upload-artifact@ea165f8d... # v4`) — match this convention for any new action; Renovate handles bumps.
- `release.yml`'s `include_neo4j` / `include_aura` are computed by diffing `git diff HEAD~1 --name-only` against changelog filenames — split into its own step (id: `changed`) so the job's `outputs:` block can wire it cleanly to downstream `workflow_run` workflows like `publish-npm.yml`.
- `workflow_run.workflows: ["<name>"]` matches by the upstream workflow's `name:` field, NOT the filename. `release.yml` declares `name: release` (lowercase) so the watcher uses `["release"]`. Mismatching never errors — it just silently never triggers.
- Cross-workflow artifact download with `actions/download-artifact@v4` requires both `github-token: ${{ secrets.GITHUB_TOKEN }}` AND `run-id: ${{ github.event.workflow_run.id }}`. Without `run-id` it looks at the current run only and 404s.
- `workflow_run` events do NOT have `inputs.*` populated; `workflow_dispatch` events do not have `github.event.workflow_run.*`. To pick a value cleanly across both triggers use the ternary pattern `${{ github.event_name == 'workflow_dispatch' && inputs.x || steps.<auto>.outputs.x }}` — short-circuit makes the unset side empty and `||` falls back.
- A `workflow_run`-triggered workflow's job-level `if:` cannot read artifact contents (artifacts haven't been downloaded yet). Gate the JOB on `github.event.workflow_run.conclusion == 'success'`, then download the meta artifact, parse with `jq`, and apply per-step `if:` gates on the parsed value.

## GoReleaser Notes

- GoReleaser v2 deprecates `archives.format` (string) — use `archives.formats` (list)
- GoReleaser v2 deprecates `format_overrides.format` — use `format_overrides.formats`
- Each `archives` entry must have a unique `id`; omitting it defaults to `"default"` and causes errors when there are multiple archive blocks
- Use `{{ .Binary }}` in `name_template` (not `{{ .ProjectName }}`) when building multiple binaries so archives are named per binary
- `-X "<importpath>.Version=..."` ldflag must match the actual package path of the Version var. If you move Version from `package main` to e.g. `neo4j-cli/app`, update the ldflag to `-X "github.com/neo4j/cli/neo4j-cli/app.Version=..."` — a stale path silently no-ops and ships `dev`.

## Repo Doc Notes

- `CLAUDE.md` is a symlink to `AGENTS.md` (`ls -la` confirms). Edit `AGENTS.md` once — both surfaces update. Don't write to `CLAUDE.md` directly.
- Contributor-facing workflows (e.g. `make generate` / add-new-CLI procedure) live in `CONTRIBUTING.md` "Development" subsections. AGENTS.md Architecture orients readers and links to CONTRIBUTING.md for the procedure rather than duplicating it.

## Repo Layout Notes

- `neo4j-cli/app/app.go` builds the neo4j-cli cobra tree and exports `Version`. `neo4j-cli/main.go` is a thin entrypoint. Generators (e.g. skill bundle) import `app` to walk the tree without main-side effects.
- `neo4j-cli/aura/aura.go` already exposes `NewCmd` (super-CLI mount) and `NewStandaloneCmd` (aura-cli binary, adds credential).
- `common/skill/` holds shared agent-skill logic (catalog, path expansion, installer). Hermetic-friendly: `DetectAgents(afero.Fs)` takes an FS; tests use `afero.NewMemMapFs` + `t.Setenv("HOME", ...)`.
- `common/skill/filesystem.go::CopyBundle(dst, dstDir, bundle fs.FS)` walks `bundle` (already scoped — generators do `fs.Sub(Bundle, "bundle")` upstream). Uses `filepath.FromSlash` on each entry so embed.FS forward slashes translate to OS separators on Windows.
- `common/skill/render.Bundle(root, opts)` returns `map[string][]byte` keyed with forward-slash paths (`SKILL.md`, `references/<sub>.md`). Uses `LocalFlags()` (not `Flags()`) when rendering subcommand flag tables to avoid duplicating root persistent flags shown in SKILL.md "Global Flags". Sorts subcommands + flag rows for byte-determinism. TOC inserted only when reference body >100 lines, between the H1 and the rest of the body.
- Render golden-file tests use a `-update` flag (`go test ./common/skill/render -update`) to regenerate `testdata/`. The gate runs without it; pass `-update` only when the renderer's output legitimately changes.
- `common/skill/installer.go` Install/Remove/List/Check: typed sentinel errors (`ErrNoAgentsDetected`, `ErrUnknownAgent`, `ErrAgentNotDetected`) for command-layer assertions. `{{VERSION}}` substituted in SKILL.md only — references stay verbatim. Install RemoveDir's the dst before copy so reinstall doesn't leave stale references. Check returns rows only for installed agents; drift=true also fires on `unknown-version` (frontmatter missing/unparseable). Frontmatter parsed via `regexp` — no YAML dep.
- `common/skill/skill.go::NewCmd(cfg, bundle, skillName)` parent `skill` cobra command. Persistent `--output` flag at the parent (mirrors `tenant.NewCmd`), validated + bound in `PersistentPreRunE` via `cfg.Aura.BindOutput`. Installer sentinel errors are wrapped via `clierr.NewUsageError` with the valid-agent list before returning to cobra. Drift on `check` renders the table/JSON THEN returns a non-nil RunE error so the user still sees the rows. JSON output is a plain array (matches `output.PrintBodyMap` posture — serialize the data directly, no envelope wrapper).
- `common/skill/` follows the dominant `<resource>.go` parent + `<action>.go` per-leaf file convention (matches `aura/internal/subcommands/instance/`, `credential/`, etc.). `skill.go` holds only `NewCmd` + persistent-flag wiring; each leaf (install/remove/list/check) lives in its own `<leaf>.go` + `<leaf>_test.go`. Cross-leaf helpers (`formatAgentErr`, `agentNames`, `printJSON`) live in `helpers.go`. Shared test fixture lives in `cmd_helpers_test.go`. New leaves go in their own file.
- Per-binary skill template lives at `<bin>/internal/skill/`: `embed.go` (`//go:generate go run ./gen` + `//go:embed bundle` into unexported `rawBundle` + `var Bundle fs.FS = mustSub(rawBundle, "bundle")`), `description.txt` (single-line third-person frontmatter description), `additions.md` (gotchas, leading HTML comment), `gen/main.go` (package main, resolves pkg dir via `runtime.Caller(0)` so it's CWD-independent, removes bundle/ before regenerating to evict stale refs), `gen/main_test.go` (round-trip into `t.TempDir`, byte-equal vs committed bundle — local guard that complements `make generate-check`). Adding a new standalone CLI = copy this template + edit description.txt/additions.md + import update in gen/main.go.
- Bootstrap order matters: `embed.go`'s `//go:embed bundle` fails to compile if `bundle/` is missing. Run `go run ./<bin>/internal/skill/gen` first to populate `bundle/`, then everything else (build, tests, `go generate`) works. Subsequent regenerations are fine because gen/ is a sibling package that compiles independently of embed.go.
- `//go:embed bundle` exposes an fs.FS rooted ABOVE the `bundle/` dir — `fs.WalkDir(rawEmbed, ".")` yields `bundle/SKILL.md`, not `SKILL.md`. The installer assumes the flat layout produced by `render.Bundle`, so `Bundle` MUST be sub-rooted: `var Bundle fs.FS = mustSub(rawBundle, "bundle")`. Naïvely exporting `var Bundle embed.FS` writes installs to `<skillsDir>/<skillName>/bundle/SKILL.md` and skips `{{VERSION}}` substitution. Fixture-based unit tests (`fstest.MapFS{"SKILL.md": ...}`) are already flat and miss this — lock the contract per-binary with an `install_e2e_test.go` that drives the REAL exported `Bundle` through `commonskill.Install` against `afero.NewMemMapFs` + `t.Setenv("HOME", t.TempDir())` and asserts (a) `<skillsDir>/<skillName>/SKILL.md` exists, (b) `{{VERSION}}` is substituted, (c) at least one references/*.md is written, plus a `TestBundleWalkAtRoot` contract test that walks `Bundle` and asserts `SKILL.md` appears at the root.
- The aura-cli generator imports `neo4j-cli/aura` and builds `NewStandaloneCmd` (NOT `NewCmd`) — `NewCmd` is the super-CLI mount point and omits `credential`. Wrong choice would mis-represent aura-cli's surface. Generator passes `clicfg.NewConfig(MemMapFs, "dev")`; the literal "dev" never surfaces because render emits `{{VERSION}}` placeholder regardless.
- `cfg.Aura.AuraBetaEnabled()` defaults false on a fresh `MemMapFs` config, so beta-gated commands (dataapi, import, deployment) are intentionally absent from the generated aura-cli bundle. Matches the default-build user surface; users who enable beta locally get a richer `--help` but the shipped skill stays stable.
- Skill cobra mount: top-level in `app.NewCmd` (super-CLI) and inside `aura.NewStandaloneCmd` (aura-cli binary), NEVER inside `aura.NewCmd`. Mounting in `NewCmd` would duplicate skill under the super-CLI's nested `aura` subtree (`neo4j-cli aura skill`). After mounting, re-run `go generate ./...` so each binary's bundle includes its own `references/skill.md`.
- Cobra prints parent help (exit 0) for an unknown subcommand of a parent group with no `RunE` — so the negative test for "no duplicate skill mount" is structural (skill absent from `Available Commands:`), not exit-code-based. Don't write a test that expects non-zero exit.
- Usage guide heading conventions diverge: `docs/usageGuide/Neo4j CLI.md` uses one h1 title + h2/h3 sections; `docs/usageGuide/A Guide To The New Aura CLI.md` uses h1-per-top-level-area + h2 leaves. Match per-file when adding sections; mismatch breaks the file's TOC shape.
- Adding a top-level command to `app.NewCmd` (e.g. `query`) invalidates the committed neo4j-cli skill bundle — `TestGenerator_RoundTrip` in `neo4j-cli/internal/skill/gen` fails until you run `make generate`. Even if your task only "scaffolds" the command, regen the bundle in the same task or `make test` won't be green. Same applies to any new top-level mount on `aura.NewStandaloneCmd` (aura-cli bundle).

## Hermetic Test Notes

- For path-expansion tests using `~` / `$XDG_CONFIG_HOME`, use `t.Setenv("HOME", "...")` and `t.Setenv("XDG_CONFIG_HOME", "")` — Go's `os.Getenv` returns "" for both unset and set-to-empty, and `t.Setenv` auto-restores after the test.
- Use `afero.DirExists` (not `Exists`) for "is the agent installed?" checks — files at the marker path shouldn't count as detected.
- `go-pretty/v6/table` upper-cases header text by default — assertions on table output should compare against `strings.ToLower(...)` for header columns, exact case for body cells.
- Lightweight cobra command tests can wire `clicfg.NewConfig(testfs.GetTestFs(...), version)` directly without the heavier `testutils.NewAuraTestHelper` — the latter pulls in API mocking and credential setup that `skill` doesn't need.
- For repo-wide gate tests that must auto-discover content (e.g. `common/skill/bundles_test.go` walking every `<bin>/internal/skill/bundle/SKILL.md`), resolve repo root via `runtime.Caller(0)` then `filepath.Walk` from there. Suffix-match paths after `filepath.ToSlash` so Windows runs match. Prune `.git`, `node_modules`, `bin`, `.changes` to keep the walk fast.

## Windows CI Gotchas

- Path-separator bugs in `expandPath`-style helpers are Windows-only. Catalog entries keep forward slashes (portable convention); helpers MUST wrap any post-substitution path through `filepath.FromSlash` (or build via `filepath.Join`) so the whole path is OS-native. A `ReplaceAll(path, "$XDG_CONFIG_HOME", xdg)` where `xdg` came from `os.Getenv` produces mixed separators on Windows (`C:\…\.config/opencode`) — fix at the helper, not the catalog.
- Test expected values that hard-code separators bake in OS assumptions. Build expected values with `filepath.Join` / `filepath.FromSlash` rather than literals when asserting cross-OS path output. MemMapFs marker paths in detection tests must also be built OS-natively so they match what the (post-fix) helper looks up.
- Committed `.md` / golden / bundle files MUST be pinned to LF via `.gitattributes` — Windows runners have `core.autocrlf=true` by default and will rewrite to CRLF on checkout. The renderer (`common/skill/render`) and `make generate-check` both assume LF; a CRLF checkout breaks byte-equal golden tests AND `git diff --exit-code`. The repo-root `.gitattributes` covers `common/skill/render/testdata/**`, `**/internal/skill/bundle/**`, `**/internal/skill/additions.md`, `**/internal/skill/description.txt`. `common/skill/bundles_test.go::TestCommittedBundlesAndTestdataAreLF` is the assertion that catches a weakened/removed attribute.

## npm Distribution Notes

- `distribution/npm/cli/bin/neo4j-cli.js` is the wrapper-package bin shim — Node stdlib only, no runtime npm deps. Resolves `@neo4j-labs/cli-<platform>-<arch>` via `require.resolve` and execs via `spawnSync(..., { stdio: 'inherit' })` so SIGINT and TTY pass through. Use `result.status ?? 1` for exit code (covers null from signal-kill AND undefined from spawn error).
- The SUPPORTED platform list is hard-coded in the shim itself (not read from `distribution/platforms.tsv`) because the shim ships standalone inside the `@neo4j-labs/cli` tarball — no extra runtime files. Adding a platform = update the shim's SUPPORTED array + add a `platforms.tsv` row + add a GoReleaser build entry (reviewers should catch all three together).
- `result.error` branch is distinct from non-zero exit — handles spawn-time failures (missing exe perms, ENOENT) before reaching the exit-code propagation path.
- `.gitignore` ignores `bin/` globally for Go build output. The `distribution/**/bin/` paths (npm wrapper shim, future platform-pkg bin layouts) need an explicit `!distribution/**/bin/` un-ignore — without it, `git status` silently hides committed-intent files and `git add` errors with "ignored by .gitignore". Verify with `git check-ignore -v <path>`.
- `package.json.tmpl` files use `__TOKEN__` placeholders (double underscores), NOT `{{TOKEN}}` or `${TOKEN}` — the former survives shell/sed metachar quoting cleanly and the templates parse as valid JSON pre-substitution (so editors lint them). Wrapper template enumerates all 8 platform pkgs in `optionalDependencies` inline, pinned EXACTLY to `__VERSION__` (no `^`) — caret would let npm pick a newer platform pkg from a different release. Verify any template change with `sed ... | jq .` to catch trailing-comma / quoting drift before commit.
- `distribution/npm/publish.sh` builds rendered packages under `dist/.npm-build/` (already-gitignored via `dist/`). Reads `platforms.tsv` via `while IFS=$'\t' read ... done < <(tail -n +2 ...)` process substitution so `set -e` doesn't choke on the header-skip pipeline. Uses bash parameter expansion `${DIRNAME_TMPL//\$\{VERSION\}/$VERSION}` (NOT sed) to substitute `${VERSION}` in the goreleaser dirname column — keeps it self-contained and avoids escaping the literal `$`. The `published()` helper always returns 0 to satisfy `set -e`; the publish helper does the conditional `[skip]` log. `TAG_FLAG`/`DRY_RUN` are intentionally word-split (with shellcheck disable comment) so the empty case omits the flag entirely.
- Wrapper README.md ship in publish.sh is gated on file existence (`[ -f ... ] && cp`) so `publish.sh` doesn't depend on the user-facing README task — once that file exists, it auto-ships in the wrapper tarball. Same pattern works for any future ship-in-tarball additions.
- Test publish.sh locally without hitting npm by stubbing `npm` (case on `$1`: `view` echo empty/version, `publish` echo args) and pointing `PATH` at the stub dir. Build a fake `dist/` tree from `platforms.tsv` mirroring GoReleaser's archive layout. Verifies dist-tag selection (alpha/beta/rc/next/latest), idempotency, ordering, and rendered package.json validity in seconds, no registry round-trip.
- npm scope is `@neo4j-labs/cli` (not `@neo4j/cli`) — matches the GitHub `neo4j-labs` org. Bin-shim `pkgName` const, both `package.json.tmpl` files, publish.sh `PKG_NAME` + final wrapper publish, and all docs pin this. The `homepage`/`repository`/`bugs` URLs in the package.json templates still point to `github.com/neo4j/cli` (legacy stub) — cleanup pending; not the npm scope.
- ASCII diagrams in `distribution/npm/README.md` use 22-char inner box width to fit `@neo4j-labs/cli-<arch>` content. Arrow-to-box-center alignment matters for readability: with three boxes at columns 0-23, 34-57, 68-91 the centers are 11.5/45.5/79.5; place arrows at integer cols 11/45/79. Verify with `python3 -c "for i,c in enumerate(line): print(i,c)"` after edits.

## golangci-lint Notes

- Version installed: v2.11.4 (via Homebrew)
- golangci-lint v2 requires `version: "2"` at the top of `.golangci.yml`
- In v2, `gofmt` is a **formatter** (not a linter); put it under `formatters.enable`, not `linters.enable`
- Use `linters.default: none` to disable auto-enabled defaults (e.g. `ineffassign`) and run only explicitly listed linters
- Config lives at `.golangci.yml` in repo root
- In CI, `golangci/golangci-lint-action@v6` is used as the lint step — it installs, caches, and runs golangci-lint using `.golangci.yml`. This is equivalent to `make lint`. Renovate will pin the SHA.

## Cobra Flag Access Notes

- `cmd.Flags().GetString("foo")` only sees LOCAL flags until `mergePersistentFlags()` runs (during `Execute` or `ParseFlags`). Calling it from a unit test that drives a function directly (without Execute) will fail with `flag accessed but not defined`. Use `cmd.Flag("foo").Value.String()` instead — `Flag()` falls through to persistent flags + parents' persistent flags via `persistentFlag()`/`updateParentsPflags()`. Same applies to GetBool — for bool defaults that overlap with "unset" (e.g. `--insecure` defaults false), gate on `cmd.Flag("name").Changed` to disambiguate.
- For first-non-empty-wins precedence pass values to a helper that picks the FIRST non-empty. For lowest→highest precedence (`.env` < env < flag) use a `last-non-empty-wins` helper instead — easy to swap accidentally; query/connect.go calls this `overlay()`.

## Query Command Notes

- `neo4j-cli/query/connect.go` exposes `httpDoer interface { Do(*http.Request) (*http.Response, error) }` as the test seam — tests inject `srv.Client()` from `httptest.NewServer`/`NewTLSServer` rather than monkey-patching.
- `loadEnvFile(fs afero.Fs, explicitPath, startDir string) (map[string]string, error)` is the pure helper for `.env` resolution; `resolveConn` glues it to `os.Getwd()` + `cfg.Aura.Fs()` in production. Walk-up uses `filepath.Dir` and stops when `parent == dir` (root reached).
- TLS unit tests: `httptest.NewTLSServer` produces a self-signed cert. The default `newHTTPClient(false)` rejects it; `newHTTPClient(true)` (i.e. `--insecure`) accepts. `srv.Client()` IS pre-configured to trust httptest's cert — for the secure-rejection assertion you MUST construct your own client (`newHTTPClient(false)`), NOT use `srv.Client()`.
- `subosito/gotenv` was already a transitive dep (via viper); `go mod tidy` after first import promotes it to direct. No new third-party deps for the query command.
- `t.Chdir(t.TempDir())` (Go 1.24+) is the hermetic primitive for tests that need a controlled cwd — go.mod `go 1.25.0` baseline allows it.
- `clicfg.AuraConfig` exposes `Output()` getter + `BindOutput(*pflag.Flag)` viper bind — there is NO `SetOutput` setter. Tests that need a specific output mode write a real `config.json` via `testfs.GetTestFs(`{"aura":{"output":"json"}}`, "{}")` (mirrors `common/skill/cmd_helpers_test.go`); the in-memory FS makes this hermetic.
- `neo4j-cli/query/output.go::renderRows` keeps strings unquoted in table cells (readability) but JSON-marshals everything else; `nil` → literal `null`. `jsonRowsResult` struct pins envelope field order (columns, rows, truncated) via Go struct declaration order — `encoding/json` honours it, no ordered-map dep needed. Pair with `rowsFromValues(columns, [][]any)` to convert API positional values into `[]map[string]any` (missing positions → nil, extras dropped).
- `neo4j-cli/query/run.go` exposes three package-level `var` test seams (`passwordReader`, `stdinIsTTY`, `stdinReader`) backed by `golang.org/x/term` + `os.Stdin`. The harness in `run_test.go` resets all three via `t.Cleanup` so production behaviour returns between tests. `runQuery` sets `cmd.SilenceUsage = true` at entry to suppress cobra's "Usage:" reprint after errors. Truncation order: `truncateArrays` runs on positional `[][]any` BEFORE `capRows` so array-truncation is unaffected by `--max-rows`. Row-cap warning is `fmt.Fprintf(cmd.ErrOrStderr(), "warning: truncated to %d rows (use --max-rows 0 for unlimited)\n", N)` — fires regardless of output mode (JSON consumers still get truncated:true in the envelope AND the stderr warning).
- Tests that drive cobra Execute() must call `cmd.SetContext(context.Background())` because `runStatement` reads `cmd.Context()` and `http.NewRequestWithContext` rejects nil ctx. Cobra leaves ctx nil unless the caller sets it.
- Password-prompt non-TTY error message MUST mention all three of `--password`, `NEO4J_PASSWORD`, and `.env` to satisfy the task-006 acceptance criterion. The "Password:" prompt itself goes to `cmd.ErrOrStderr()` (NOT stdout) so it never pollutes a JSON envelope on a piped stdout.
- `golang.org/x/term` is a separate module from `golang.org/x/sys`; `go get golang.org/x/term` adds it. `go mod tidy` may demote it back to indirect until a Go file imports it; importing in `run.go` then re-running `go mod tidy` lands it as direct.
- `neo4j-cli/query/schema.go` `:schema` cobra leaf — cobra accepts the literal `:schema` as `Use` (verified: dispatches normally as `query :schema`). RunE reuses `resolveConn` + `promptPassword` from run.go. Required queries (1–5) propagate errors; optional `CALL dbms.components()` + `SHOW SETTINGS YIELD ... 'db.query.default_language'` swallow errors via `fetchDatabaseInfo` returning nil-or-partial. Per-relType MATCH builds the pattern with `fmt.Sprintf("MATCH (n)-[r:`%s`]->(m)...", stripped)` after `stripRelTypeWrap` peels colons/backticks/quotes (driver returns ":`OWNS`"); `uniqueRelTypes` dedups + sorts for deterministic call order.
- `:schema` table mode emits FIVE sub-tables (Nodes, Relationships, Relationship Paths, Indexes, Constraints) per PRD REQ-F-013 — Database is JSON-only. JSON envelope marshals the `schemaResult` struct directly so field order is fixed by struct declaration. `--max-rows`/`--truncate-arrays-over` do NOT apply to `:schema`.
- Routing httptest pattern for multi-statement tests: one `httptest.NewServer` whose handler reads `req["statement"]` and substring-matches against a `map[string]cannedResponse` route table. `errBody` field on the route triggers a 4xx + `errors[]` envelope so the existing `runStatement` error path exercises end-to-end. Required-query failure tests inject a route override on the happy-path map; the test stays one assertion away from the canned setup.
- Avoid using `any` as a local variable name even though syntactically valid post-Go 1.18 — Go's `any` is an alias for `interface{}`, shadowing it confuses readers. Rename to `gotAny`/`hasAny`.
- Raw string literals (backticks) cannot contain backticks; use `"..." + "..."` regular-string concatenation when test fixtures need to embed Cypher backtick-quoting (e.g. `:"`ACTED_IN`"`). gofmt detects this as a parse error before tests can even run.
- `make generate-check` diffs the WHOLE repo, not just bundle paths — any unrelated edit (e.g. an in-progress task YAML status flip) makes it fail with "ERROR: generated files are stale" even when the bundle is clean. To verify the bundle in isolation, `git stash push -- <unrelated-files>` first, run the check, then `git stash pop`. Don't mistake the false positive for real drift.
- Bundle regen happens organically in feature-add tasks (because new cobra subcommands appear in the tree the moment they're mounted), so a dedicated "regen bundle" verification task is typically a no-op gate. Don't expect a bundle diff if the prior task already mounted the command and ran `make generate`.
- Hand-authored changie YAMLs in `.changes/unreleased/` use single-quoted `body:`. To embed an apostrophe (e.g. `'query'`), double it: `body: 'Add ''query'' command ...'` — YAML's only single-quote escape rule. `time:` accepts RFC3339 with `Z` suffix or numeric offset; existing entries mix both styles.

---

_This AGENTS.md was generated using agent-based project discovery._

<!-- END GENERATED: AGENTS-MD -->

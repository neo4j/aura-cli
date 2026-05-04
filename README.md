# Neo4j CLI

## Installation

### Prebuilt archives (recommended)

Download the archive for your OS/arch from the [releases](https://github.com/neo4j-labs/neo4j-cli/releases/latest) page and extract the `neo4j-cli` binary.

### npm

```bash
npm i -g @neo4j-labs/cli
```

Installs a small Node wrapper plus a prebuilt `neo4j-cli` binary matching your OS and CPU. `pnpm add -g @neo4j-labs/cli` and `yarn global add @neo4j-labs/cli` work the same way. Prereleases via `npm i -g @neo4j-labs/cli@alpha` (also `@beta`, `@rc`, `@next`). See [`distribution/npm/cli/README.md`](./distribution/npm/cli/README.md) for the supported platform matrix.

## Usage

Extract the executable to a directory of your choosing.

Create Aura API Credentials in your [Account Settings](https://console.neo4j.io/#account), and note down the client ID and secret.

Add these credentials into the CLI with a name of your choosing:

```bash
./neo4j-cli aura credential add --name "Aura API Credentials" --client-id <client-id> --client-secret <client-secret>
```

This will add and set the credential as the default credential for use.

You can then, for example, list your instances in a table format:

```bash
./neo4j-cli aura instance list --output table
```

If you would rather just type `neo4j-cli` then move the neo4j-cli binary into the file path of your computer.  
Windows:

```bash
move neo4j-cli.exe c:\windows\system32
```

Mac:

```bash
sudo mv neo4j-cli /usr/local/bin
```

To see all of the available commands:

```bash
./neo4j-cli
```

Help for each command is accessed by using it without any flags or options. For example, to see help creating an instance:

```bash
./neo4j-cli aura instance create
```

## Agent skills

`neo4j-cli` ships an embedded skill bundle (`SKILL.md` + per-subcommand references) that teaches AI coding agents how to drive the CLI. `skill install` drops that bundle into the supported agents' skill directories so the agent picks it up on next run.

Supported agents:

-   Claude Code
-   Cursor
-   Windsurf
-   Copilot
-   Gemini CLI
-   Cline
-   Codex
-   Pi
-   OpenCode
-   Junie

Install into every detected agent (or pass an agent name to target one):

```bash
./neo4j-cli skill install
./neo4j-cli skill install claude-code
```

List supported agents and per-agent install state:

```bash
./neo4j-cli skill list
```

Check installed bundles for version drift against the running binary:

```bash
./neo4j-cli skill check
```

Remove the installed bundle (idempotent):

```bash
./neo4j-cli skill remove
./neo4j-cli skill remove claude-code
```

The same `skill` subcommand is available on `aura-cli`, scoped to the standalone Aura surface.

Beta features (commands gated by `AuraBetaEnabled`, e.g. `dataapi`, `import`, `deployment`) are not included in the generated bundle — the bundle reflects the default-config command surface.

## Feedback / Issues

Please use [GitHub issues](https://github.com/neo4j-labs/neo4j-cli/issues) to provide feedback and report any issues that you have encountered.

## Building locally

Clone the repository and run:

```bash
make build
```

This produces `bin/aura-cli` and `bin/neo4j-cli`. To run without building:

```bash
make run-aura   # standalone aura-cli
make run-neo4j  # neo4j-cli super CLI
```

## Developing and contributing

Read [CONTRIBUTING.md](./CONTRIBUTING.md)

# Contributing

Thanks for your interest in contributing to the Neo4j Aura CLI, [issues](https://github.com/neo4j/aura-cli/issues) and [pull requests](https://github.com/neo4j/aura-cli/pulls) are welcome.

If you want to contribute code, make sure to [sign the CLA](https://neo4j.com/developer/contributing-code/#sign-cla).

## Development

### Testing

The full suite of tests can be run using the following command:

```bash
go test ./...
```

### Local running

The CLI can be run locally without building by running the following command:

```bash
go run neo4j-cli/main.go aura-cli
```

### Pull requests

As well as your code changes, pull requests need a changelog entry. These are added using the tool [`changie`](https://changie.dev/). You will need to install this using the following command:

```bash
go install github.com/miniscruff/changie@latest
```

With this installed, the following command will guide through the process of adding a changelog entry:

```bash
changie new
```

Simply commit the file that this command produces and you're done!

If changie is not available, you may need to add /go/bin to your path: `export PATH="$HOME/go/bin:$PATH"`

### Building

Builds for releases are handled in GitHub Actions. If you want to create local builds, there are a couple of approaches.

To create a simply binary using `go` directly, you can execute the following command:

```bash
go build -o bin/ ./...
```

If you want to build binaries for all varieties of platforms, you can do so with the following command:

```bash
GORELEASER_CURRENT_TAG=dev goreleaser release --snapshot --clean
```

In the above command, `GORELEASER_CURRENT_TAG` can be substituted for any version of your choosing.

## CLI Guidelines

The Aura CLI aims to provide a consistent and reliable experience to the end user. Any change made to the CLI must comform to the following guidelines:

-   All commands must be singular
    -   ✅ `aura-cli instance`
    -   ❌ `aura-cli instances`
-   Output should support the following options:
    -   `json`: Provides the raw JSON output of the API, formatted to be human-readable.
    -   `table`: Provides a subset of the output, formatted to be human readable on a table. Try to keep the table output below 120 characters to avoid overflowing the screen.
-   Verbs and nouns should be separate, with the action at the end
    -   ✅ `aura-cli instance list`
    -   ❌ `aura-cli list-instance`
    -   ❌ `aura-cli list instance`
-   Only one argument should be used, if more than one is needed, use flags instead. This is to avoid confusion when passing parameters without enough context
    -   ✅ `aura-cli instance get <id>`
    -   ❌ `aura-cli instance get <id> <deployment-id>`
    -   ✅ `aura-cli instance get <id> --deployment-id <deployment-id>`
    -   ⚠️ `aura-cli instance get --instance-id <id> --deployment-id <deployment-id>`  
        This valid, but the option above is preferred as it is more concise
-   The argument must always refer to the closest noun
    -   ❌ `aura-cli instance snapshot list <instance-id>`
    -   ✅ `aura-cli instance snapshot list --instance-id <instance-id>`
-   No arguments between commands
    -   ❌ `aura-cli tenant <tenant-id> instance get <id>`
    -   ✅ `aura-cli instance get <id> --tenant-id <tenant-id>`
-   Flags, if set, take precedence over global configuration or default values

> These guidelines are based on https://clig.dev

### Structure

Aura CLI is divided in top level commands, for example:

-   `instance`
-   `config`

Each of these commands handle a certain resource of the API and have several subcommands for the actions, for example:

-   `instance list`
-   `instance get`

Nested subcommands are also allowed, for example:

-   `instance snapshot list`

### Common subcommands

Most commands targetting API resources contain some of the following subcommands as actions:

-   `get`
-   `list`
-   `delete`
-   `create`

Commands may also have some extra, specific commands, such as `instance pause`.

For asynchronous operations (i.e. operations that trigger a job that won't be finished in the same request), the flag `--await` can be used to wait until the operation has been completed, generally polling for the status. If this flag is not set, all operations must finish when the request has been completed, even if a job is pending.

## Resources

-   [CLI Usage Guide](./docs/usageGuide/A%20Guide%20To%20The%20New%20Aura%20CLI.md).
-   [Neo4j Aura API](https://neo4j.com/docs/aura/platform/api/specification/)
-   https://clig.dev

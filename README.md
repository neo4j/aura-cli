# Neo4j CLI

## Installation

Downloadable binaries are available from the [releases](https://github.com/neo4j/cli/releases/latest) page.

Download the appropriate archive for for operating system and architecture.

## Usage

Extract the executable to a directory of your choosing.

Create Aura API Credentials in your [Account Settings](https://console.neo4j.io/#account), and note down the client ID and secret.

Add these credentials into the CLI with a name of your choosing:

```bash
./aura-cli credential add --name "Aura API Credentials" --client-id <client-id> --client-secret <client-secret>
```

This will add and set the credential as the default credential for use.

You can then, for example, list your instances in a table format:

```bash
./aura-cli instance list --output table
```

If you would rather just type ```aura-cli``` then move the aura-cli binary into the file path of your computer.   
Windows:
```bash
move aura-cli c:\windows\system32
```
Mac:
```bash
sudo mv aura-cli /usr/local/bin
```

To see all of the available commands:
```bash
./aura-cli
```
Help for each command is accessed by using it without any flags or options.  For example, to see help creating an instance
```bash
./aura-cli instance create
```

## Feedback / Issues
Please use [GitHub issues](https://github.com/neo4j/aura-cli/issues) to provide feedback and report any issues that you have encountered.

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

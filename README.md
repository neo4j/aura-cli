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

If you would rather just type `aura-cli` then move the aura-cli binary into the file path of your computer.  
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

Help for each command is accessed by using it without any flags or options. For example, to see help creating an instance

```bash
./aura-cli instance create
```

## Feedback / Issues

Please use [GitHub issues](https://github.com/neo4j/aura-cli/issues) to provide feedback and report any issues that you have encountered.

## Developing and contributing

Read [CONTRIBUTING.md](./CONTRIBUTING.md)

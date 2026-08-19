// Copyright (c) "Neo4j"
// Neo4j Sweden AB [http://neo4j.com]

package redact

// capturedArgs stores the redacted command-line arguments.
// This is set by main() at the top of each CLI entrypoint and read by error/panic handlers.
var capturedArgs []string

func SetCapturedArgs(args []string) {
	capturedArgs = args
}

func CapturedArgs() []string {
	return capturedArgs
}

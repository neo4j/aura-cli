// Copyright (c) "Neo4j"
// Neo4j Sweden AB [http://neo4j.com]

package redact

// capturedArgs stores the redacted command-line arguments.
// This is set by main() at the top of each CLI entrypoint and read by error/panic handlers.
var capturedArgs []string

// SetCapturedArgs stores the redacted command-line arguments for access by error handlers.
func SetCapturedArgs(args []string) {
	capturedArgs = args
}

// CapturedArgs returns the redacted command-line arguments previously set by SetCapturedArgs,
// or nil if not yet set.
func CapturedArgs() []string {
	return capturedArgs
}

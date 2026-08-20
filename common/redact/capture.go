// Copyright (c) "Neo4j"
// Neo4j Sweden AB [http://neo4j.com]

package redact

var capturedArgs []string

func CaptureArgs(args []string) {
	capturedArgs = MaskArgs(args)
}

func CapturedArgs() []string {
	return capturedArgs
}

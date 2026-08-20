// Copyright (c) "Neo4j"
// Neo4j Sweden AB [http://neo4j.com]

package main

import (
	"fmt"
	"os"

	"github.com/neo4j/cli/common/clicfg"
	"github.com/neo4j/cli/common/redact"
	"github.com/neo4j/cli/neo4j-cli/aura"
	"github.com/spf13/afero"
)

var Version = "dev"

func main() {
	redact.CaptureArgs(os.Args[1:])

	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Unexpected error running CLI with args %s, please report an issue in https://github.com/neo4j/aura-cli\n\n", redact.CapturedArgs())

			panic(r)
		}
	}()

	cfg := clicfg.NewConfig(afero.NewOsFs(), Version)

	cmd := aura.NewCmd(cfg)
	cmd.SetOut(os.Stdout)
	cmd.SetErr(os.Stderr)
	cmd.Execute()
}

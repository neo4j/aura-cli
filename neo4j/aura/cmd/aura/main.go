/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/neo4j/cli/common/clicfg"
	"github.com/neo4j/cli/common/clictx"
	"github.com/neo4j/cli/neo4j/aura"
	"github.com/spf13/afero"
)

var Version = "dev"

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Unexpected error running CLI with args %s, please report an issue in https://github.com/neo4j/cli\n\n", os.Args[1:])

			panic(r)
		}
	}()

	cmd := aura.NewCmd()

	cmd.SetOut(os.Stdout)
	cfg, err := clicfg.NewConfig(afero.NewOsFs())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	ctx, err := clictx.NewContext(context.Background(), cfg, Version)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cmd.ExecuteContext(ctx)
}

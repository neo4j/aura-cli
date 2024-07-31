/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/neo4j/cli/pkg/clicfg"
	"github.com/neo4j/cli/pkg/clictx"
	"github.com/neo4j/cli/pkg/neo4j"
	"github.com/spf13/afero"
)

var Version = "dev"

func main() {
	cmd := neo4j.NewCmd()

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

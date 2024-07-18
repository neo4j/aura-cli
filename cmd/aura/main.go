/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/neo4j/cli/pkg/aura"
	"github.com/neo4j/cli/pkg/clicfg"
	"github.com/neo4j/cli/pkg/clictx"
)

var Version = "dev"

func main() {
	cfg, err := clicfg.NewConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	ctx, err := clictx.NewContext(context.Background(), cfg, Version)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	aura.Execute(ctx)
}

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

	neo4j.Execute(ctx)
}

/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package neo4j

import (
	"context"

	"github.com/neo4j/cli/pkg/aura"
	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "neo4j",
	Short: "Allows you to manage Neo4j resources",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(ctx context.Context) error {
	return RootCmd.ExecuteContext(ctx)
}

func init() {
	RootCmd.AddCommand(aura.Cmd)
}

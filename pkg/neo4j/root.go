/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package neo4j

import (
	"github.com/neo4j/cli/pkg/aura"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "neo4j",
		Short: "Allows you to manage Neo4j resources",
	}

	cmd.AddCommand(aura.NewCmd())
	return cmd
}

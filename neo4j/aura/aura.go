/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package aura

import (
	"github.com/neo4j/cli/neo4j/aura/internal/subcommands/config"
	"github.com/neo4j/cli/neo4j/aura/internal/subcommands/credential"
	"github.com/neo4j/cli/neo4j/aura/internal/subcommands/customermanagedkey"
	"github.com/neo4j/cli/neo4j/aura/internal/subcommands/instance"
	"github.com/neo4j/cli/neo4j/aura/internal/subcommands/tenant"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "aura",
		Short: "Allows you to programmatically provision and manage your Aura instances",
	}

	cmd.AddCommand(config.NewCmd())
	cmd.AddCommand(credential.NewCmd())
	cmd.AddCommand(customermanagedkey.NewCmd())
	cmd.AddCommand(instance.NewCmd())
	cmd.AddCommand(tenant.NewCmd())

	return cmd
}

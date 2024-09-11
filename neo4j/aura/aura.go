/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package aura

import (
	"github.com/neo4j/cli/common/clicfg"
	"github.com/neo4j/cli/neo4j/aura/internal/subcommands/config"
	"github.com/neo4j/cli/neo4j/aura/internal/subcommands/credential"
	"github.com/neo4j/cli/neo4j/aura/internal/subcommands/customermanagedkey"
	"github.com/neo4j/cli/neo4j/aura/internal/subcommands/instance"
	"github.com/neo4j/cli/neo4j/aura/internal/subcommands/tenant"
	"github.com/spf13/cobra"
)

func NewCmd(cfg *clicfg.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "aura",
		Short:   "Allows you to programmatically provision and manage your Aura instances",
		Version: cfg.Version,
	}

	cmd.AddCommand(config.NewCmd(cfg))
	cmd.AddCommand(credential.NewCmd(cfg))
	cmd.AddCommand(customermanagedkey.NewCmd(cfg))
	cmd.AddCommand(instance.NewCmd(cfg))
	cmd.AddCommand(tenant.NewCmd(cfg))

	return cmd
}

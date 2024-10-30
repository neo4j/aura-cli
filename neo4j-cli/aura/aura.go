/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package aura

import (
	"github.com/spf13/cobra"

	"github.com/neo4j/cli/common/clicfg"
	"github.com/neo4j/cli/neo4j-cli/aura/internal/subcommands/config"
	"github.com/neo4j/cli/neo4j-cli/aura/internal/subcommands/credential"
	"github.com/neo4j/cli/neo4j-cli/aura/internal/subcommands/customermanagedkey"
	"github.com/neo4j/cli/neo4j-cli/aura/internal/subcommands/dataapi"
	"github.com/neo4j/cli/neo4j-cli/aura/internal/subcommands/instance"
	"github.com/neo4j/cli/neo4j-cli/aura/internal/subcommands/tenant"
)

func NewCmd(cfg *clicfg.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "aura",
		Short:   "Allows you to programmatically provision and manage your Aura resources",
		Version: cfg.Version,
	}

	cmd.AddCommand(config.NewCmd(cfg))
	cmd.AddCommand(credential.NewCmd(cfg))
	cmd.AddCommand(customermanagedkey.NewCmd(cfg))
	cmd.AddCommand(instance.NewCmd(cfg))
	cmd.AddCommand(tenant.NewCmd(cfg))
	if cfg.Aura.AuraBetaEnabled() {
		cmd.AddCommand(dataapi.NewCmd(cfg))
	}

	return cmd
}

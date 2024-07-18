/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package aura

import (
	"context"

	"github.com/neo4j/cli/pkg/aura/cmd/config"
	"github.com/neo4j/cli/pkg/aura/cmd/credential"
	"github.com/neo4j/cli/pkg/aura/cmd/customermanagedkey"
	"github.com/neo4j/cli/pkg/aura/cmd/instance"
	"github.com/neo4j/cli/pkg/aura/cmd/tenant"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "aura",
	Short: "Allows you to programmatically provision and manage your Aura instances",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(ctx context.Context) error {
	return Cmd.ExecuteContext(ctx)
}

func init() {
	Cmd.AddCommand(config.Cmd)
	Cmd.AddCommand(credential.Cmd)
	Cmd.AddCommand(customermanagedkey.Cmd)
	Cmd.AddCommand(instance.Cmd)
	Cmd.AddCommand(tenant.Cmd)
}

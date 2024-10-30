package graphql

import (
	"github.com/neo4j/cli/common/clicfg"
	"github.com/spf13/cobra"

	"github.com/neo4j/cli/neo4j-cli/aura/internal/subcommands/dataapi/graphql/authprovider"
)

func NewCmd(cfg *clicfg.Config) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "graphql",
		Short: "Allows you to programmatically provision and manage your GraphQL Data APIs",
	}

	cmd.AddCommand(authprovider.NewCmd(cfg))
	cmd.AddCommand(NewListCmd(cfg))
	cmd.AddCommand(NewGetCmd(cfg))
	cmd.AddCommand(NewCreateCmd(cfg))
	cmd.AddCommand(NewDeleteCmd(cfg))
	cmd.AddCommand(NewResumeCmd(cfg))
	cmd.AddCommand(NewPauseCmd(cfg))

	return cmd
}

package graphql

import (
	"github.com/neo4j/cli/common/clicfg"
	"github.com/spf13/cobra"
)

func NewCmd(cfg *clicfg.Config) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "graphql",
		Short: "Allows you to programmatically provision and manage your GraphQL Data APIs",
	}

	cmd.AddCommand(NewListCmd(cfg))
	cmd.AddCommand(NewGetCmd(cfg))
	cmd.AddCommand(NewCreateCmd(cfg))
	cmd.AddCommand(NewDeleteCmd(cfg))
	cmd.AddCommand(NewResumeCmd(cfg))
	cmd.AddCommand(NewPauseCmd(cfg))

	return cmd
}

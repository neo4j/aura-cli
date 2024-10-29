package dataapi

import (
	"github.com/neo4j/cli/common/clicfg"
	"github.com/spf13/cobra"

	"github.com/neo4j/cli/neo4j/aura/internal/subcommands/dataapi/graphql"
)

func NewCmd(cfg *clicfg.Config) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "data-api",
		Short: "Allows you to programmatically provision and manage your Data APIs",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			cfg.Aura.BindBaseUrl(cmd.Flags().Lookup("base-url"))
			cfg.Aura.BindAuthUrl(cmd.Flags().Lookup("auth-url"))
			cfg.Aura.BindOutput(cmd.Flags().Lookup("output"))

			return nil
		},
	}

	cmd.AddCommand(graphql.NewCmd(cfg))

	cmd.PersistentFlags().String("auth-url", "", "")
	cmd.PersistentFlags().String("base-url", "", "")
	cmd.PersistentFlags().String("output", "", "")

	return cmd
}

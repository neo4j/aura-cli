package authprovider

import (
	"github.com/neo4j/cli/common/clicfg"
	"github.com/spf13/cobra"
)

func NewCmd(cfg *clicfg.Config) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "auth-provider",
		Short: "Allows you to programmatically manage authentication providers for a specific GraphQL Data API",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := cfg.Aura.BindBaseUrl(cmd.Flags().Lookup("base-url")); err != nil {
				return err
			}
			if err := cfg.Aura.BindAuthUrl(cmd.Flags().Lookup("auth-url")); err != nil {
				return err
			}
			if err := cfg.Aura.BindOutput(cmd.Flags().Lookup("output")); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.AddCommand(NewListCmd(cfg))
	cmd.AddCommand(NewGetCmd(cfg))

	cmd.PersistentFlags().String("auth-url", "", "")
	cmd.PersistentFlags().String("base-url", "", "")
	cmd.PersistentFlags().String("output", "", "")

	return cmd
}

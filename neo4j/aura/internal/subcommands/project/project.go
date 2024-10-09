package project

import (
	"github.com/spf13/cobra"

	"github.com/neo4j/cli/common/clicfg"
)

func NewCmd(cfg *clicfg.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "project",
		Short: "Relates to an Aura Project",
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

	cmd.PersistentFlags().String("auth-url", "", "")
	cmd.PersistentFlags().String("base-url", "", "")
	cmd.PersistentFlags().String("output", "", "")

	cmd.AddCommand(NewGetCmd(cfg))
	cmd.AddCommand(NewListCmd(cfg))

	return cmd
}

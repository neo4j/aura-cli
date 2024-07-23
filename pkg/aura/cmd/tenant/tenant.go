package tenant

import (
	"errors"

	"github.com/neo4j/cli/pkg/clictx"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tenant",
		Short: "Relates to an Aura Tenant",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			config, ok := clictx.Config(cmd.Context())

			if !ok {
				return errors.New("error fetching cli configuration values")
			}

			if err := config.BindPFlag("aura.base-url", cmd.Flags().Lookup("base-url")); err != nil {
				return err
			}
			if err := config.BindPFlag("aura.auth-url", cmd.Flags().Lookup("auth-url")); err != nil {
				return err
			}
			if err := config.BindPFlag("aura.output", cmd.Flags().Lookup("output")); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.PersistentFlags().String("auth-url", "", "")
	cmd.PersistentFlags().String("base-url", "", "")
	cmd.PersistentFlags().String("output", "", "")

	cmd.AddCommand(NewGetCmd())
	cmd.AddCommand(NewListCmd())

	return cmd
}

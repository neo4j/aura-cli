package customermetricsintegration

import (
	"github.com/spf13/cobra"

	"github.com/neo4j/cli/common/clicfg"
)

func NewCmd(cfg *clicfg.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "customer-metrics-integration",
		Short: "View customer metrics integration endpoints",
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

	cmd.AddCommand(NewGetTenantEndpointUrlCmd(cfg))

	return cmd
}

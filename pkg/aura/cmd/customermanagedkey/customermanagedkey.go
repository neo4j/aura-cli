package customermanagedkey

import (
	"github.com/neo4j/cli/pkg/clicfg"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:     "customer-managed-key",
	Short:   "Relates to Customer Managed Keys",
	Aliases: []string{"cmk"},
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if err := clicfg.Viper.BindPFlag("aura.base-url", cmd.Flags().Lookup("base-url")); err != nil {
			return err
		}
		if err := clicfg.Viper.BindPFlag("aura.auth-url", cmd.Flags().Lookup("auth-url")); err != nil {
			return err
		}
		if err := clicfg.Viper.BindPFlag("aura.output", cmd.Flags().Lookup("output")); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	Cmd.PersistentFlags().String("auth-url", "", "")
	Cmd.PersistentFlags().String("base-url", "", "")
	Cmd.PersistentFlags().String("output", "", "")

	Cmd.AddCommand(CreateCmd)
	Cmd.AddCommand(DeleteCmd)
	Cmd.AddCommand(GetCmd)
	Cmd.AddCommand(ListCmd)
}

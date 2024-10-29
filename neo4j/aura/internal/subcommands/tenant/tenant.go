package tenant

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/neo4j/cli/common/clicfg"
)

func NewCmd(cfg *clicfg.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tenant",
		Short: "Relates to an Aura Tenant",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			cfg.Aura.BindBaseUrl(cmd.Flags().Lookup("base-url"))

			cfg.Aura.BindAuthUrl(cmd.Flags().Lookup("auth-url"))

			outputValue := cmd.Flags().Lookup("output").Value.String()
			if outputValue != "" {
				validOutputValue := false
				for _, v := range clicfg.ValidOutputValues {
					if v == outputValue {
						validOutputValue = true
						break
					}
				}
				if !validOutputValue {
					return fmt.Errorf("invalid output value specified: %s", outputValue)
				}
			}

			cfg.Aura.BindOutput(cmd.Flags().Lookup("output"))

			return nil
		},
	}

	cmd.PersistentFlags().String("auth-url", "", "")
	cmd.PersistentFlags().String("base-url", "", "")
	cmd.PersistentFlags().String("output", "", fmt.Sprintf("Format to print console output in, from a choice of [%s]", strings.Join(clicfg.ValidOutputValues[:], ", ")))

	cmd.AddCommand(NewGetCmd(cfg))
	cmd.AddCommand(NewListCmd(cfg))

	return cmd
}

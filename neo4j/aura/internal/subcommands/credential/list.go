package credential

import (
	"github.com/neo4j/cli/common/clicfg"
	"github.com/spf13/cobra"
)

func NewListCmd(cfg *clicfg.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "list credentials",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cfg.Credentials.Aura.Print(cmd.OutOrStdout()); err != nil {
				return err
			}

			return nil
		},
	}
}

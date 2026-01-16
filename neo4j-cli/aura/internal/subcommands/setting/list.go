package setting

import (
	"github.com/neo4j/cli/common/clicfg"
	"github.com/spf13/cobra"
)

func NewListCmd(cfg *clicfg.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "list settings",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cfg.Settings.Aura.Print(cmd.OutOrStdout()); err != nil {
				return err
			}

			return nil
		},
	}
}

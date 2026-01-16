package setting

import (
	"github.com/neo4j/cli/common/clicfg"
	"github.com/spf13/cobra"
)

func NewCmd(cfg *clicfg.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setting",
		Short: "Manage and view setting values",
	}

	cmd.AddCommand(NewAddCmd(cfg))
	cmd.AddCommand(NewUseCmd(cfg))
	cmd.AddCommand(NewListCmd(cfg))
	cmd.AddCommand(NewRemoveCmd(cfg))

	return cmd
}

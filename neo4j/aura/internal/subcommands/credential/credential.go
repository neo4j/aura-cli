package credential

import (
	"github.com/neo4j/cli/common/clicfg"
	"github.com/spf13/cobra"
)

func NewCmd(cfg *clicfg.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "credential",
		Short: "Manage and view credential values",
	}

	cmd.AddCommand(NewAddCmd(cfg))
	cmd.AddCommand(NewRemoveCmd(cfg))
	cmd.AddCommand(NewUseCmd(cfg))
	cmd.AddCommand(NewListCmd(cfg))

	return cmd
}

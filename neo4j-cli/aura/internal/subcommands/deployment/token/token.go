package token

import (
	"github.com/neo4j/cli/common/clicfg"
	"github.com/spf13/cobra"
)

func NewCmd(cfg *clicfg.Config) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "token",
		Short: "Relates to deployment tokens",
	}

	cmd.AddCommand(NewCreateCmd(cfg))
	cmd.AddCommand(NewUpdateCmd(cfg))
	cmd.AddCommand(NewDeleteCmd(cfg))

	return cmd
}

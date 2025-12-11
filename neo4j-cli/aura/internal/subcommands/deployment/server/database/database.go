package serverdatabase

import (
	"github.com/neo4j/cli/common/clicfg"
	"github.com/spf13/cobra"
)

func NewCmd(cfg *clicfg.Config) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "database",
		Short: "Relates to deployment server databases",
	}

	cmd.AddCommand(NewListCmd(cfg))

	return cmd
}

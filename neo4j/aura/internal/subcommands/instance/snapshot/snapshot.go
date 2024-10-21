package snapshot

import (
	"github.com/neo4j/cli/common/clicfg"
	"github.com/spf13/cobra"
)

func NewCmd(cfg *clicfg.Config) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "snapshot",
		Short: "Relates to an instance snapshots",
	}

	cmd.AddCommand(NewListCmd(cfg))
	cmd.AddCommand(NewCreateCmd(cfg))
	cmd.AddCommand(NewGetCmd(cfg))

	return cmd
}

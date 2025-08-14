package job

import (
	"github.com/neo4j/cli/common/clicfg"
	"github.com/spf13/cobra"
)

func NewCmd(cfg *clicfg.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use: "job",
	}

	cmd.AddCommand(NewSpawnCmd(cfg))
	cmd.AddCommand(NewGetCmd(cfg))
	return cmd
}

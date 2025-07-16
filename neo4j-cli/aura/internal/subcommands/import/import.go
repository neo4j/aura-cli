package _import

import (
	"github.com/neo4j/cli/common/clicfg"
	"github.com/spf13/cobra"
)

func NewCmd(cfg *clicfg.Config) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "import",
		Short: "Allows you to import your data into Aura instances and manage your import jobs",
	}
	return cmd
}

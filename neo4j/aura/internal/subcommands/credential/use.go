package credential

import (
	"github.com/neo4j/cli/common/clicfg"
	"github.com/spf13/cobra"
)

func NewUseCmd(cfg *clicfg.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "use <name>",
		Short: "Sets the default credential to be used",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cfg.Aura.SetDefaultCredential(args[0])
		},
	}
}

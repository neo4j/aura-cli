package config

import (
	"github.com/neo4j/cli/common/clicfg"
	"github.com/spf13/cobra"
)

func NewGetCmd(cfg *clicfg.Config) *cobra.Command {
	return &cobra.Command{
		Use:       "get <key>",
		Short:     "Displays the specified configuration value",
		ValidArgs: cfg.Aura.ValidConfigKeys[:],
		Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		RunE: func(cmd *cobra.Command, args []string) error {
			value := cfg.Aura.Get(args[0])

			cmd.Println(value)

			return nil
		},
	}
}

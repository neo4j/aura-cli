package config

import (
	"fmt"

	"github.com/neo4j/cli/common/clicfg"
	"github.com/spf13/cobra"
)

func NewSetCmd(cfg *clicfg.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "set <key> <value>",
		Short: "Sets the specified configuration value to the provided value",
		Args: func(cmd *cobra.Command, args []string) error {
			if err := cobra.ExactArgs(2)(cmd, args); err != nil {
				return err
			}

			if cfg.Aura.IsValidConfigKey(args[0]) {
				return nil
			}

			return fmt.Errorf("invalid config key specified: %s", args[0])
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg.Aura.Set(args[0], args[1])

			return nil
		},
	}
}

package config

import (
	"errors"
	"fmt"

	"github.com/neo4j/cli/common/clictx"
	"github.com/spf13/cobra"
)

func NewSetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "set",
		Short: "Sets the specified configuration value to the provided value",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, ok := clictx.Config(cmd.Context())

			if !ok {
				return errors.New("error fetching configuration values")
			}

			config.Set(fmt.Sprintf("aura.%s", args[0]), args[1])

			config.Write()

			return nil
		},
	}
}

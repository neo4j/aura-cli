package config

import (
	"errors"
	"fmt"

	"github.com/neo4j/cli/pkg/clictx"
	"github.com/spf13/cobra"
)

func NewGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:       "get",
		Short:     "Displays the specified configuration value",
		ValidArgs: []string{"auth-url", "base-url"},
		Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, ok := clictx.Config(cmd.Context())

			if !ok {
				return errors.New("error fetching cli configuration")
			}

			value, err := config.Get(fmt.Sprintf("aura.%s", args[0]))

			if err != nil {
				return err
			}

			cmd.Println(value)

			return nil
		},
	}
}

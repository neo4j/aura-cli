package config

import (
	"errors"

	"github.com/neo4j/cli/common/clictx"
	"github.com/spf13/cobra"
)

func NewListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "Lists the current configuration of the Aura CLI subcommand",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			config, ok := clictx.Config(cmd.Context())

			if !ok {
				return errors.New("error fetching configuration values")
			}

			if err := config.Aura.Print(cmd); err != nil {
				return err
			}

			return nil
		},
	}
}

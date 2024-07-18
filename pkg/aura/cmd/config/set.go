package config

import (
	"errors"
	"fmt"

	"github.com/neo4j/cli/pkg/clictx"
	"github.com/spf13/cobra"
)

var SetCmd = &cobra.Command{
	Use:   "set",
	Short: "Returns details about a specific Aura Instance",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.ExactArgs(2),
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

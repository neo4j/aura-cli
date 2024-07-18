package config

import (
	"errors"
	"fmt"

	"github.com/neo4j/cli/pkg/clictx"
	"github.com/spf13/cobra"
)

var GetCmd = &cobra.Command{
	Use:   "get",
	Short: "Returns details about a specific Aura Instance",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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

package credential

import (
	"errors"

	"github.com/neo4j/cli/pkg/clictx"
	"github.com/spf13/cobra"
)

func NewUseCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "use",
		Short: "Sets the default credential to be used",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, ok := clictx.Config(cmd.Context())

			if !ok {
				return errors.New("error fetching configuration values")
			}

			err := config.Aura.SetDefaultCredential(args[0])
			if err != nil {
				return err
			}

			err = config.Write()
			if err != nil {
				return err
			}

			return nil
		},
	}
}

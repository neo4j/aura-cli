package credential

import (
	"errors"

	"github.com/neo4j/cli/pkg/clictx"
	"github.com/spf13/cobra"
)

func NewAddCmd() *cobra.Command {
	var (
		name         string
		clientId     string
		clientSecret string
	)

	const (
		nameFlag         = "name"
		clientIdFlag     = "client-id"
		clientSecretFlag = "client-secret"
	)

	cmd := &cobra.Command{
		Use:   "add",
		Short: "Adds a credential",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			config, ok := clictx.Config(cmd.Context())

			if !ok {
				return errors.New("error fetching configuration values")
			}

			err := config.Aura.AddCredential(name, clientId, clientSecret)
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

	cmd.Flags().StringVar(&name, nameFlag, "", "Name")
	cmd.MarkFlagRequired(nameFlag)

	cmd.Flags().StringVar(&clientId, clientIdFlag, "", "Client ID")
	cmd.MarkFlagRequired(clientIdFlag)

	cmd.Flags().StringVar(&clientSecret, clientSecretFlag, "", "Client secret")
	cmd.MarkFlagRequired(clientSecretFlag)

	return cmd
}

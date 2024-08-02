package credential

import (
	"errors"

	"github.com/neo4j/cli/common/clictx"
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

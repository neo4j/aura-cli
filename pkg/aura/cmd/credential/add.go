package credential

import (
	"errors"

	"github.com/neo4j/cli/pkg/clictx"
	"github.com/spf13/cobra"
)

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

var AddCmd = &cobra.Command{
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

func init() {
	AddCmd.Flags().StringVar(&name, nameFlag, "", "Name")
	AddCmd.MarkFlagRequired(nameFlag)

	AddCmd.Flags().StringVar(&clientId, clientIdFlag, "", "Client ID")
	AddCmd.MarkFlagRequired(clientIdFlag)

	AddCmd.Flags().StringVar(&clientSecret, clientSecretFlag, "", "Client secret")
	AddCmd.MarkFlagRequired(clientSecretFlag)
}

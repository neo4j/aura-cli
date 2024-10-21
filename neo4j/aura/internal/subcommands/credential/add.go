package credential

import (
	"github.com/neo4j/cli/common/clicfg"
	"github.com/spf13/cobra"
)

func NewAddCmd(cfg *clicfg.Config) *cobra.Command {
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
			// return cfg.Aura.AddCredential(name, clientId, clientSecret)
			return cfg.Credentials.Aura.Add(name, clientId, clientSecret)
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

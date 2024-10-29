package authprovider

import (
	"fmt"
	"net/http"

	"github.com/neo4j/cli/common/clicfg"
	"github.com/neo4j/cli/neo4j/aura/internal/api"
	"github.com/neo4j/cli/neo4j/aura/internal/output"
	"github.com/spf13/cobra"
)

func NewListCmd(cfg *clicfg.Config) *cobra.Command {
	var (
		instanceId string
		dataApiId  string
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "Returns a list of authentication providers of a specific GraphQL Data API",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			path := fmt.Sprintf("/instances/%s/data-apis/graphql/%s/auth-providers", instanceId, dataApiId)

			resBody, statusCode, err := api.MakeRequest(cfg, path, &api.RequestConfig{Method: http.MethodGet})
			if err != nil {
				return err
			}

			if statusCode == http.StatusOK {
				output.PrintBody(cmd, cfg, resBody, []string{"id", "name", "type", "enabled", "url"})
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&instanceId, "instance-id", "", "The ID of the instance the GraphQL Data API is connected to")
	cmd.MarkFlagRequired("instance-id")

	cmd.Flags().StringVar(&dataApiId, "data-api-id", "", "The ID of the GraphQL Data API to list the authentication providers of")
	cmd.MarkFlagRequired("data-api-id")

	return cmd
}

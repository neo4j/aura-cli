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
		Short: "Returns a list of auth providers of a GraphQL Data API",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			path := fmt.Sprintf("/instances/%s/data-apis/graphql/%s/auth-provider", instanceId, dataApiId)
			resBody, statusCode, err := api.MakeRequest(cfg, path, &api.RequestConfig{Method: http.MethodGet})
			if err != nil {
				return err
			}

			if statusCode == http.StatusOK {
				err = output.PrintBody(cmd, cfg, resBody, []string{"id", "name", "type", "enabled", "url"})
				if err != nil {
					return err
				}
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&instanceId, "instance-id", "", "The ID of the instance the Data API is connected to")
	cmd.MarkFlagRequired("instance-id")

	cmd.Flags().StringVar(&dataApiId, "data-api-id", "", "The ID of the Data API to list the auth providers of")
	cmd.MarkFlagRequired("instance-id")

	return cmd
}

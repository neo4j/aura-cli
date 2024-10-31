package authprovider

import (
	"fmt"
	"net/http"

	"github.com/neo4j/cli/common/clicfg"
	"github.com/neo4j/cli/neo4j-cli/aura/internal/api"
	"github.com/neo4j/cli/neo4j-cli/aura/internal/output"
	"github.com/spf13/cobra"
)

func NewDeleteCmd(cfg *clicfg.Config) *cobra.Command {
	var (
		instanceId string
		dataApiId  string
	)

	cmd := &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a GraphQL Data API authentication provider",
		Long:  "Deletes a GraphQL Data API authentication provider. This action can not be undone.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			path := fmt.Sprintf("/instances/%s/data-apis/graphql/%s/auth-providers/%s", instanceId, dataApiId, args[0])

			resBody, statusCode, err := api.MakeRequest(cfg, path, &api.RequestConfig{
				Method: http.MethodDelete,
			})
			if err != nil {
				return err
			}

			// NOTE: delete should not return OK (200), it always returns 202, checking both just in case
			if statusCode == http.StatusAccepted || statusCode == http.StatusOK {
				//
				// TODO: check what is returned. Also update in test if need be
				//
				output.PrintBody(cmd, cfg, resBody, []string{"id", "name", "type", "enabled", "url"})
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&instanceId, "instance-id", "", "The ID of the instance to delete the Data API for")
	cmd.MarkFlagRequired("instance-id")

	cmd.Flags().StringVar(&dataApiId, "data-api-id", "", "The ID of the GraphQL Data API to delete the Authentication provider for")
	cmd.MarkFlagRequired("data-api-id")

	return cmd
}

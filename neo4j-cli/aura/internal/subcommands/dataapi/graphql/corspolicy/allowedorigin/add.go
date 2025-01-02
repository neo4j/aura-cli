package allowedorigin

import (
	"fmt"
	"net/http"

	"github.com/neo4j/cli/common/clicfg"
	"github.com/neo4j/cli/neo4j-cli/aura/internal/api"
	"github.com/neo4j/cli/neo4j-cli/aura/internal/output"
	"github.com/spf13/cobra"
)

func NewAddCmd(cfg *clicfg.Config) *cobra.Command {
	const (
		instanceIdFlag = "instance-id"
		dataApiIdFlag  = "data-api-id"
		awaitFlag      = "await"
	)

	var (
		instanceId string
		dataApiId  string
		await      bool
	)

	cmd := &cobra.Command{
		Use:   "add <origin>",
		Short: "Adds a new allowed origin to the CORS policy",
		Long: `This command adds a new allowed origin to the Cross-Origin Resource Sharing (CORS) policy of a GraphQL Data API.

Updating the CORS policy of a GraphQL Data API is an asynchronous operation. Use the --await flag to wait for the GraphQL Data API to be ready. Once the status transitions from "updating" to "ready" you may begin to use your GraphQL Data API.

Adding a new allowed origin to the CORS policy of a GraphQL Data API allows browsers to make requests to the GraphQL Data API from a web app that is served from the specified origin.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			newOrigin := args[0]

			existingOrigins, err := getGetExistingOrigins(cfg, dataApiId, instanceId)
			if err != nil {
				return err
			}

			for _, origin := range existingOrigins {
				if origin == newOrigin {
					cmd.Println("Origin already exists in allowed origins:", newOrigin)
					return nil
				}
			}

			newOrigins := append(existingOrigins, newOrigin)

			fmt.Println("Existing allowed origins:", existingOrigins)
			fmt.Println("Adding new origin:", newOrigin)
			fmt.Println("New allowed origins:", newOrigins)

			cmd.SilenceUsage = true
			body := map[string]any{
				"security": map[string]any{
					"cors_policy": map[string]any{
						"allowed_origins": newOrigins,
					},
				},
			}
			path := fmt.Sprintf("/instances/%s/data-apis/graphql/%s", instanceId, dataApiId)
			resBody, statusCode, err := api.MakeRequest(cfg, path, &api.RequestConfig{
				PostBody: body,
				Method:   http.MethodPatch,
			})
			if err != nil {
				return err
			}

			// NOTE: Update should not return OK (200), it always returns 202, checking both just in case
			if statusCode == http.StatusAccepted || statusCode == http.StatusOK {
				output.PrintBody(cmd, cfg, resBody, []string{"id", "name", "status", "url"})
				if await {
					cmd.Println("Waiting for GraphQL Data API to be ready...")
					pollResponse, err := api.PollGraphQLDataApi(cfg, instanceId, dataApiId, api.GraphQLDataApiStatusCreating)
					if err != nil {
						return err
					}

					cmd.Println("GraphQL Data API Status:", pollResponse.Data.Status)
				}
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&instanceId, instanceIdFlag, "", "(required) The ID of the instance to create the GraphQL Data API for")
	cmd.MarkFlagRequired(instanceIdFlag)

	cmd.Flags().StringVar(&dataApiId, dataApiIdFlag, "", "(required) The ID of the GraphQL Data API to create the authentication provider for")
	cmd.MarkFlagRequired(dataApiIdFlag)

	cmd.Flags().BoolVar(&await, awaitFlag, false, "Waits until updated GraphQL Data API is ready.")

	return cmd
}

package authprovider

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/neo4j/cli/common/clicfg"
	"github.com/neo4j/cli/neo4j/aura/internal/api"
	"github.com/neo4j/cli/neo4j/aura/internal/output"
	"github.com/spf13/cobra"
)

func NewUpdateCmd(cfg *clicfg.Config) *cobra.Command {
	const (
		instanceIdFlag = "instance-id"
		dataApiIdFlag  = "data-api-id"
		nameFlag       = "name"
		enabledFlag    = "enabled"
		urlFlag        = "url"
		awaitFlag      = "await"
	)

	var (
		instanceId string
		dataApiId  string
		name       string
		enabled    string
		url        string
		await      bool
	)

	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Edit a GraphQL Data API authentication provider",
		Long: `This endpoint edits a specific GraphQL Data API authentication provider.
		
Updating a GraphQL Data API authentication provider is an asynchronous operation. Use the --await flag to wait for the GraphQL Data API to be ready again. Once the status transitions from "updating" to "ready" you may continue to use your GraphQL Data API.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			body := map[string]any{}

			if name != "" {
				body["name"] = name
			}

			if enabled != "" {
				isEnabled, err := strconv.ParseBool(enabled)
				if err != nil {
					return fmt.Errorf("invalid value for boolean enabled, err: %s", err.Error())
				}
				body["enabled"] = isEnabled
			}

			if url != "" {
				body["url"] = url
			}

			cmd.SilenceUsage = true
			path := fmt.Sprintf("/instances/%s/data-apis/graphql/%s/auth-providers/%s", instanceId, dataApiId, args[0])

			resBody, statusCode, err := api.MakeRequest(cfg, path, &api.RequestConfig{
				Method:   http.MethodPatch,
				PostBody: body,
			})
			if err != nil {
				return err
			}

			// NOTE: GraphQL Data API update should not return OK (200), it always returns 202, checking both just in case
			if statusCode == http.StatusAccepted || statusCode == http.StatusOK {
				output.PrintBody(cmd, cfg, resBody, []string{"id", "name", "type", "enabled", "url"})

				if await {
					cmd.Println("Waiting for GraphQL Data API to be updated...")
					pollResponse, err := api.PollGraphQLDataApi(cfg, instanceId, args[0], api.GraphQLDataApiStatusUpdating)
					if err != nil {
						return err
					}

					cmd.Println("GraphQL Data API Status:", pollResponse.Data.Status)
				}
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&instanceId, instanceIdFlag, "", "The ID of the instance the GraphQL Data API is connected to")
	cmd.MarkFlagRequired(instanceIdFlag)

	cmd.Flags().StringVar(&dataApiId, dataApiIdFlag, "", "The ID of the GraphQL Data API to update the authentication providers for")
	cmd.MarkFlagRequired(dataApiIdFlag)

	cmd.Flags().StringVar(&name, nameFlag, "", "The name of the authentication provider")

	cmd.Flags().StringVar(&enabled, enabledFlag, "", "Wether or not the authentication provider is enabled")

	cmd.Flags().StringVar(&url, urlFlag, "", "The url for the JWKS endpoint, NOTE: only applicable for authentication provider type 'jwks'")

	cmd.Flags().BoolVar(&await, awaitFlag, false, "Waits until updated GraphQL Data API is ready again.")

	return cmd
}

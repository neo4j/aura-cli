package graphql

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/neo4j/cli/common/clicfg"
	"github.com/neo4j/cli/neo4j/aura/internal/api"
	"github.com/neo4j/cli/neo4j/aura/internal/output"
	"github.com/spf13/cobra"
)

func NewCreateCmd(cfg *clicfg.Config) *cobra.Command {
	const (
		instanceIdFlag             = "instance-id"
		nameFlag                   = "name"
		instanceUsernameFlag       = "instance-username"
		instancePasswordFlag       = "instance-password"
		typeDefsFlag               = "type-definitions"
		featureSubGraphEnabledFlag = "feature-subgraph-enabled"
		// authProviderFlag             = "auth-provider"
		awaitFlag = "await"

		featureSubGraphEnabledDefault = false
	)

	var (
		instanceId             string
		name                   string
		instanceUsername       string
		instancePassword       string
		typeDefs               string
		featureSubGraphEnabled bool = featureSubGraphEnabledDefault
		// authProvider             []string
		await bool
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a new GraphQL Data API",
		Long: `This endpoint starts the creation process of an Aura GraphQL Data API.

Creating a GraphQL Data API is an asynchronous operation. You can poll the current status of this operation by periodically getting the GraphQL Data API details for the GraphQL Data API ID using the GET /data-apis/graphql/{data-apiId} endpoint. Once the status transitions from "creating" to "ready" you may begin to use your GraphQL Data API.

This endpoint returns your GraphQL Data API ID, API key, and connection URL in the response body for you to use once the GraphQL Data API is running. It is important to store the API key as it is not currently possible to get this or update it.

If you lose your API key, you will need to delete the GraphQL Data API and re-create it. This will not result in any loss of data.`,
		PreRunE: func(cmd *cobra.Command, args []string) error {

			// if auth-provider individual flags mark xxx as required
			// if auth-provider as "slice flag" then do not mark as required

			// typeValue, _ := cmd.Flags().GetString("type")
			// if typeValue != "free-db" {
			// 	cmd.MarkFlagRequired(memoryFlag)
			// 	cmd.MarkFlagRequired(regionFlag)
			// }

			// if cfg.Aura.DefaultTenant() == "" {
			// 	cmd.MarkFlagRequired(tenantIdFlag)
			// }

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			body := map[string]any{
				"name": name,
				"aura_instance": map[string]string{
					"username": instanceUsername,
					"password": instancePassword,
				},
				"features": map[string]bool{
					"subgraph": featureSubGraphEnabled,
				},
			}

			if !IsBase64(typeDefs) {
				return errors.New("Type definitions are not valid base64")
			}

			cmd.SilenceUsage = true
			path := fmt.Sprintf("/instances/%s/data-apis/graphql", instanceId)
			resBody, statusCode, err := api.MakeRequest(cfg, path, &api.RequestConfig{
				PostBody: body,
				Method:   http.MethodPost,
			})
			if err != nil {
				return err
			}

			// NOTE: GraphQL Data API create should not return OK (200), it always returns 202, checking both just in case
			if statusCode == http.StatusAccepted || statusCode == http.StatusOK {

				// TODO: print a "safe this key message" if there is an api-key in there

				err = output.PrintBody(cmd, cfg, resBody, []string{"id", "name", "status", "url", "type_definitions", "auth_provider"})
				if err != nil {
					return err
				}

				if await {
					cmd.Println("Waiting for GraphQL Data API to be ready...")
					var response api.CreateGraphQLDataApiResponse
					if err := json.Unmarshal(resBody, &response); err != nil {
						return err
					}

					pollResponse, err := api.PollGraphQLDataApi(cfg, instanceId, response.Data.Id)
					if err != nil {
						return err
					}

					cmd.Println("GraphQL Data API Status:", pollResponse.Data.Status)
				}
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&instanceId, instanceIdFlag, "", "The ID of the instance to list the GraphQL Data APIs of")
	cmd.MarkFlagRequired(instanceIdFlag)

	cmd.Flags().StringVar(&name, nameFlag, "", "The name of the Data API")
	cmd.MarkFlagRequired(nameFlag)

	cmd.Flags().StringVar(&instanceUsername, instanceUsernameFlag, "", "The username of the instance this GraphQL Data API will be connected to")
	cmd.MarkFlagRequired(instanceUsernameFlag)

	cmd.Flags().StringVar(&instancePassword, instancePasswordFlag, "", "The password of the instance this GraphQL Data API will be connected to")
	cmd.MarkFlagRequired(instancePasswordFlag)

	cmd.Flags().StringVar(&typeDefs, typeDefsFlag, "", "The GraphQL type definitions, NOTE: must be base64 encoded")
	cmd.MarkFlagRequired(typeDefsFlag)

	featureSubGraphHelpMsg := fmt.Sprintf("Wether or not GraphQL subgraph is enabled, default is %t", featureSubGraphEnabledDefault)
	cmd.Flags().BoolVar(&featureSubGraphEnabled, featureSubGraphEnabledFlag, featureSubGraphEnabledDefault, featureSubGraphHelpMsg)

	cmd.Flags().BoolVar(&await, awaitFlag, false, "Waits until created GraphQL Data API is ready.")

	return cmd
}

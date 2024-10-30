package authprovider

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/neo4j/cli/common/clicfg"
	"github.com/neo4j/cli/neo4j/aura/internal/api"
	"github.com/neo4j/cli/neo4j/aura/internal/output"
	"github.com/spf13/cobra"
)

func NewCreateCmd(cfg *clicfg.Config) *cobra.Command {
	const (
		instanceIdFlag = "instance-id"
		dataApiIdFlag  = "data-api-id"
		typeFlag       = "type"
		nameFlag       = "name"
		enabledFlag    = "enabled"
		urlFlag        = "url"
		awaitFlag      = "await"
	)

	var (
		instanceId string
		dataApiId  string
		_type      string
		name       string
		enabled    string
		url        string
		await      bool
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a new GraphQL Data API authentication provider",
		Long: `This command creates a new GraphQL Data API authentication provider.

Creating a GraphQL Data API authentication provider is an asynchronous operation. Use the --await flag to wait for the GraphQL Data API to be ready. Once the status transitions from "creating" to "ready" you may begin to use your GraphQL Data API.

If you create a 'api-key' Authentication provider, a API key will be created. It is important to store the API key as it is not currently possible to get it or update it.

If you lose your API key, you will need to create a new Authentication provider. This will not result in any loss of data.`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			typeValue, _ := cmd.Flags().GetString(typeFlag)
			if typeValue == api.GraphQLDataApiAuthProviderTypeJwks {
				cmd.MarkFlagRequired(urlFlag)
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if _type != api.GraphQLDataApiAuthProviderTypeJwks && _type != api.GraphQLDataApiAuthProviderTypeApiKey {
				msg := strings.ToLower(fmt.Sprintf("invalid authentication provider type, got '%s', expected '%s' or '%s'", _type, api.GraphQLDataApiAuthProviderTypeJwks, api.GraphQLDataApiAuthProviderTypeApiKey))
				return errors.New(msg)
			}

			body := map[string]any{
				"type": _type,
				"name": name,
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
			path := fmt.Sprintf("/instances/%s/data-apis/graphql/%s/auth-providers", instanceId, dataApiId)
			resBody, statusCode, err := api.MakeRequest(cfg, path, &api.RequestConfig{
				PostBody: body,
				Method:   http.MethodPost,
			})
			if err != nil {
				return err
			}

			// NOTE: Auth provider create should not return OK (200), it always returns 202, checking both just in case
			if statusCode == http.StatusAccepted || statusCode == http.StatusOK {

				if _type == api.GraphQLDataApiAuthProviderTypeApiKey {
					cmd.Println("###############################")
					cmd.Println("# It is important to store the created API key! If you lose your API key, you will need to create a new Authentication provider. This will not result in any loss of data.")
					cmd.Println("###############################")
				}

				output.PrintBody(cmd, cfg, resBody, []string{"id", "name", "type", "enabled", "key", "url"})

				if await {
					cmd.Println("Waiting for GraphQL Data API to be ready...")
					//
					// TODO: is it really api.GraphQLDataApiStatusCreating? or are there new statuses?
					//
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

	cmd.Flags().StringVar(&instanceId, instanceIdFlag, "", "The ID of the instance to create the GraphQL Data API for")
	cmd.MarkFlagRequired(instanceIdFlag)

	cmd.Flags().StringVar(&dataApiId, dataApiIdFlag, "", "The ID of the GraphQL Data API to update the authentication providers for")
	cmd.MarkFlagRequired(dataApiIdFlag)

	cmd.Flags().StringVar(&_type, typeFlag, "", "The type of the authentication provider, one of 'api-key' or 'jwks'")
	cmd.MarkFlagRequired(typeFlag)

	cmd.Flags().StringVar(&name, nameFlag, "", "The name of the authentication provider")
	cmd.MarkFlagRequired(nameFlag)

	cmd.Flags().StringVar(&enabled, enabledFlag, "", "Wether or not the authentication provider is enabled")

	cmd.Flags().StringVar(&url, urlFlag, "", "The url for the JWKS endpoint, NOTE: only applicable for authentication provider type 'jwks'")

	cmd.Flags().BoolVar(&await, awaitFlag, false, "Waits until created GraphQL Data API is ready.")

	return cmd
}

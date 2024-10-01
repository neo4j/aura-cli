package tenant

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/spf13/cobra"

	"github.com/neo4j/cli/common/clicfg"
	"github.com/neo4j/cli/neo4j/aura/internal/api"
	"github.com/neo4j/cli/neo4j/aura/internal/output"
)

func NewGetCmd(cfg *clicfg.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Returns tenant details",
		Long:  "This subcommand returns details about a specific Aura Tenant.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			tenantId := args[0]
			path := fmt.Sprintf("/tenants/%s", tenantId)

			cmd.SilenceUsage = true
			resBody, statusCode, err := api.MakeRequest(cfg, path, &api.RequestConfig{
				Method: http.MethodGet,
			})
			if err != nil {
				return err
			}

			if statusCode == http.StatusOK {
				responseData, err := api.ParseBody(resBody)
				if err != nil {
					return err
				}
				fields, values, err := postProcessResponseValues(cfg, tenantId, responseData)
				if err != nil {
					return err
				}
				err = output.PrintBodyMap(cmd, cfg, values, fields)
				if err != nil {
					return err
				}
			}

			return nil
		},
	}
}

func postProcessResponseValues(cfg *clicfg.Config, tenantId string, responseData api.ResponseData) ([]string, api.ResponseData, error) {
	resBody, statusCode, err := api.MakeRequest(cfg, http.MethodGet, fmt.Sprintf("/tenants/%s/metrics-integration", tenantId), nil)
	if err != nil {
		return nil, api.ResponseData{}, err
	}
	if statusCode == http.StatusOK {
		metricsIntegrationResponse, err := api.ParseBody(resBody)
		if err != nil {
			return nil, api.ResponseData{}, err
		}
		metricsIntegration, err := metricsIntegrationResponse.GetOne()
		if err != nil {
			return nil, api.ResponseData{}, err
		}
		fields := []string{"id", "name"}
		switch cmiEndpointUrl := metricsIntegration["endpoint"].(type) {
		case string:
			if len(cmiEndpointUrl) > 0 {
				tenant, err := responseData.GetOne()
				if err != nil {
					return nil, api.ResponseData{}, err
				}
				tenant["metrics_integration_url"] = cmiEndpointUrl
				return append(fields, "metrics_integration_url"), api.NewSingleValueResponseData(tenant), nil
			}
		default:
			return fields, responseData, nil
		}
		return nil, api.ResponseData{}, err
	} else {
		return nil, api.ResponseData{}, errors.New(fmt.Sprintf("Unexpected statusCode %d", statusCode))
	}
}

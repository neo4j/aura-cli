package tenant

import (
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
	metricsIntegrationEndpointUrl, err := getMetricsIntegrationEndpointUrl(cfg, tenantId)
	if err != nil {
		return nil, nil, err
	}
	fields := []string{"id", "name"}
	if len(metricsIntegrationEndpointUrl) > 0 {
		tenant, err := responseData.GetOne()
		if err != nil {
			return nil, nil, err
		}
		tenant["metrics_integration_url"] = metricsIntegrationEndpointUrl
		return append(fields, "metrics_integration_url"), api.NewSingleValueResponseData(tenant), nil
	} else {
		return fields, responseData, nil
	}
}

func getMetricsIntegrationEndpointUrl(cfg *clicfg.Config, tenantId string) (string, error) {
	resBody, statusCode, err := api.MakeRequest(cfg, http.MethodGet, fmt.Sprintf("/tenants/%s/metrics-integration", tenantId), nil)
	// Aura API (in fact Console API returns HTTP 400 when CMI endpoint is not available for the tenant)
	if err != nil && statusCode != http.StatusBadRequest {
		return "", err
	}
	switch {
	case statusCode == http.StatusOK:
		metricsIntegrationResponse, err := api.ParseBody(resBody)
		if err != nil {
			return "", err
		}
		metricsIntegration, err := metricsIntegrationResponse.GetOne()
		if err != nil {
			return "", err
		}
		switch endpointUrl := metricsIntegration["endpoint"].(type) {
		case string:
			if len(endpointUrl) > 0 {
				return endpointUrl, nil
			}
		default:
			return "", nil
		}
		return "", err
	case statusCode == http.StatusBadRequest:
		return "", nil
	default:
		return "", fmt.Errorf("Unexpected statusCode %d", statusCode)
	}
}

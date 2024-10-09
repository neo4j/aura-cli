package project

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
		Short: "Returns project details",
		Long:  "This subcommand returns details about a specific Aura Project.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			projectId := args[0]
			path := fmt.Sprintf("/tenants/%s", projectId)

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
				fields, values, err := postProcessResponseValues(cfg, projectId, responseData)
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

func postProcessResponseValues(cfg *clicfg.Config, projectId string, responseData api.ResponseData) ([]string, api.ResponseData, error) {
	metricsIntegrationEndpointUrl, err := getMetricsIntegrationEndpointUrl(cfg, projectId)
	if err != nil {
		return nil, nil, err
	}
	fields := []string{"id", "name"}
	if len(metricsIntegrationEndpointUrl) > 0 {
		project, err := responseData.GetSingleOrError()
		if err != nil {
			return nil, nil, err
		}
		project["metrics_integration_url"] = metricsIntegrationEndpointUrl
		return append(fields, "metrics_integration_url"), api.NewSingleValueResponseData(project), nil
	} else {
		return fields, responseData, nil
	}
}

func getMetricsIntegrationEndpointUrl(cfg *clicfg.Config, projectId string) (string, error) {
	resBody, statusCode, err := api.MakeRequest(cfg, fmt.Sprintf("/tenants/%s/metrics-integration", projectId), &api.RequestConfig{
		Method: http.MethodGet,
	})
	// Aura API (in fact Console API returns HTTP 400 when CMI endpoint is not available for the project)
	if err != nil && statusCode != http.StatusBadRequest {
		return "", err
	}
	switch {
	case statusCode == http.StatusOK:
		metricsIntegrationResponse, err := api.ParseBody(resBody)
		if err != nil {
			return "", err
		}
		metricsIntegration, err := metricsIntegrationResponse.GetSingleOrError()
		if err != nil {
			return "", err
		}
		if endpointUrl, ok := metricsIntegration["endpoint"].(string); ok {
			if len(endpointUrl) > 0 {
				return endpointUrl, nil
			}
		}
		return "", nil
	case statusCode == http.StatusBadRequest:
		return "", nil
	default:
		return "", fmt.Errorf("unexpected statusCode %d", statusCode)
	}
}

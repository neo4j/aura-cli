package instance

import (
	"fmt"
	"net/http"

	"github.com/neo4j/cli/common/clicfg"
	"github.com/neo4j/cli/neo4j/aura/internal/api"
	"github.com/neo4j/cli/neo4j/aura/internal/output"
	"github.com/spf13/cobra"
)

func NewGetCmd(cfg *clicfg.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Returns instance details",
		Long:  "This endpoint returns details about a specific Aura Instance.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			instanceId := args[0]
			path := fmt.Sprintf("/instances/%s", instanceId)

			cmd.SilenceUsage = true
			resBody, statusCode, err := api.MakeRequest(cfg, path, &api.RequestConfig{
				Method: http.MethodGet,
			})
			if err != nil {
				return err
			}

			if statusCode == http.StatusOK {
				fields, err := getFields(resBody)
				if err != nil {
					return err
				}
				err = output.PrintBody(cmd, cfg, resBody, fields)
				if err != nil {
					return err
				}
			}

			return nil
		},
	}
}

func getFields(resBody []byte) ([]string, error) {
	instances, err := api.ParseBody(resBody)
	if err != nil {
		return nil, err
	}
	if len(instances) != 1 {
		return nil, fmt.Errorf("expected 1 instance, got %d", len(instances))
	}
	fields := []string{"id", "name", "tenant_id", "status", "connection_url", "cloud_provider", "region", "type", "memory", "storage", "customer_managed_key_id"}
	if HasCmiEndpoint(instances[0]) {
		fields = append(fields, "metrics_integration_url")
	}
	return fields, nil
}

func HasCmiEndpoint(instance map[string]any) bool {
	cmiEndpointUrl := instance["metrics_integration_url"]
	switch cmiEndpointUrl := cmiEndpointUrl.(type) {
	case string:
		return len(cmiEndpointUrl) > 0
	}
	return false
}

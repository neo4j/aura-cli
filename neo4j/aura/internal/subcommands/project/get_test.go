package project_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/neo4j/cli/neo4j/aura/internal/test/testutils"
)

func TestGetProjectWithoutIntegrationEndpoint(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	projectId := "6981ace7-efe8-4f5c-b7c5-267b5162ce91"

	getMockHandler := helper.NewRequestHandlerMock(fmt.Sprintf("/v1/tenants/%s", projectId), http.StatusOK, `{
			"data": {
				"id": "6981ace7-efe8-4f5c-b7c5-267b5162ce91",
				"name": "Production",
				"instance_configurations": []
			}
		}`)

	metricsIntegrationMockHandler := helper.NewRequestHandlerMock(fmt.Sprintf("/v1/tenants/%s/metrics-integration", projectId), http.StatusBadRequest, `{
			"errors": [
				{
					"message": "This project has no instances eligible for metrics integration",
					"reason": "tenant-incapable-of-action"
				}
			]
		}`)

	helper.ExecuteCommand(fmt.Sprintf("project get %s", projectId))

	getMockHandler.AssertCalledTimes(1)
	getMockHandler.AssertCalledWithMethod(http.MethodGet)
	metricsIntegrationMockHandler.AssertCalledTimes(1)
	metricsIntegrationMockHandler.AssertCalledWithMethod(http.MethodGet)

	helper.AssertOutJson(`{
		"data": {
			"id": "6981ace7-efe8-4f5c-b7c5-267b5162ce91",
			"instance_configurations": [],
			"name": "Production"
		}
	}
	`)
}
func TestGetProjectWithMetricsIntegrationEndpoint(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	projectId := "6981ace7-efe8-4f5c-b7c5-267b5162ce91"

	getMockHandler := helper.NewRequestHandlerMock(fmt.Sprintf("/v1/tenants/%s", projectId), http.StatusOK, `{
			"data": {
				"id": "6981ace7-efe8-4f5c-b7c5-267b5162ce91",
				"name": "Production",
				"instance_configurations": []
			}
		}`)

	metricsIntegrationMockHandler := helper.NewRequestHandlerMock(fmt.Sprintf("/v1/tenants/%s/metrics-integration", projectId), http.StatusOK, `{
			"data": {
				"endpoint": "https://customer-metrics-api-devnommrr.neo4j-dev.io/api/v1/ca7bc96c-204c-546e-9736-f4a578d53f64/metrics"
			}
		}`)

	helper.ExecuteCommand(fmt.Sprintf("project get %s", projectId))

	getMockHandler.AssertCalledTimes(1)
	getMockHandler.AssertCalledWithMethod(http.MethodGet)
	metricsIntegrationMockHandler.AssertCalledTimes(1)
	metricsIntegrationMockHandler.AssertCalledWithMethod(http.MethodGet)

	helper.AssertOutJson(`{
		"data": {
			"id": "6981ace7-efe8-4f5c-b7c5-267b5162ce91",
			"instance_configurations": [],
			"metrics_integration_url": "https://customer-metrics-api-devnommrr.neo4j-dev.io/api/v1/ca7bc96c-204c-546e-9736-f4a578d53f64/metrics",
			"name": "Production"
		}
	}
	`)
}

func TestGetProjectNotFoundError(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	projectId := "6981ace7-efe8-4f5c-b7c5-267b5162ce91"

	mockHandler := helper.NewRequestHandlerMock(fmt.Sprintf("/v1/tenants/%s", projectId), http.StatusNotFound, `{
		"errors": [
			{
			"message": "The tenant you specified could not be found",
			"reason": "tenant-not-found"
			}
		]
		}`)

	helper.ExecuteCommand(fmt.Sprintf("project get %s", projectId))

	mockHandler.AssertCalledTimes(1)
	mockHandler.AssertCalledWithMethod(http.MethodGet)

	helper.AssertErr("Error: [The tenant you specified could not be found]")
}

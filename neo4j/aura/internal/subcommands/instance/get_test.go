package instance_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/neo4j/cli/neo4j/aura/internal/test/testutils"
)

func TestGetInstance(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	instanceId := "2f49c2b3"

	mockHandler := helper.NewRequestHandlerMock(fmt.Sprintf("/v1/instances/%s", instanceId), http.StatusOK, `{
			"data": {
				"id": "2f49c2b3",
				"name": "Production",
				"status": "running",
				"tenant_id": "YOUR_TENANT_ID",
				"cloud_provider": "gcp",
				"connection_url": "YOUR_CONNECTION_URL",
				"metrics_integration_url": "YOUR_METRICS_INTEGRATION_ENDPOINT",
				"region": "europe-west1",
				"type": "enterprise-db",
				"memory": "8GB",
				"storage": "16GB"
			}
		}`)

	helper.ExecuteCommand(fmt.Sprintf("instance get %s", instanceId))

	mockHandler.AssertCalledTimes(1)
	mockHandler.AssertCalledWithMethod(http.MethodGet)

	helper.AssertOutJson(`{
		"data": {
			"id": "2f49c2b3",
			"name": "Production",
			"status": "running",
			"tenant_id": "YOUR_TENANT_ID",
			"cloud_provider": "gcp",
			"connection_url": "YOUR_CONNECTION_URL",
			"metrics_integration_url": "YOUR_METRICS_INTEGRATION_ENDPOINT",
			"region": "europe-west1",
			"type": "enterprise-db",
			"memory": "8GB",
			"storage": "16GB"
		}
	}
	`)
}

func TestGetInstanceWithTableOutput(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	instanceId := "2f49c2b3"

	mockHandler := helper.NewRequestHandlerMock(fmt.Sprintf("/v1/instances/%s", instanceId), http.StatusOK, `{
			"data": {
				"id": "2f49c2b3",
				"name": "Production",
				"status": "running",
				"tenant_id": "YOUR_TENANT_ID",
				"cloud_provider": "gcp",
				"connection_url": "YOUR_CONNECTION_URL",
				"metrics_integration_url": "YOUR_METRICS_INTEGRATION_ENDPOINT",
				"region": "europe-west1",
				"type": "enterprise-db",
				"memory": "8GB",
				"storage": "16GB"
			}
		}`)

	// TODO: Make a better way to override config
	helper.SetConfigValue("aura.output", "default")

	helper.ExecuteCommand(fmt.Sprintf("instance get %s", instanceId))

	mockHandler.AssertCalledTimes(1)
	mockHandler.AssertCalledWithMethod(http.MethodGet)

	helper.AssertOut(`
┌──────────┬────────────┬────────────────┬─────────┬─────────────────────┬────────────────┬──────────────┬───────────────┬────────┬─────────┬─────────────────────────┬───────────────────────────────────┐
│ ID       │ NAME       │ TENANT_ID      │ STATUS  │ CONNECTION_URL      │ CLOUD_PROVIDER │ REGION       │ TYPE          │ MEMORY │ STORAGE │ CUSTOMER_MANAGED_KEY_ID │ METRICS_INTEGRATION_URL           │
├──────────┼────────────┼────────────────┼─────────┼─────────────────────┼────────────────┼──────────────┼───────────────┼────────┼─────────┼─────────────────────────┼───────────────────────────────────┤
│ 2f49c2b3 │ Production │ YOUR_TENANT_ID │ running │ YOUR_CONNECTION_URL │ gcp            │ europe-west1 │ enterprise-db │ 8GB    │ 16GB    │                         │ YOUR_METRICS_INTEGRATION_ENDPOINT │
└──────────┴────────────┴────────────────┴─────────┴─────────────────────┴────────────────┴──────────────┴───────────────┴────────┴─────────┴─────────────────────────┴───────────────────────────────────┘
`)

}

func TestGetInstanceNotFoundError(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	instanceId := "2f49c2b3"

	mockHandler := helper.NewRequestHandlerMock(fmt.Sprintf("/v1/instances/%s", instanceId), http.StatusNotFound, fmt.Sprintf(`{
		"errors": [
			{
			"message": "DB not found: %s",
			"reason": "db-not-found"
			}
		]
	}`, instanceId))

	helper.ExecuteCommand(fmt.Sprintf("instance get %s", instanceId))

	mockHandler.AssertCalledTimes(1)
	mockHandler.AssertCalledWithMethod(http.MethodGet)

	helper.AssertErr(fmt.Sprintf("Error: [DB not found: %s]", instanceId))
}

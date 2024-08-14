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

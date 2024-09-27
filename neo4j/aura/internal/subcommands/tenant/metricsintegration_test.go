package tenant_test

import (
	"net/http"
	"testing"

	"github.com/neo4j/cli/neo4j/aura/internal/test/testutils"
)

func TestMetricsIntegration(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	mockHandler := helper.NewRequestHandlerMock("/v1/tenants/TENANT_ID/metrics-integration", http.StatusOK, `{
			"data": {
				"endpoint": "MY_ENDPOINT_URL"
			}
		}`)

	helper.ExecuteCommand("tenant metrics-integration --tenant-id TENANT_ID")

	mockHandler.AssertCalledTimes(1)
	mockHandler.AssertCalledWithMethod(http.MethodGet)

	helper.AssertOutJson(`{
			"data": {
				"endpoint": "MY_ENDPOINT_URL"
			}
	}
	`)
}

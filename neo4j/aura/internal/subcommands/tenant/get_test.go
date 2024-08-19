package tenant_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/neo4j/cli/neo4j/aura/internal/test/testutils"
)

func TestGetTenant(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	tenantId := "6981ace7-efe8-4f5c-b7c5-267b5162ce91"

	mockHandler := helper.NewRequestHandlerMock(fmt.Sprintf("/v1/tenants/%s", tenantId), http.StatusOK, `{
			"data": {
				"id": "6981ace7-efe8-4f5c-b7c5-267b5162ce91",
				"name": "Production",
				"instance_configurations": []
			}
		}`)

	helper.ExecuteCommand(fmt.Sprintf("tenant get %s", tenantId))

	mockHandler.AssertCalledTimes(1)
	mockHandler.AssertCalledWithMethod(http.MethodGet)

	helper.AssertOutJson(`{
		"data": {
			"id": "6981ace7-efe8-4f5c-b7c5-267b5162ce91",
			"name": "Production",
			"instance_configurations": []
		}
	}
	`)
}

func TestGetTenantNotFoundError(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	tenantId := "6981ace7-efe8-4f5c-b7c5-267b5162ce91"

	mockHandler := helper.NewRequestHandlerMock(fmt.Sprintf("/v1/tenants/%s", tenantId), http.StatusNotFound, `{
		"errors": [
			{
			"message": "The tenant you specified could not be found",
			"reason": "tenant-not-found"
			}
		]
		}`)

	helper.ExecuteCommand(fmt.Sprintf("tenant get %s", tenantId))

	mockHandler.AssertCalledTimes(1)
	mockHandler.AssertCalledWithMethod(http.MethodGet)

	helper.AssertErr("Error: [The tenant you specified could not be found]")
}

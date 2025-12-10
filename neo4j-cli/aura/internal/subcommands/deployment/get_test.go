package deployment_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/neo4j/cli/neo4j-cli/aura/internal/test/testutils"
)

func TestGetDeployment(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	organizationId := "81e4ae5c-171b-4700-b243-8d1dd34f7321"
	projectId := "ef7faf53-fb7e-4994-8d0f-64ae56e91c42"
	deploymentId := "9a1e6181-7d0b-48a2-bc2b-4250c36b5cc2"

	mockHandler := helper.NewRequestHandlerMock(fmt.Sprintf("/v2beta1/organizations/%s/projects/%s/fleet-manager/deployments/%s", organizationId, projectId, deploymentId), http.StatusOK, `{
		"data": {
			"id": "9a1e6181-7d0b-48a2-bc2b-4250c36b5cc2",
			"name": "Test Deployment",
			"status": "running",
			"created_by": "941d32f6-6abf-42d7-beb8-012341376dc6",
			"dbms": {
				"edition": "enterprise",
				"metric_collection_enabled": true,
				"packaging": "PACKAGING"
			},
			"token": {
				"id": "941d32f6-6abf-42d7-beb8-012341376dc6",
				"claimed_time": "CLAIMED_BY",
				"expiry_time": "EXPIRY_TIME",
				"last_used_time": "LAST_USED",
				"release_time": "RELEASE_TIME",
				"auto_rotated": true,
				"created_by": "941d32f6-6abf-42d7-beb8-012341376dc6",
				"creation_time": "CREATION_TIME"
			}
		}
	}`)

	helper.SetConfigValue("aura.beta-enabled", true)
	helper.SetConfigValue("aura.output", "json")
	helper.ExecuteCommand(fmt.Sprintf("deployment get %s --organization-id=%s --project-id=%s", deploymentId, organizationId, projectId))

	mockHandler.AssertCalledTimes(1)
	mockHandler.AssertCalledWithMethod(http.MethodGet)

	helper.AssertOutJson(`{
		"data": {
			"created_by": "941d32f6-6abf-42d7-beb8-012341376dc6",
			"dbms": {
				"edition": "enterprise",
				"metric_collection_enabled": true,
				"packaging": "PACKAGING"
			},
			"id": "9a1e6181-7d0b-48a2-bc2b-4250c36b5cc2",
			"name": "Test Deployment",
			"status": "running",
			"token": {
				"auto_rotated": true,
				"claimed_time": "CLAIMED_BY",
				"created_by": "941d32f6-6abf-42d7-beb8-012341376dc6",
				"creation_time": "CREATION_TIME",
				"expiry_time": "EXPIRY_TIME",
				"id": "941d32f6-6abf-42d7-beb8-012341376dc6",
				"last_used_time": "LAST_USED",
				"release_time": "RELEASE_TIME"
			}
		}
	}`)
}

func TestGetDeploymentWithTableOutput(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	organizationId := "81e4ae5c-171b-4700-b243-8d1dd34f7321"
	projectId := "ef7faf53-fb7e-4994-8d0f-64ae56e91c42"
	deploymentId := "9a1e6181-7d0b-48a2-bc2b-4250c36b5cc2"

	mockHandler := helper.NewRequestHandlerMock(fmt.Sprintf("/v2beta1/organizations/%s/projects/%s/fleet-manager/deployments/%s", organizationId, projectId, deploymentId), http.StatusOK, `{
		"data": {
			"id": "9a1e6181-7d0b-48a2-bc2b-4250c36b5cc2",
			"name": "Test Deployment",
			"status": "running",
			"created_by": "941d32f6-6abf-42d7-beb8-012341376dc6",
			"dbms": {
				"edition": "enterprise",
				"metric_collection_enabled": true,
				"packaging": "PACKAGING"
			},
			"token": {
				"id": "941d32f6-6abf-42d7-beb8-012341376dc6",
				"claimed_time": "CLAIMED_BY",
				"expiry_time": "EXPIRY_TIME",
				"last_used_time": "LAST_USED",
				"release_time": "RELEASE_TIME",
				"auto_rotated": true,
				"created_by": "941d32f6-6abf-42d7-beb8-012341376dc6",
				"creation_time": "CREATION_TIME"
			}
		}
	}`)

	helper.SetConfigValue("aura.beta-enabled", true)
	helper.SetConfigValue("aura.output", "default")
	helper.ExecuteCommand(fmt.Sprintf("deployment get %s --organization-id=%s --project-id=%s", deploymentId, organizationId, projectId))

	mockHandler.AssertCalledTimes(1)
	mockHandler.AssertCalledWithMethod(http.MethodGet)

	helper.AssertOut(`
┌──────────────────────────────────────┬─────────────────┬──────────────┬────────────────┬───────────────────┬────────────────────┬─────────────────────┐
│ ID                                   │ NAME            │ DBMS:EDITION │ DBMS:PACKAGING │ TOKEN:EXPIRY_TIME │ TOKEN:AUTO_ROTATED │ TOKEN:CREATION_TIME │
├──────────────────────────────────────┼─────────────────┼──────────────┼────────────────┼───────────────────┼────────────────────┼─────────────────────┤
│ 9a1e6181-7d0b-48a2-bc2b-4250c36b5cc2 │ Test Deployment │ enterprise   │ PACKAGING      │ EXPIRY_TIME       │ true               │ CREATION_TIME       │
└──────────────────────────────────────┴─────────────────┴──────────────┴────────────────┴───────────────────┴────────────────────┴─────────────────────┘
	`)
}

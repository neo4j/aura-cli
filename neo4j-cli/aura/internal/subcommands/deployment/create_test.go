package deployment_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/neo4j/cli/neo4j-cli/aura/internal/test/testutils"
)

func TestCreateDeployment(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	organizationId := "81e4ae5c-171b-4700-b243-8d1dd34f7321"
	projectId := "ef7faf53-fb7e-4994-8d0f-64ae56e91c42"
	name := "Test Deployment"

	mockHandler := helper.NewRequestHandlerMock(fmt.Sprintf("/v2beta1/organizations/%s/projects/%s/deployments", organizationId, projectId), http.StatusCreated, `{
		"data": {
			"id": "9a1e6181-7d0b-48a2-bc2b-4250c36b5cc2"
		}
	}`)

	helper.SetConfigValue("aura.beta-enabled", true)
	helper.SetConfigValue("aura.output", "json")
	helper.ExecuteCommand(fmt.Sprintf("deployment create --name \"%s\" --organization-id %s --project-id %s", name, organizationId, projectId))

	mockHandler.AssertCalledTimes(1)
	mockHandler.AssertCalledWithMethod(http.MethodPost)
	mockHandler.AssertCalledWithBody(fmt.Sprintf(`{"connection_url":"","name":"%s"}`, name))

	helper.AssertOutJson(`{
		"data": {
			"id": "9a1e6181-7d0b-48a2-bc2b-4250c36b5cc2"
		}
	}`)
}

func TestCreateDeploymentWithConnectionUrl(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	organizationId := "81e4ae5c-171b-4700-b243-8d1dd34f7321"
	projectId := "ef7faf53-fb7e-4994-8d0f-64ae56e91c42"
	name := "Test Deployment"

	mockHandler := helper.NewRequestHandlerMock(fmt.Sprintf("/v2beta1/organizations/%s/projects/%s/deployments", organizationId, projectId), http.StatusCreated, `{
		"data": {
			"id": "9a1e6181-7d0b-48a2-bc2b-4250c36b5cc2"
		}
	}`)

	helper.SetConfigValue("aura.beta-enabled", true)
	helper.SetConfigValue("aura.output", "json")
	helper.ExecuteCommand(fmt.Sprintf("deployment create --name \"%s\" --connection-url \"http://localhost:7876\" --organization-id %s --project-id %s", name, organizationId, projectId))

	mockHandler.AssertCalledTimes(1)
	mockHandler.AssertCalledWithMethod(http.MethodPost)
	mockHandler.AssertCalledWithBody(fmt.Sprintf(`{"connection_url":"http://localhost:7876","name":"%s"}`, name))

	helper.AssertOutJson(`{
		"data": {
			"id": "9a1e6181-7d0b-48a2-bc2b-4250c36b5cc2"
		}
	}`)
}

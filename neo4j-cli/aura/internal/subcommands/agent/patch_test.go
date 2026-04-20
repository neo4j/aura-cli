// Copyright (c) "Neo4j"
// Neo4j Sweden AB [http://neo4j.com]

package agent_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/neo4j/cli/neo4j-cli/aura/internal/test/testutils"
)

func TestPatchAgentName(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	organizationId := "81e4ae5c-171b-4700-b243-8d1dd34f7321"
	projectId := "ef7faf53-fb7e-4994-8d0f-64ae56e91c42"
	agentId := "f47ac10b-58cc-4372-a567-0e02b2c3d479"

	mockHandler := helper.NewRequestHandlerMock(fmt.Sprintf("/v2beta1/organizations/%s/projects/%s/agents/%s", organizationId, projectId, agentId), http.StatusOK, `{
		"id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
		"name": "Renamed Agent",
		"description": "An agent that queries the database",
		"dbid": "a1b2c3d4",
		"is_private": false,
		"is_mcp_enabled": false,
		"enabled": true
	}`)

	helper.SetConfigValue("aura.beta-enabled", true)
	helper.SetConfigValue("aura.output", "json")
	helper.ExecuteCommand(fmt.Sprintf(
		`agent patch %s --name "Renamed Agent" --organization-id %s --project-id %s`,
		agentId, organizationId, projectId,
	))

	mockHandler.AssertCalledTimes(1)
	mockHandler.AssertCalledWithMethod(http.MethodPatch)
	mockHandler.AssertCalledWithBody(`{"name": "Renamed Agent"}`)

	helper.AssertOutJson(`{
		"id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
		"name": "Renamed Agent",
		"description": "An agent that queries the database",
		"dbid": "a1b2c3d4",
		"is_private": false,
		"is_mcp_enabled": false,
		"enabled": true
	}`)
}

func TestPatchAgentMultipleFields(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	organizationId := "81e4ae5c-171b-4700-b243-8d1dd34f7321"
	projectId := "ef7faf53-fb7e-4994-8d0f-64ae56e91c42"
	agentId := "f47ac10b-58cc-4372-a567-0e02b2c3d479"

	mockHandler := helper.NewRequestHandlerMock(fmt.Sprintf("/v2beta1/organizations/%s/projects/%s/agents/%s", organizationId, projectId, agentId), http.StatusOK, `{
		"id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
		"name": "My Agent",
		"description": "An agent that queries the database",
		"dbid": "a1b2c3d4",
		"is_private": true,
		"is_mcp_enabled": true,
		"enabled": true
	}`)

	helper.SetConfigValue("aura.beta-enabled", true)
	helper.SetConfigValue("aura.output", "json")
	helper.ExecuteCommand(fmt.Sprintf(
		`agent patch %s --is-private --is-mcp-enabled --organization-id %s --project-id %s`,
		agentId, organizationId, projectId,
	))

	mockHandler.AssertCalledTimes(1)
	mockHandler.AssertCalledWithMethod(http.MethodPatch)
	mockHandler.AssertCalledWithBody(`{"is_mcp_enabled": true, "is_private": true}`)

	helper.AssertOutJson(`{
		"id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
		"name": "My Agent",
		"description": "An agent that queries the database",
		"dbid": "a1b2c3d4",
		"is_private": true,
		"is_mcp_enabled": true,
		"enabled": true
	}`)
}

func TestPatchAgentTools(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	organizationId := "81e4ae5c-171b-4700-b243-8d1dd34f7321"
	projectId := "ef7faf53-fb7e-4994-8d0f-64ae56e91c42"
	agentId := "f47ac10b-58cc-4372-a567-0e02b2c3d479"

	mockHandler := helper.NewRequestHandlerMock(fmt.Sprintf("/v2beta1/organizations/%s/projects/%s/agents/%s", organizationId, projectId, agentId), http.StatusOK, `{
		"id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
		"name": "My Agent",
		"description": "An agent that queries the database",
		"dbid": "a1b2c3d4",
		"is_private": false,
		"is_mcp_enabled": false,
		"enabled": true
	}`)

	helper.SetConfigValue("aura.beta-enabled", true)
	helper.SetConfigValue("aura.output", "json")
	helper.ExecuteCommand(fmt.Sprintf(
		`agent patch %s --tools '%s' --organization-id %s --project-id %s`,
		agentId, testTools, organizationId, projectId,
	))

	mockHandler.AssertCalledTimes(1)
	mockHandler.AssertCalledWithMethod(http.MethodPatch)
	mockHandler.AssertCalledWithBody(fmt.Sprintf(`{"tools": %s}`, testTools))

	helper.AssertOutJson(`{
		"id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
		"name": "My Agent",
		"description": "An agent that queries the database",
		"dbid": "a1b2c3d4",
		"is_private": false,
		"is_mcp_enabled": false,
		"enabled": true
	}`)
}

func TestPatchAgentWithOrganizationAndProjectIdFromConfig(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	organizationId := "81e4ae5c-171b-4700-b243-8d1dd34f7321"
	projectId := "ef7faf53-fb7e-4994-8d0f-64ae56e91c42"
	agentId := "f47ac10b-58cc-4372-a567-0e02b2c3d479"

	mockHandler := helper.NewRequestHandlerMock(fmt.Sprintf("/v2beta1/organizations/%s/projects/%s/agents/%s", organizationId, projectId, agentId), http.StatusOK, `{
		"id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
		"name": "Renamed Agent",
		"description": "An agent that queries the database",
		"dbid": "a1b2c3d4",
		"is_private": false,
		"is_mcp_enabled": false,
		"enabled": true
	}`)

	helper.SetConfigValue("aura.beta-enabled", true)
	helper.SetConfigValue("aura.output", "json")
	helper.SetDefaultProjectInConfig(organizationId, projectId)
	helper.ExecuteCommand(fmt.Sprintf(`agent patch %s --name "Renamed Agent"`, agentId))

	mockHandler.AssertCalledTimes(1)
	mockHandler.AssertCalledWithMethod(http.MethodPatch)
	mockHandler.AssertCalledWithBody(`{"name": "Renamed Agent"}`)

	helper.AssertOutJson(`{
		"id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
		"name": "Renamed Agent",
		"description": "An agent that queries the database",
		"dbid": "a1b2c3d4",
		"is_private": false,
		"is_mcp_enabled": false,
		"enabled": true
	}`)
}

func TestPatchAgentWithNoFields(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	organizationId := "81e4ae5c-171b-4700-b243-8d1dd34f7321"
	projectId := "ef7faf53-fb7e-4994-8d0f-64ae56e91c42"
	agentId := "f47ac10b-58cc-4372-a567-0e02b2c3d479"

	helper.SetConfigValue("aura.beta-enabled", true)
	helper.ExecuteCommand(fmt.Sprintf("agent patch %s --organization-id %s --project-id %s", agentId, organizationId, projectId))

	helper.AssertErr("Error: at least one field must be specified")
}

func TestPatchAgentWithInvalidToolsJSON(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	organizationId := "81e4ae5c-171b-4700-b243-8d1dd34f7321"
	projectId := "ef7faf53-fb7e-4994-8d0f-64ae56e91c42"
	agentId := "f47ac10b-58cc-4372-a567-0e02b2c3d479"

	helper.SetConfigValue("aura.beta-enabled", true)
	helper.ExecuteCommand(fmt.Sprintf(
		`agent patch %s --tools "not-valid-json" --organization-id %s --project-id %s`,
		agentId, organizationId, projectId,
	))

	helper.AssertErr("Error: invalid tools JSON: invalid character 'o' in literal null (expecting 'u')")
}

func TestPatchAgentNotFound(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	organizationId := "81e4ae5c-171b-4700-b243-8d1dd34f7321"
	projectId := "ef7faf53-fb7e-4994-8d0f-64ae56e91c42"
	agentId := "non-existent-agent-id"

	mockHandler := helper.NewRequestHandlerMock(fmt.Sprintf("/v2beta1/organizations/%s/projects/%s/agents/%s", organizationId, projectId, agentId), http.StatusNotFound, `{
		"errors": [{"message": "Agent not found"}]
	}`)

	helper.SetConfigValue("aura.beta-enabled", true)
	helper.ExecuteCommand(fmt.Sprintf(
		`agent patch %s --name "New Name" --organization-id %s --project-id %s`,
		agentId, organizationId, projectId,
	))

	mockHandler.AssertCalledTimes(1)
	mockHandler.AssertCalledWithMethod(http.MethodPatch)

	helper.AssertErr("Error: [Agent not found]")
}

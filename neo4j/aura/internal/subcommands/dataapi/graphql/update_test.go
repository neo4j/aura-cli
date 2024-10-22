package graphql_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/neo4j/cli/neo4j/aura/internal/test/testutils"
)

func TestUpdateGraphQLDataApiOneTypeDefs(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	helper.SetConfigValue("aura.beta-enabled", "true")

	instanceId := "2f49c2b3"
	dataApiId := "afdb4e9d"

	helper.ExecuteCommand(fmt.Sprintf("data-api graphql update --output json --instance-id %s --type-definitions bla --type-definitions-file blabla %s", instanceId, dataApiId))

	helper.AssertErr("Error: only one of '--type-definitions' or '--type-definitions-file' flag can be provided")
}

func TestUpdateGraphQLDataApiNewName(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	helper.SetConfigValue("aura.beta-enabled", "true")

	instanceId := "2f49c2b3"
	dataApiId := "afdb4e9d"
	newName := "friendly-name-4"
	mockHandler := helper.NewRequestHandlerMock(fmt.Sprintf("/v1/instances/%s/data-apis/graphql/%s", instanceId, dataApiId), http.StatusAccepted, `{
			"data": {
                "id": "afdb4e9d",
                "name": "friendly-name-4",
                "status": "ready",
                "url": "https://afdb4e9d.28be6e4d8d3e836019.graphql.neo4j.io/graphql"
        	}
		}`)

	helper.ExecuteCommand(fmt.Sprintf("data-api graphql update --output json --instance-id %s --name %s %s", instanceId, newName, dataApiId))

	mockHandler.AssertCalledTimes(1)
	mockHandler.AssertCalledWithMethod(http.MethodPatch)
	mockHandler.AssertCalledWithBody(`{"name":"friendly-name-4"}`)

	helper.AssertOutJson(`{
		"data": {
			"id": "afdb4e9d",
			"name": "friendly-name-4",
			"status": "ready",
			"url": "https://afdb4e9d.28be6e4d8d3e836019.graphql.neo4j.io/graphql"
        }
	}`)
}

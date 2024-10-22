package graphql_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/neo4j/cli/neo4j/aura/internal/test/testutils"
)

func TestUpdateGraphQLDataApi(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	helper.SetConfigValue("aura.beta-enabled", "true")

	instanceId := "2f49c2b3"
	dataApiId := "afdb4e9d"
	mockHandler := helper.NewRequestHandlerMock(fmt.Sprintf("/v1/instances/%s/data-apis/graphql/%s", instanceId, dataApiId), http.StatusOK, `{
			"data": {
                "features": {
                        "subgraph": false
                },
                "id": "4f78488f",
                "name": "friendly-name",
                "status": "ready",
                "type_definitions": "dHlwZSBBY3RvciB7CiAgbmFtZTogU3RyaW5nCiAgbW92aWVzOiBbTW92aWUhXSEgQHJlbGF0aW9uc2hpcCh0eXBlOiAiQUNURURfSU4iLCBkaXJlY3Rpb246IE9VVCkKfQoKdHlwZSBNb3ZpZSB7CiAgdGl0bGU6IFN0cmluZwogIGFjdG9yczogW0FjdG9yIV0hIEByZWxhdGlvbnNoaXAodHlwZTogIkFDVEVEX0lOIiwgZGlyZWN0aW9uOiBJTikKfQ==",
                "url": "https://4f78488f.28be6e4d8d3e836019.graphql.neo4j.io/graphql"
        	}
		}`)

	helper.ExecuteCommand(fmt.Sprintf("data-api graphql get --output json --instance-id %s %s", instanceId, dataApiId))

	mockHandler.AssertCalledTimes(1)
	mockHandler.AssertCalledWithMethod(http.MethodGet)

	helper.AssertOutJson(`{
		"data": {
			"features": {
					"subgraph": false
			},
			"id": "4f78488f",
			"name": "friendly-name",
			"status": "ready",
			"type_definitions": "dHlwZSBBY3RvciB7CiAgbmFtZTogU3RyaW5nCiAgbW92aWVzOiBbTW92aWUhXSEgQHJlbGF0aW9uc2hpcCh0eXBlOiAiQUNURURfSU4iLCBkaXJlY3Rpb246IE9VVCkKfQoKdHlwZSBNb3ZpZSB7CiAgdGl0bGU6IFN0cmluZwogIGFjdG9yczogW0FjdG9yIV0hIEByZWxhdGlvbnNoaXAodHlwZTogIkFDVEVEX0lOIiwgZGlyZWN0aW9uOiBJTikKfQ==",
			"url": "https://4f78488f.28be6e4d8d3e836019.graphql.neo4j.io/graphql"
        }
	}`)
}

package graphql_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/neo4j/cli/neo4j-cli/aura/internal/test/testutils"
)

func TestListGraphQLDataApis(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	helper.SetConfigValue("aura.beta-enabled", true)

	instanceId := "2f49c2b3"
	mockHandler := helper.NewRequestHandlerMock(fmt.Sprintf("/v1beta5/instances/%s/data-apis/graphql", instanceId), http.StatusOK, `{
		"data": [
			{
				"id": "7261d20a",
				"name": "friendly-name",
				"status": "creating",
				"url": "https://23423.453489590fdsgs34.test.com/graphql"
			}
		]	
	}`)

	helper.ExecuteCommand(fmt.Sprintf("data-api graphql list --instance-id %s", instanceId))

	mockHandler.AssertCalledTimes(1)
	mockHandler.AssertCalledWithMethod(http.MethodGet)

	helper.AssertOutJson(`{
		"data": [
			{
				"id": "7261d20a",
				"name": "friendly-name",
				"status": "creating",
				"url": "https://23423.453489590fdsgs34.test.com/graphql"
			}
		]
	}
	`)
}

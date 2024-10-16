package instance_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/neo4j/cli/neo4j/aura/internal/test/testutils"
)

func TestUpdateMemory(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	instanceId := "2f49c2b3"

	mockHandler := helper.NewRequestHandlerMock(fmt.Sprintf("/v1/instances/%s", instanceId), http.StatusAccepted, `{
		"data": {
			"id": "2f49c2b3",
			"name": "Production",
			"status": "updating",
			"connection_url": "YOUR_CONNECTION_URL",
			"tenant_id": "YOUR_TENANT_ID",
			"cloud_provider": "gcp",
			"memory": "8GB",
			"region": "europe-west1",
			"type": "enterprise-db"
		}
	}`)

	helper.ExecuteCommand(fmt.Sprintf("instance update %s --memory 8GB", instanceId))

	mockHandler.AssertCalledTimes(1)
	mockHandler.AssertCalledWithMethod(http.MethodPatch)
	mockHandler.AssertCalledWithBody(`{"memory":"8GB"}`)

	helper.AssertOutJson(`{
	  "data": {
		"cloud_provider": "gcp",
		"connection_url": "YOUR_CONNECTION_URL",
		"id": "2f49c2b3",
		"memory": "8GB",
		"name": "Production",
		"region": "europe-west1",
		"status": "updating",
		"tenant_id": "YOUR_TENANT_ID",
		"type": "enterprise-db"
	  }
	}`)
}

func TestUpdateName(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	instanceId := "2f49c2b3"

	mockHandler := helper.NewRequestHandlerMock(fmt.Sprintf("/v1/instances/%s", instanceId), http.StatusOK, `{
		"data": {
			"id": "2f49c2b3",
			"name": "New Name",
			"status": "updating",
			"connection_url": "YOUR_CONNECTION_URL",
			"tenant_id": "YOUR_TENANT_ID",
			"cloud_provider": "gcp",
			"memory": "4GB",
			"region": "europe-west1",
			"type": "enterprise-db"
		}
	}`)

	helper.ExecuteCommand(fmt.Sprintf(`instance update %s --name "New Name"`, instanceId))

	mockHandler.AssertCalledTimes(1)
	mockHandler.AssertCalledWithMethod(http.MethodPatch)
	mockHandler.AssertCalledWithBody(`{"name":"New Name"}`)

	helper.AssertOutJson(`{
	  "data": {
		"cloud_provider": "gcp",
		"connection_url": "YOUR_CONNECTION_URL",
		"id": "2f49c2b3",
		"memory": "4GB",
		"name": "New Name",
		"region": "europe-west1",
		"status": "updating",
		"tenant_id": "YOUR_TENANT_ID",
		"type": "enterprise-db"
	  }
	}`)
}

func TestUpdateMemoryAndName(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	instanceId := "2f49c2b3"

	mockHandler := helper.NewRequestHandlerMock(fmt.Sprintf("/v1/instances/%s", instanceId), http.StatusAccepted, `{
		"data": {
			"id": "2f49c2b3",
			"name": "New Name",
			"status": "updating",
			"connection_url": "YOUR_CONNECTION_URL",
			"tenant_id": "YOUR_TENANT_ID",
			"cloud_provider": "gcp",
			"memory": "8GB",
			"region": "europe-west1",
			"type": "enterprise-db"
		}
	}`)

	helper.ExecuteCommand(fmt.Sprintf(`instance update %s --name "New Name" --memory 8GB`, instanceId))

	mockHandler.AssertCalledTimes(1)
	mockHandler.AssertCalledWithMethod(http.MethodPatch)
	mockHandler.AssertCalledWithBody(`{"memory":"8GB","name":"New Name"}`)

	helper.AssertOutJson(`{
	  "data": {
		"cloud_provider": "gcp",
		"connection_url": "YOUR_CONNECTION_URL",
		"id": "2f49c2b3",
		"memory": "8GB",
		"name": "New Name",
		"region": "europe-west1",
		"status": "updating",
		"tenant_id": "YOUR_TENANT_ID",
		"type": "enterprise-db"
	  }
	}`)
}

func TestUpdateErrorsWithNoFlags(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	instanceId := "2f49c2b3"

	mockHandler := helper.NewRequestHandlerMock(fmt.Sprintf("/v1/instances/%s", instanceId), http.StatusAccepted, "")

	helper.ExecuteCommand(fmt.Sprintf(`instance update %s`, instanceId))

	mockHandler.AssertCalledTimes(0)

	helper.AssertErr(`Error: at least one of the flags in the group [memory name] is required
`)
}

func TestUpdateInstanceError(t *testing.T) {
	testCases := []struct {
		statusCode    int
		expectedError string
		returnBody    string
	}{
		{
			statusCode:    http.StatusNotFound,
			expectedError: "Error: [DB not found: 24d18db5]",
			returnBody: `{
			"errors": [
			  {
				"message": "DB not found: 24d18db5",
				"reason": "db-not-found"
			  }
			]
		  }`,
		},
		{
			statusCode:    http.StatusConflict,
			expectedError: "Error: [The database is current undergoing an operation: resuming]",
			returnBody: `{
				"errors": [
				  {
					"message": "The database is current undergoing an operation: resuming",
					"reason": "ongoing-database-operation"
				  }
				]
			}`,
		},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("StatusCode%d", testCase.statusCode), func(t *testing.T) {
			helper := testutils.NewAuraTestHelper(t)
			defer helper.Close()

			instanceId := "2f49c2b3"

			mockHandler := helper.NewRequestHandlerMock(fmt.Sprintf("/v1/instances/%s", instanceId), testCase.statusCode, testCase.returnBody)

			helper.ExecuteCommand(fmt.Sprintf(`instance update %s --name "New Name" --memory 8GB`, instanceId))

			mockHandler.AssertCalledTimes(1)
			mockHandler.AssertCalledWithMethod(http.MethodPatch)

			helper.AssertOut("")
			helper.AssertErr(testCase.expectedError)
		})
	}
}

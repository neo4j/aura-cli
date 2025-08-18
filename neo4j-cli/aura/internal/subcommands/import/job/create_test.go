package job_test

import (
	"fmt"
	"github.com/neo4j/cli/neo4j-cli/aura/internal/test/testutils"
	"net/http"
	"testing"
)

func TestCreateImportJob(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	mockHandler := helper.NewRequestHandlerMock("/v2beta1/projects/f607bebe-0cc0-4166-b60c-b4eed69ee7ee/import/jobs", http.StatusCreated, `
		{
			"data": {"id": "87d485b4-73fc-4a7f-bb03-720f4672947e"}
		}
	`)

	helper.SetConfigValue("aura.beta-enabled", true)

	helper.ExecuteCommand("import job create --project-id=f607bebe-0cc0-4166-b60c-b4eed69ee7ee --import-model-id=e01cdc6d-2f50-4f46-b04b-8ec8fc8de839 --aura-db-id=07e49cf5")

	mockHandler.AssertCalledTimes(1)
	mockHandler.AssertCalledWithMethod(http.MethodPost)
	mockHandler.AssertCalledWithBody(`{
		"importModelId": "e01cdc6d-2f50-4f46-b04b-8ec8fc8de839",
		"auraCredentials": {
			"dbId": "07e49cf5"
		}
	}`)

	helper.AssertErr("")
	helper.AssertOutJson(`
		{
			"data": {"id": "87d485b4-73fc-4a7f-bb03-720f4672947e"}
		}
	`)
}

func TestCreateImportJobError(t *testing.T) {
	testCases := []struct {
		statusCode    int
		expectedError string
		returnBody    string
	}{
		{
			statusCode:    http.StatusBadRequest,
			expectedError: "Error: [DataSourceId: Import model data source id is required]",
			returnBody: `{
				"errors": [
					{
					"message": "Import model data source id is required",
					"reason": "Required",
					"field": "DataSourceId"
					}
				]
			}`,
		},
		{
			statusCode:    http.StatusMethodNotAllowed,
			expectedError: "Error: [string]",
			returnBody: `{
				"errors": [
					{
					"message": "string",
					"reason": "string",
					"field": "string"
					}
				]
			}`,
		},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("StatusCode%d", testCase.statusCode), func(t *testing.T) {
			helper := testutils.NewAuraTestHelper(t)
			defer helper.Close()

			mockHandler := helper.NewRequestHandlerMock("/v2beta1/projects/f607bebe-0cc0-4166-b60c-b4eed69ee7ee/import/jobs", testCase.statusCode, testCase.returnBody)

			helper.SetConfigValue("aura.beta-enabled", true)

			helper.ExecuteCommand("import job create --project-id=f607bebe-0cc0-4166-b60c-b4eed69ee7ee --import-model-id=e01cdc6d-2f50-4f46-b04b-8ec8fc8de839 --aura-db-id=07e49cf5")

			mockHandler.AssertCalledTimes(1)
			mockHandler.AssertCalledWithMethod(http.MethodPost)

			helper.AssertErr(testCase.expectedError)
		})
	}
}

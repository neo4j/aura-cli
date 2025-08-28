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
	projectId := "f607bebe-0cc0-4166-b60c-b4eed69ee7ee"
	importModelId := "e01cdc6d-2f50-4f46-b04b-8ec8fc8de839"
	auraDbId := "07e49cf5"
	testCases := map[string]struct {
		executeCommand      string
		expectedCalledTimes int
		statusCode          int
		expectedError       string
		returnBody          string
	}{
		"correct command with error response 1": {
			executeCommand:      fmt.Sprintf("import job create --project-id=%s --import-model-id=%s --aura-db-id=%s", projectId, importModelId, auraDbId),
			expectedCalledTimes: 1,
			statusCode:          http.StatusBadRequest,
			expectedError:       "Error: [DataSourceId: Import model data source id is required]",
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
		"correct command with error response 2": {
			executeCommand:      fmt.Sprintf("import job create --project-id=%s --import-model-id=%s --aura-db-id=%s", projectId, importModelId, auraDbId),
			expectedCalledTimes: 1,
			statusCode:          http.StatusMethodNotAllowed,
			expectedError:       "Error: [string]",
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
		"incorrect command with missing projectId": {
			executeCommand:      fmt.Sprintf("import job create --import-model-id=%s --aura-db-id=%s", importModelId, auraDbId),
			expectedCalledTimes: 0,
			statusCode:          http.StatusNotFound,
			expectedError:       "Error: required flag(s) \"project-id\" not set",
			returnBody:          ``,
		},
		"incorrect command with missing importModelId": {
			executeCommand: fmt.Sprintf("import job create --project-id=%s --aura-db-id=%s", projectId, auraDbId),
			statusCode:     http.StatusNotFound,
			expectedError:  "Error: required flag(s) \"import-model-id\" not set",
			returnBody:     ``,
		},
		"incorrect command with missing auraDbId": {
			executeCommand: fmt.Sprintf("import job create --project-id=%s --import-model-id=%s", projectId, importModelId),
			statusCode:     http.StatusNotFound,
			expectedError:  "Error: required flag(s) \"aura-db-id\" not set",
			returnBody:     ``,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			helper := testutils.NewAuraTestHelper(t)
			defer helper.Close()

			mockHandler := helper.NewRequestHandlerMock(fmt.Sprintf("/v2beta1/projects/%s/import/jobs", projectId), testCase.statusCode, testCase.returnBody)

			helper.SetConfigValue("aura.beta-enabled", true)

			helper.ExecuteCommand(testCase.executeCommand)

			mockHandler.AssertCalledTimes(testCase.expectedCalledTimes)
			if testCase.expectedCalledTimes > 0 {
				mockHandler.AssertCalledWithMethod(http.MethodPost)
			}

			helper.AssertErr(testCase.expectedError)
		})
	}
}

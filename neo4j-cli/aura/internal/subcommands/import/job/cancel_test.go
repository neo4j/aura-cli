package job_test

import (
	"fmt"
	"github.com/neo4j/cli/neo4j-cli/aura/internal/test/testutils"
	"net/http"
	"testing"
)

func TestCancelImportJob(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	mockHandler := helper.NewRequestHandlerMock("/v2beta1/projects/f607bebe-0cc0-4166-b60c-b4eed69ee7ee/import/jobs/87d485b4-73fc-4a7f-bb03-720f4672947e/cancel", http.StatusOK, `
		{
			"data": {"id": "87d485b4-73fc-4a7f-bb03-720f4672947e"}
		}
	`)

	helper.SetConfigValue("aura.beta-enabled", true)

	helper.ExecuteCommand("import job cancel --project-id=f607bebe-0cc0-4166-b60c-b4eed69ee7ee 87d485b4-73fc-4a7f-bb03-720f4672947e")

	mockHandler.AssertCalledTimes(1)
	mockHandler.AssertCalledWithMethod(http.MethodPatch)

	helper.AssertErr("")
	helper.AssertOutJson(`
		{
			"data": {"id": "87d485b4-73fc-4a7f-bb03-720f4672947e"}
		}
	`)
}

func TestCancelImportJobError(t *testing.T) {
	testCases := []struct {
		statusCode    int
		expectedError string
		returnBody    string
	}{
		{
			statusCode:    http.StatusBadRequest,
			expectedError: "Error: [The job 87d485b4-73fc-4a7f-bb03-720f4672947e has requested to cancel]",
			returnBody: `{
				"errors": [
					{
					"message": "The job 87d485b4-73fc-4a7f-bb03-720f4672947e has requested to cancel",
					"reason": "Requested"
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

			mockHandler := helper.NewRequestHandlerMock("/v2beta1/projects/f607bebe-0cc0-4166-b60c-b4eed69ee7ee/import/jobs/87d485b4-73fc-4a7f-bb03-720f4672947e/cancel", testCase.statusCode, testCase.returnBody)

			helper.SetConfigValue("aura.beta-enabled", true)

			helper.ExecuteCommand("import job cancel --project-id=f607bebe-0cc0-4166-b60c-b4eed69ee7ee 87d485b4-73fc-4a7f-bb03-720f4672947e")

			mockHandler.AssertCalledTimes(1)
			mockHandler.AssertCalledWithMethod(http.MethodPatch)

			helper.AssertErr(testCase.expectedError)
		})
	}
}

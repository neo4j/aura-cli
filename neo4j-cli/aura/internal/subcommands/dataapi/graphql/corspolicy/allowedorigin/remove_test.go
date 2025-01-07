package allowedorigin_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/neo4j/cli/neo4j-cli/aura/internal/test/testutils"
)

func TestRemoveAllowedOriginFlagsValidation(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	helper.SetConfigValue("aura.beta-enabled", true)

	instanceId := "2f49c2b3"
	dataApiId := "e157301d"
	allowedOrigin := "https://test.com"

	tests := map[string]struct {
		executedCommand string
		expectedError   string
	}{
		"missing all flags": {
			executedCommand: fmt.Sprintf("data-api graphql cors-policy allowed-origin remove %s", allowedOrigin),
			expectedError:   "Error: required flag(s) \"data-api-id\", \"instance-id\" not set",
		},
		"missing origin": {
			executedCommand: fmt.Sprintf("data-api graphql cors-policy allowed-origin remove --data-api-id %s --instance-id %s", dataApiId, instanceId),
			expectedError:   "Error: accepts 1 arg(s), received 0",
		},
		"missing data api id flag": {
			executedCommand: fmt.Sprintf("data-api graphql cors-policy allowed-origin remove %s --instance-id %s", allowedOrigin, instanceId),
			expectedError:   "Error: required flag(s) \"data-api-id\" not set",
		},
		"missing instance id flag": {
			executedCommand: fmt.Sprintf("data-api graphql cors-policy allowed-origin remove %s --data-api-id %s", allowedOrigin, dataApiId),
			expectedError:   "Error: required flag(s) \"instance-id\" not set",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			helper.ExecuteCommand(tt.executedCommand)
			helper.AssertErr(tt.expectedError)
		})
	}
}

func TestRemoveAllowedOriginWithResponse(t *testing.T) {
	instanceId := "2f49c2b3"
	dataApiId := "e157301d"
	allowedOrigin := "https://test.com"

	mockGetResponseNoOrigins := `{
		"data": {
			"id": "2f49c2b3",
			"name": "my-data-api-1",
			"status": "ready",
			"url": "https://2f49c2b3.28be6e4d8d3e8360197cb6c1fa1d25d1.graphql.neo4j-dev.io/graphql",
			"security": {
				"cors_policy": {
					"allowed_origins": []
				}
			}
		}
	}`

	mockGetResponseWithOrigins := fmt.Sprintf(`{
		"data": {
			"id": "2f49c2b3",
			"name": "my-data-api-1",
			"status": "ready",
			"url": "https://2f49c2b3.28be6e4d8d3e8360197cb6c1fa1d25d1.graphql.neo4j-dev.io/graphql",
			"security": {
				"cors_policy": {
					"allowed_origins": ["https://test1.com", "https://test2.com", "%s"]
				}
			}
		}
	}`, allowedOrigin)

	mockGetResponseWithLastExistingOrigin := fmt.Sprintf(`{
		"data": {
			"id": "2f49c2b3",
			"name": "my-data-api-1",
			"status": "ready",
			"url": "https://2f49c2b3.28be6e4d8d3e8360197cb6c1fa1d25d1.graphql.neo4j-dev.io/graphql",
			"security": {
				"cors_policy": {
					"allowed_origins": ["%s"]
				}
			}
		}
	}`, allowedOrigin)

	mockPatchResponse := `{
		"data": {
			"id": "2f49c2b3",
			"name": "my-data-api-1",
			"status": "ready",
			"url": "https://2f49c2b3.28be6e4d8d3e8360197cb6c1fa1d25d1.graphql.neo4j-dev.io/graphql"
		}
	}`

	expectedResponseExistingOrigin := `New allowed origins: ["https://test1.com", "https://test2.com"]
{
	"data": {
		"id": "2f49c2b3",
		"name": "my-data-api-1",
		"status": "ready",
		"url": "https://2f49c2b3.28be6e4d8d3e8360197cb6c1fa1d25d1.graphql.neo4j-dev.io/graphql"
	}
}`
	expectedResponseNoRemainingOrigins := `New allowed origins: []
{
	"data": {
		"id": "2f49c2b3",
		"name": "my-data-api-1",
		"status": "ready",
		"url": "https://2f49c2b3.28be6e4d8d3e8360197cb6c1fa1d25d1.graphql.neo4j-dev.io/graphql"
	}
}`

	expectedResponseTable := `New allowed origins: ["https://test1.com", "https://test2.com"]
┌──────────┬───────────────┬────────┬────────────────────────────────────────────────────────────────────────────────┐
│ ID       │ NAME          │ STATUS │ URL                                                                            │
├──────────┼───────────────┼────────┼────────────────────────────────────────────────────────────────────────────────┤
│ 2f49c2b3 │ my-data-api-1 │ ready  │ https://2f49c2b3.28be6e4d8d3e8360197cb6c1fa1d25d1.graphql.neo4j-dev.io/graphql │
└──────────┴───────────────┴────────┴────────────────────────────────────────────────────────────────────────────────┘
`

	tests := map[string]struct {
		mockGetResponse     string
		mockPatchResponse   string
		executeCommand      string
		expectedRequestBody string
		expectedResponse    string
		expectedErr         string
	}{
		"remove allowed origin successfully": {
			mockGetResponse:     mockGetResponseWithOrigins,
			mockPatchResponse:   mockPatchResponse,
			executeCommand:      fmt.Sprintf("data-api graphql cors-policy allowed-origin remove %s --instance-id %s --data-api-id %s", allowedOrigin, instanceId, dataApiId),
			expectedRequestBody: "{\"security\":{\"cors_policy\":{\"allowed_origins\":[\"https://test1.com\",\"https://test2.com\"]}}}",
			expectedResponse:    expectedResponseExistingOrigin,
		},
		"remove allowed origin with no existing origins": {
			mockGetResponse: mockGetResponseNoOrigins,
			executeCommand:  fmt.Sprintf("data-api graphql cors-policy allowed-origin remove %s --instance-id %s --data-api-id %s", allowedOrigin, instanceId, dataApiId),
			expectedErr:     fmt.Sprintf("Error: Origin \"%s\" not found in allowed origins", allowedOrigin),
		},
		"remove last allowed origin": {
			mockGetResponse:     mockGetResponseWithLastExistingOrigin,
			mockPatchResponse:   mockPatchResponse,
			executeCommand:      fmt.Sprintf("data-api graphql cors-policy allowed-origin remove %s --instance-id %s --data-api-id %s", allowedOrigin, instanceId, dataApiId),
			expectedRequestBody: "{\"security\":{\"cors_policy\":{\"allowed_origins\":[]}},\"test\":\"ignore me\"}",
			expectedResponse:    expectedResponseNoRemainingOrigins,
		},
		"remove allowed origin with output table": {
			mockGetResponse:     mockGetResponseWithOrigins,
			mockPatchResponse:   mockPatchResponse,
			executeCommand:      fmt.Sprintf("data-api graphql cors-policy allowed-origin remove %s --instance-id %s --data-api-id %s --output table", allowedOrigin, instanceId, dataApiId),
			expectedRequestBody: "{\"security\":{\"cors_policy\":{\"allowed_origins\":[\"https://test1.com\",\"https://test2.com\"]}}}",
			expectedResponse:    expectedResponseTable,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			helper := testutils.NewAuraTestHelper(t)
			defer helper.Close()

			helper.SetConfigValue("aura.beta-enabled", true)

			mockHandler := helper.NewRequestHandlerMock(fmt.Sprintf("/v1/instances/%s/data-apis/graphql/%s", instanceId, dataApiId), http.StatusOK, tt.mockGetResponse)
			mockHandler.AddResponse(http.StatusAccepted, tt.mockPatchResponse)

			helper.ExecuteCommand(tt.executeCommand)

			expectedCalls := 0
			if tt.mockPatchResponse != "" {
				expectedCalls += 1
			}
			if tt.mockGetResponse != "" {
				expectedCalls += 1
			}

			mockHandler.AssertCalledTimes(expectedCalls)
			if tt.mockGetResponse != "" {
				mockHandler.AssertCalledWithMethod(http.MethodGet)
			}
			if tt.mockPatchResponse != "" {
				mockHandler.AssertCalledWithMethod(http.MethodPatch)
				mockHandler.AssertCalledWithBody(tt.expectedRequestBody)
			}

			helper.AssertOut(tt.expectedResponse)
			helper.AssertErr(tt.expectedErr)
		})
	}
}

package authprovider_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/neo4j/cli/neo4j-cli/aura/internal/test/testutils"
)

func TestUpdateAuthProviderFlagsValidation(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	helper.SetConfigValue("aura.beta-enabled", true)

	instanceId := "2f49c2b3"
	dataApiId := "afdb4e9d"
	authProviderId := "929d4180-4d28-48fc-8e6e-6018613f0b52"

	tests := map[string]struct {
		executedCommand string
		expectedError   string
	}{
		"enabled flag has invalid value": {
			executedCommand: fmt.Sprintf("data-api graphql auth-provider update --output json --instance-id %s --data-api-id %s --enabled fg %s", instanceId, dataApiId, authProviderId),
			expectedError:   "Error: invalid value for boolean enabled, err: strconv.ParseBool: parsing \"fg\": invalid syntax",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			helper.ExecuteCommand(tt.executedCommand)
			helper.AssertErr(tt.expectedError)
		})
	}
}

func TestUpdateAuthProviderWithResponse(t *testing.T) {
	instanceId := "2f49c2b3"
	dataApiId := "75a234b5"
	authProviderId := "929d4180-4d28-48fc-8e6e-6018613f0b52"
	name := "api-key-32"
	enabled := "false"
	enabledShortHand := "t"
	url := "https://test.com/.well-known/jwks.json"

	mockResponse := `{
		"data": {
			"id": "929d4180-4d28-48fc-8e6e-6018613f0b52",
			"name": "api-key-32",
			"enabled": "false",
			"type": "jwks",
			"url": "https://test.com/.well-known/jwks.json"
		}
	}`
	expectedResponse := `{
		"data": {
			"enabled": "false",
			"id": "929d4180-4d28-48fc-8e6e-6018613f0b52",
			"name": "api-key-32",
			"type": "jwks",
			"url": "https://test.com/.well-known/jwks.json"
        }
	}`

	tests := map[string]struct {
		mockResponse        string
		executeCommand      string
		expectedRequestBody string
		expectedResponse    string
	}{
		"update the name": {
			mockResponse:        mockResponse,
			executeCommand:      fmt.Sprintf("data-api graphql auth-provider update --output json --instance-id %s --data-api-id %s --name %s %s", instanceId, dataApiId, name, authProviderId),
			expectedRequestBody: `{"name":"api-key-32"}`,
			expectedResponse:    expectedResponse,
		},
		"update enabled value": {
			mockResponse:        mockResponse,
			executeCommand:      fmt.Sprintf("data-api graphql auth-provider update --output json --instance-id %s --data-api-id %s --enabled %s %s", instanceId, dataApiId, enabled, authProviderId),
			expectedRequestBody: `{"enabled":false}`,
			expectedResponse:    expectedResponse,
		},
		"update enabled value short hand": {
			mockResponse:        mockResponse,
			executeCommand:      fmt.Sprintf("data-api graphql auth-provider update --output json --instance-id %s --data-api-id %s --enabled %s %s", instanceId, dataApiId, enabledShortHand, authProviderId),
			expectedRequestBody: `{"enabled":true}`,
			expectedResponse:    expectedResponse,
		},
		"update url": {
			mockResponse:        mockResponse,
			executeCommand:      fmt.Sprintf("data-api graphql auth-provider update --output json --instance-id %s --data-api-id %s --url %s %s", instanceId, dataApiId, url, authProviderId),
			expectedRequestBody: `{"url":"https://test.com/.well-known/jwks.json"}`,
			expectedResponse:    expectedResponse,
		},
		"update all possible values in one request": {
			mockResponse:        mockResponse,
			executeCommand:      fmt.Sprintf("data-api graphql auth-provider update --output json --instance-id %s --data-api-id %s --name %s --enabled %s --url %s %s", instanceId, dataApiId, name, enabled, url, authProviderId),
			expectedRequestBody: `{"enabled":false,"name":"api-key-32","url":"https://test.com/.well-known/jwks.json"}`,
			expectedResponse:    expectedResponse,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			helper := testutils.NewAuraTestHelper(t)
			defer helper.Close()

			helper.SetConfigValue("aura.beta-enabled", true)

			mockHandler := helper.NewRequestHandlerMock(fmt.Sprintf("/v1/instances/%s/data-apis/graphql/%s/auth-providers/%s", instanceId, dataApiId, authProviderId), http.StatusAccepted, tt.mockResponse)

			helper.ExecuteCommand(tt.executeCommand)

			mockHandler.AssertCalledTimes(1)
			mockHandler.AssertCalledWithMethod(http.MethodPatch)
			mockHandler.AssertCalledWithBody(tt.expectedRequestBody)

			helper.AssertOutJson(tt.expectedResponse)
		})
	}
}

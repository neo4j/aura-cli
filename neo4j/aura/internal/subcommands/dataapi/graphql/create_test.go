package graphql_test

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/neo4j/cli/neo4j/aura/internal/subcommands/dataapi/graphql"
	"github.com/neo4j/cli/neo4j/aura/internal/test/testutils"
	"github.com/stretchr/testify/assert"
)

func TestResolveTypeDefsFileFlagValue(t *testing.T) {
	// Create temporary files
	tmpGraphQLFile, err := os.CreateTemp("", "testTypeDefs.*.graphql")
	if err != nil {
		t.Fatalf("failed to create temp graphql file: %v", err)
	}
	tmpTextFile, err := os.CreateTemp("", "test.*.txt")
	if err != nil {
		t.Fatalf("failed to create temp text file: %v", err)
	}
	defer os.Remove(tmpGraphQLFile.Name()) // clean up the file
	defer os.Remove(tmpTextFile.Name())

	// Change permissions of the tmp graphql file to make the file unreadable
	if err := os.Chmod(tmpGraphQLFile.Name(), 0200); err != nil {
		t.Fatalf("failed to chmod file: %v", err)
	}

	tests := map[string]struct {
		flagValue        string
		expectedValue    string
		expectedErrorMsg string
	}{"correct path to file": {
		flagValue:        "../../../test/assets/typeDefs.graphql",
		expectedValue:    "dHlwZSBNb3ZpZSB7CiAgdGl0bGU6IFN0cmluZwp9Cg==",
		expectedErrorMsg: "",
	}, "invalid file type": {
		flagValue:        tmpTextFile.Name(),
		expectedValue:    "",
		expectedErrorMsg: "must have file type '.graphql'",
	},
		"invalid path": {
			flagValue:        "../../test/assets/typeDefs.graphql",
			expectedValue:    "",
			expectedErrorMsg: "type definitions file '../../test/assets/typeDefs.graphql' does not exist",
		}, "empty file": {
			flagValue:        "../../../test/assets/empty.graphql",
			expectedValue:    "",
			expectedErrorMsg: "read type definitions file is empty",
		}, "unreadable file": {
			flagValue:        tmpGraphQLFile.Name(),
			expectedValue:    "",
			expectedErrorMsg: "reading type definitions file failed with error",
		}}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			val, err := graphql.ResolveTypeDefsFileFlagValue(test.flagValue)
			if test.expectedErrorMsg != "" {
				assert.Contains(t, err.Error(), test.expectedErrorMsg)
			} else {
				assert.Equal(t, test.expectedValue, val)
			}
		})
	}
}

func TestCreateGraphQLDataApiFlagsValidation(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	helper.SetConfigValue("aura.beta-enabled", true)

	instanceId := "2f49c2b3"
	instanceUsername := "neo4j"
	instancePassword := "dfjglhssdopfrow"
	name := "my-data-api-1"
	typeDefs := "dHlwZSBNb3ZpZSB7CiAgdGl0bGU6IFN0cmluZwp9"
	invalidBase64TypeDefs := "df"
	typeDefsFile := "../../../test/assets/typeDefs.graphql"
	invalidTypeDefsFile := "../invalid/typeDefs.graphql"

	tests := map[string]struct {
		executedCommand string
		expectedError   string
	}{"no type defs flag provided": {
		executedCommand: fmt.Sprintf("data-api graphql create --instance-id %s", instanceId),
		expectedError:   "Error: either '--type-definitions' or '--type-definitions-file' flag needs to be provided",
	},
		"only one type defs flag can be provided": {
			executedCommand: fmt.Sprintf("data-api graphql create --instance-id %s --type-definitions %s --type-definitions-file %s", instanceId, typeDefs, typeDefsFile),
			expectedError:   "Error: only one of '--type-definitions' or '--type-definitions-file' flag can be provided",
		},
		"missing almost all flags": {
			executedCommand: fmt.Sprintf("data-api graphql create --instance-id %s --type-definitions %s", instanceId, typeDefs),
			expectedError:   "Error: required flag(s) \"instance-password\", \"instance-username\", \"name\" not set",
		},
		"missing instance password flag": {
			executedCommand: fmt.Sprintf("data-api graphql create --instance-id %s --instance-username %s --name %s --type-definitions %s", instanceId, instanceUsername, name, typeDefs),
			expectedError:   "Error: required flag(s) \"instance-password\" not set",
		},
		"missing instance username flag": {
			executedCommand: fmt.Sprintf("data-api graphql create --instance-id %s --instance-password %s --name %s --type-definitions %s", instanceId, instancePassword, name, typeDefs),
			expectedError:   "Error: required flag(s) \"instance-username\" not set",
		},
		"missing name flag": {
			executedCommand: fmt.Sprintf("data-api graphql create --instance-id %s --instance-username %s --instance-password %s --type-definitions %s", instanceId, instanceUsername, instancePassword, typeDefs),
			expectedError:   "Error: required flag(s) \"name\" not set",
		},
		"invalid base64 for type defs": {
			executedCommand: fmt.Sprintf("data-api graphql create --instance-id %s --instance-username %s --instance-password %s --name %s --type-definitions %s", instanceId, instanceUsername, instancePassword, name, invalidBase64TypeDefs),
			expectedError:   "Error: provided type definitions are not valid base64",
		},
		"invalid type defs file": {
			executedCommand: fmt.Sprintf("data-api graphql create --instance-id %s --instance-username %s --instance-password %s --name %s --type-definitions-file %s", instanceId, instanceUsername, instancePassword, name, invalidTypeDefsFile),
			expectedError:   "Error: type definitions file '../invalid/typeDefs.graphql' does not exist",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			helper.ExecuteCommand(tt.executedCommand)
			helper.AssertErr(tt.expectedError)
		})
	}
}

func TestCreateGraphQLDataApiWithResponse(t *testing.T) {
	instanceId := "2f49c2b3"
	instanceUsername := "neo4j"
	instancePassword := "dfjglhssdopfrow"
	name := "my-data-api-1"
	typeDefs := "dHlwZSBNb3ZpZSB7CiAgdGl0bGU6IFN0cmluZwp9"
	typeDefsFile := "../../../test/assets/typeDefs.graphql"

	mockResponse := `{
		"data": {
			"id": "2f49c2b3",
			"name": "my-data-api-1",
			"status": "creating",
			"url": "https://2f49c2b3.28be6e4d8d3e8360197cb6c1fa1d25d1.graphql.neo4j-dev.io/graphql",
			"authentication_providers": [
				{
					"id": "1ad1b794-e40e-41f7-8e8c-5638130317ed",
					"name": "default",
					"type": "api-key",
					"enabled": true,
					"key": "ublHwKxm2ylsc1HlkuL8NAcMfZnEVP1g"
				}
			]
		}
	}`

	expectedResponseJson := `###############################
# An API key was created. It is important to _store_ the API key as it is not currently possible to get it or update it.
#
# If you lose your API key, you will need to create a new Authentication provider.
# This will not result in any loss of data.
###############################
{
	"data": {
		"authentication_providers": [
			{
				"enabled": true,
				"id": "1ad1b794-e40e-41f7-8e8c-5638130317ed",
				"key": "ublHwKxm2ylsc1HlkuL8NAcMfZnEVP1g",
				"name": "default",
				"type": "api-key"
			}
		],
		"id": "2f49c2b3",
		"name": "my-data-api-1",
		"status": "creating",
		"url": "https://2f49c2b3.28be6e4d8d3e8360197cb6c1fa1d25d1.graphql.neo4j-dev.io/graphql"
	}
}`
	expectedResponseTable := `###############################
# An API key was created. It is important to _store_ the API key as it is not currently possible to get it or update it.
#
# If you lose your API key, you will need to create a new Authentication provider.
# This will not result in any loss of data.
###############################
┌──────────┬───────────────┬──────────┬────────────────────────────────────────────────────────────────────────────────┬───────────────────────────────────────────────────┐
│ ID       │ NAME          │ STATUS   │ URL                                                                            │ AUTHENTICATION_PROVIDERS                          │
├──────────┼───────────────┼──────────┼────────────────────────────────────────────────────────────────────────────────┼───────────────────────────────────────────────────┤
│ 2f49c2b3 │ my-data-api-1 │ creating │ https://2f49c2b3.28be6e4d8d3e8360197cb6c1fa1d25d1.graphql.neo4j-dev.io/graphql │ [                                                 │
│          │               │          │                                                                                │   {                                               │
│          │               │          │                                                                                │     "enabled": true,                              │
│          │               │          │                                                                                │     "id": "1ad1b794-e40e-41f7-8e8c-5638130317ed", │
│          │               │          │                                                                                │     "key": "ublHwKxm2ylsc1HlkuL8NAcMfZnEVP1g",    │
│          │               │          │                                                                                │     "name": "default",                            │
│          │               │          │                                                                                │     "type": "api-key"                             │
│          │               │          │                                                                                │   }                                               │
│          │               │          │                                                                                │ ]                                                 │
└──────────┴───────────────┴──────────┴────────────────────────────────────────────────────────────────────────────────┴───────────────────────────────────────────────────┘
	`

	tests := map[string]struct {
		mockResponse        string
		executeCommand      string
		expectedRequestBody string
		expectedResponse    string
	}{"create with default auth provider": {
		mockResponse:        mockResponse,
		executeCommand:      fmt.Sprintf("data-api graphql create --instance-id %s --instance-username %s --instance-password %s --name %s --type-definitions %s", instanceId, instanceUsername, instancePassword, name, typeDefs),
		expectedRequestBody: `{"aura_instance":{"password":"dfjglhssdopfrow","username":"neo4j"},"name":"my-data-api-1","security":{"authentication_providers":[{"enabled":true,"name":"default","type":"api-key"}]},"type_definitions":"dHlwZSBNb3ZpZSB7CiAgdGl0bGU6IFN0cmluZwp9"}`,
		expectedResponse:    expectedResponseJson,
	}, "create with default auth provider and output as table": {
		mockResponse:        mockResponse,
		executeCommand:      fmt.Sprintf("data-api graphql create --output table --instance-id %s --instance-username %s --instance-password %s --name %s --type-definitions %s ", instanceId, instanceUsername, instancePassword, name, typeDefs),
		expectedRequestBody: `{"aura_instance":{"password":"dfjglhssdopfrow","username":"neo4j"},"name":"my-data-api-1","security":{"authentication_providers":[{"enabled":true,"name":"default","type":"api-key"}]},"type_definitions":"dHlwZSBNb3ZpZSB7CiAgdGl0bGU6IFN0cmluZwp9"}`,
		expectedResponse:    expectedResponseTable,
	},
		"providing type defs as file": {
			mockResponse:        mockResponse,
			executeCommand:      fmt.Sprintf("data-api graphql create --instance-id %s --instance-username %s --instance-password %s --name %s --type-definitions-file %s", instanceId, instanceUsername, instancePassword, name, typeDefsFile),
			expectedRequestBody: `{"aura_instance":{"password":"dfjglhssdopfrow","username":"neo4j"},"name":"my-data-api-1","security":{"authentication_providers":[{"enabled":true,"name":"default","type":"api-key"}]},"type_definitions":"dHlwZSBNb3ZpZSB7CiAgdGl0bGU6IFN0cmluZwp9Cg=="}`,
			expectedResponse:    expectedResponseJson,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			helper := testutils.NewAuraTestHelper(t)
			defer helper.Close()

			helper.SetConfigValue("aura.beta-enabled", true)

			mockHandler := helper.NewRequestHandlerMock(fmt.Sprintf("/v1/instances/%s/data-apis/graphql", instanceId), http.StatusAccepted, tt.mockResponse)

			helper.ExecuteCommand(tt.executeCommand)

			mockHandler.AssertCalledTimes(1)
			mockHandler.AssertCalledWithMethod(http.MethodPost)
			mockHandler.AssertCalledWithBody(tt.expectedRequestBody)

			helper.AssertOut(tt.expectedResponse)
		})
	}
}

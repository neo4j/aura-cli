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
	// Create a temporary file
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

	// Change permissions to make the file unreadable
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

	helper.SetConfigValue("aura.beta-enabled", "true")

	instanceId := "2f49c2b3"
	instanceUsername := "neo4j"
	instancePassword := "dfjglhssdopfrow"
	name := "my-data-api-1"
	typeDefs := "dHlwZSBBY3RvciB7CiAgbmFtZTogU3RyaW5nCiAgbW92aWVzOiBbTW92aWUhXSEgQHJlbGF0aW9uc2hpcCh0eXBlOiAiQUNURURfSU4iLCBkaXJlY3Rpb246IE9VVCkKfQoKdHlwZSBNb3ZpZSB7CiAgdGl0bGU6IFN0cmluZwogIGFjdG9yczogW0FjdG9yIV0hIEByZWxhdGlvbnNoaXAodHlwZTogIkFDVEVEX0lOIiwgZGlyZWN0aW9uOiBJTikKfQ=="
	invalidTypeDefs := "df"
	typeDefsFile := "../../../test/assets/typeDefs.graphql"
	invalidTypeDefsFile := "../invalid/typeDefs.graphql"
	secAuthProviderName := "provider-1"
	secAuthProviderTypeApiKey := "api-key"
	secAuthProviderTypeJwks := "jwks"
	invalidSecAuthProviderType := "non-existing-type"

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
			expectedError:   "Error: required flag(s) \"instance-password\", \"instance-username\", \"name\", \"security-auth-provider-name\", \"security-auth-provider-type\" not set",
		},
		"missing instance password flag": {
			executedCommand: fmt.Sprintf("data-api graphql create --instance-id %s --instance-username %s --name %s --type-definitions %s --security-auth-provider-name %s --security-auth-provider-type %s", instanceId, instanceUsername, name, typeDefs, secAuthProviderName, secAuthProviderTypeApiKey),
			expectedError:   "Error: required flag(s) \"instance-password\" not set",
		},
		"missing instance username flag": {
			executedCommand: fmt.Sprintf("data-api graphql create --instance-id %s --instance-password %s --name %s --type-definitions %s --security-auth-provider-name %s --security-auth-provider-type %s", instanceId, instancePassword, name, typeDefs, secAuthProviderName, secAuthProviderTypeApiKey),
			expectedError:   "Error: required flag(s) \"instance-username\" not set",
		},
		"missing name flag": {
			executedCommand: fmt.Sprintf("data-api graphql create --instance-id %s --instance-username %s --instance-password %s --type-definitions %s --security-auth-provider-name %s --security-auth-provider-type %s", instanceId, instanceUsername, instancePassword, typeDefs, secAuthProviderName, secAuthProviderTypeApiKey),
			expectedError:   "Error: required flag(s) \"name\" not set",
		},
		"missing auth provider name flag": {
			executedCommand: fmt.Sprintf("data-api graphql create --instance-id %s --instance-username %s --instance-password %s --name %s --type-definitions %s --security-auth-provider-type %s", instanceId, instanceUsername, instancePassword, name, typeDefs, secAuthProviderTypeApiKey),
			expectedError:   "Error: required flag(s) \"security-auth-provider-name\" not set",
		},
		"missing auth provider type flag": {
			executedCommand: fmt.Sprintf("data-api graphql create --instance-id %s --instance-username %s --instance-password %s --name %s --type-definitions %s --security-auth-provider-name %s", instanceId, instanceUsername, instancePassword, name, typeDefs, secAuthProviderName),
			expectedError:   "Error: required flag(s) \"security-auth-provider-type\" not set",
		},
		"invalid type defs": {
			executedCommand: fmt.Sprintf("data-api graphql create --instance-id %s --instance-username %s --instance-password %s --name %s --type-definitions %s --security-auth-provider-name %s --security-auth-provider-type %s", instanceId, instanceUsername, instancePassword, name, invalidTypeDefs, secAuthProviderName, secAuthProviderTypeApiKey),
			expectedError:   "Error: provided type definitions are not valid base64",
		},
		"invalid type defs file": {
			executedCommand: fmt.Sprintf("data-api graphql create --instance-id %s --instance-username %s --instance-password %s --name %s --type-definitions-file %s --security-auth-provider-name %s --security-auth-provider-type %s", instanceId, instanceUsername, instancePassword, name, invalidTypeDefsFile, secAuthProviderName, secAuthProviderTypeApiKey),
			expectedError:   "Error: type definitions file '../invalid/typeDefs.graphql' does not exist",
		},
		"invalid auth provider type": {
			executedCommand: fmt.Sprintf("data-api graphql create --instance-id %s --instance-username %s --instance-password %s --name %s --type-definitions %s --security-auth-provider-name %s --security-auth-provider-type %s", instanceId, instanceUsername, instancePassword, name, typeDefs, secAuthProviderName, invalidSecAuthProviderType),
			expectedError:   "Error: invalid security auth provider type, got 'non-existing-type', expect 'api-key' or 'jwks'",
		},
		"auth provider jwks missing url": {
			executedCommand: fmt.Sprintf("data-api graphql create --instance-id %s --instance-username %s --instance-password %s --name %s --type-definitions %s --security-auth-provider-name %s --security-auth-provider-type %s", instanceId, instanceUsername, instancePassword, name, typeDefs, secAuthProviderName, secAuthProviderTypeJwks),
			expectedError:   "Error: required flag(s) \"security-auth-provider-url\" not set",
		},
		"invalid bool value for feature subgraph enabled": {
			executedCommand: fmt.Sprintf("data-api graphql create --instance-id %s --instance-username %s --instance-password %s --name %s --type-definitions %s --security-auth-provider-name %s --security-auth-provider-type %s --feature-subgraph-enabled yaya", instanceId, instanceUsername, instancePassword, name, typeDefs, secAuthProviderName, secAuthProviderTypeApiKey),
			expectedError:   "Error: strconv.ParseBool: parsing \"yaya\": invalid syntax",
		},
		"invalid bool value for auth provider enabled": {
			executedCommand: fmt.Sprintf("data-api graphql create --instance-id %s --instance-username %s --instance-password %s --name %s --type-definitions %s --security-auth-provider-name %s --security-auth-provider-type %s --security-auth-provider-enabled yeye", instanceId, instanceUsername, instancePassword, name, typeDefs, secAuthProviderName, secAuthProviderTypeApiKey),
			expectedError:   "Error: strconv.ParseBool: parsing \"yeye\": invalid syntax",
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
	typeDefs := "dHlwZSBBY3RvciB7CiAgbmFtZTogU3RyaW5nCiAgbW92aWVzOiBbTW92aWUhXSEgQHJlbGF0aW9uc2hpcCh0eXBlOiAiQUNURURfSU4iLCBkaXJlY3Rpb246IE9VVCkKfQoKdHlwZSBNb3ZpZSB7CiAgdGl0bGU6IFN0cmluZwogIGFjdG9yczogW0FjdG9yIV0hIEByZWxhdGlvbnNoaXAodHlwZTogIkFDVEVEX0lOIiwgZGlyZWN0aW9uOiBJTikKfQ=="
	typeDefsFile := "../../../test/assets/typeDefs.graphql"
	secAuthProviderName := "provider-1"
	secAuthProviderTypeApiKey := "api-key"
	secAuthProviderTypeJwks := "jwks"
	secAuthProviderUrl := "https://test.com/.well-known/jwks.json"

	tests := map[string]struct {
		mockResponse        string
		executeCommand      string
		expectedRequestBody string
		expectedResponse    string
	}{"one auth provider - api-key": {
		mockResponse: `{
		"data": {
			"id": "2f49c2b3",
			"name": "my-data-api-1",
			"status": "creating",
			"url": "https://2f49c2b3.28be6e4d8d3e8360197cb6c1fa1d25d1.graphql.neo4j-dev.io/graphql",
			"authentication_providers": [
				{
					"id": "1ad1b794-e40e-41f7-8e8c-5638130317ed",
					"name": "provider-1",
					"type": "api-key",
					"enabled": true,
					"key": "ublHwKxm2ylsc1HlkuL8NAcMfZnEVP1g"
				}
			]
		}
	}`,
		executeCommand:      fmt.Sprintf("data-api graphql create --instance-id %s --instance-username %s --instance-password %s --name %s --type-definitions %s --security-auth-provider-name %s --security-auth-provider-type %s --security-auth-provider-enabled false --feature-subgraph-enabled true", instanceId, instanceUsername, instancePassword, name, typeDefs, secAuthProviderName, secAuthProviderTypeApiKey),
		expectedRequestBody: `{"aura_instance":{"password":"dfjglhssdopfrow","username":"neo4j"},"features":{"subgraph":true},"name":"my-data-api-1","security":{"authentication_providers":[{"enabled":false,"name":"provider-1","type":"api-key"}]},"type_definitions":"dHlwZSBBY3RvciB7CiAgbmFtZTogU3RyaW5nCiAgbW92aWVzOiBbTW92aWUhXSEgQHJlbGF0aW9uc2hpcCh0eXBlOiAiQUNURURfSU4iLCBkaXJlY3Rpb246IE9VVCkKfQoKdHlwZSBNb3ZpZSB7CiAgdGl0bGU6IFN0cmluZwogIGFjdG9yczogW0FjdG9yIV0hIEByZWxhdGlvbnNoaXAodHlwZTogIkFDVEVEX0lOIiwgZGlyZWN0aW9uOiBJTikKfQ=="}`,
		expectedResponse: `###############################
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
				"name": "provider-1",
				"type": "api-key"
			}
		],
		"id": "2f49c2b3",
		"name": "my-data-api-1",
		"status": "creating",
		"url": "https://2f49c2b3.28be6e4d8d3e8360197cb6c1fa1d25d1.graphql.neo4j-dev.io/graphql"
	}
}`,
	}, "one auth provider - api-key - output as table": {
		mockResponse: `{
		"data": {
			"id": "2f49c2b3",
			"name": "my-data-api-1",
			"status": "creating",
			"url": "https://2f49c2b3.28be6e4d8d3e8360197cb6c1fa1d25d1.graphql.neo4j-dev.io/graphql",
			"authentication_providers": [
				{
					"id": "1ad1b794-e40e-41f7-8e8c-5638130317ed",
					"name": "provider-1",
					"type": "api-key",
					"enabled": true,
					"key": "ublHwKxm2ylsc1HlkuL8NAcMfZnEVP1g"
				}
			]
		}
	}`,
		executeCommand:      fmt.Sprintf("data-api graphql create --output table --instance-id %s --instance-username %s --instance-password %s --name %s --type-definitions %s --security-auth-provider-name %s --security-auth-provider-type %s", instanceId, instanceUsername, instancePassword, name, typeDefs, secAuthProviderName, secAuthProviderTypeApiKey),
		expectedRequestBody: `{"aura_instance":{"password":"dfjglhssdopfrow","username":"neo4j"},"features":{"subgraph":false},"name":"my-data-api-1","security":{"authentication_providers":[{"enabled":true,"name":"provider-1","type":"api-key"}]},"type_definitions":"dHlwZSBBY3RvciB7CiAgbmFtZTogU3RyaW5nCiAgbW92aWVzOiBbTW92aWUhXSEgQHJlbGF0aW9uc2hpcCh0eXBlOiAiQUNURURfSU4iLCBkaXJlY3Rpb246IE9VVCkKfQoKdHlwZSBNb3ZpZSB7CiAgdGl0bGU6IFN0cmluZwogIGFjdG9yczogW0FjdG9yIV0hIEByZWxhdGlvbnNoaXAodHlwZTogIkFDVEVEX0lOIiwgZGlyZWN0aW9uOiBJTikKfQ=="}`,
		expectedResponse: `###############################
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
│          │               │          │                                                                                │     "name": "provider-1",                         │
│          │               │          │                                                                                │     "type": "api-key"                             │
│          │               │          │                                                                                │   }                                               │
│          │               │          │                                                                                │ ]                                                 │
└──────────┴───────────────┴──────────┴────────────────────────────────────────────────────────────────────────────────┴───────────────────────────────────────────────────┘
	`,
	},
		"one auth provider - jwks": {
			mockResponse: `{
		"data": {
			"id": "2f49c2b3",
			"name": "my-data-api-1",
			"status": "creating",
			"url": "https://2f49c2b3.28be6e4d8d3e8360197cb6c1fa1d25d1.graphql.neo4j-dev.io/graphql",
			"authentication_providers": [
				{
					"id": "5170b65b-1ea6-4d59-8df6-7fd02b77fc75",
					"name": "provider-1",
					"type": "jwks",
					"enabled": true,
					"url": "https://test.com/.well-known/jwks.json"
				}
			]
		}
	}`,
			executeCommand:      fmt.Sprintf("data-api graphql create --instance-id %s --instance-username %s --instance-password %s --name %s --type-definitions %s --security-auth-provider-name %s --security-auth-provider-type %s --security-auth-provider-url %s --feature-subgraph-enabled false", instanceId, instanceUsername, instancePassword, name, typeDefs, secAuthProviderName, secAuthProviderTypeJwks, secAuthProviderUrl),
			expectedRequestBody: `{"aura_instance":{"password":"dfjglhssdopfrow","username":"neo4j"},"features":{"subgraph":false},"name":"my-data-api-1","security":{"authentication_providers":[{"enabled":true,"name":"provider-1","type":"jwks","url":"https://test.com/.well-known/jwks.json"}]},"type_definitions":"dHlwZSBBY3RvciB7CiAgbmFtZTogU3RyaW5nCiAgbW92aWVzOiBbTW92aWUhXSEgQHJlbGF0aW9uc2hpcCh0eXBlOiAiQUNURURfSU4iLCBkaXJlY3Rpb246IE9VVCkKfQoKdHlwZSBNb3ZpZSB7CiAgdGl0bGU6IFN0cmluZwogIGFjdG9yczogW0FjdG9yIV0hIEByZWxhdGlvbnNoaXAodHlwZTogIkFDVEVEX0lOIiwgZGlyZWN0aW9uOiBJTikKfQ=="}`,
			expectedResponse: `{
	"data": {
		"authentication_providers": [
			{
				"enabled": true,
				"id": "5170b65b-1ea6-4d59-8df6-7fd02b77fc75",
				"name": "provider-1",
				"type": "jwks",
				"url": "https://test.com/.well-known/jwks.json"
			}
		],
		"id": "2f49c2b3",
		"name": "my-data-api-1",
		"status": "creating",
		"url": "https://2f49c2b3.28be6e4d8d3e8360197cb6c1fa1d25d1.graphql.neo4j-dev.io/graphql"
	}
}`,
		},
		"providing type defs as file": {
			mockResponse: `{
		"data": {
			"id": "2f49c2b3",
			"name": "my-data-api-1",
			"status": "creating",
			"url": "https://2f49c2b3.28be6e4d8d3e8360197cb6c1fa1d25d1.graphql.neo4j-dev.io/graphql",
			"authentication_providers": [
				{
					"id": "5170b65b-1ea6-4d59-8df6-7fd02b77fc75",
					"name": "provider-1",
					"type": "jwks",
					"enabled": true,
					"url": "https://test.com/.well-known/jwks.json"
				}
			]
		}
	}`,
			executeCommand:      fmt.Sprintf("data-api graphql create --instance-id %s --instance-username %s --instance-password %s --name %s --type-definitions-file %s --security-auth-provider-name %s --security-auth-provider-type %s --security-auth-provider-url %s --feature-subgraph-enabled false", instanceId, instanceUsername, instancePassword, name, typeDefsFile, secAuthProviderName, secAuthProviderTypeJwks, secAuthProviderUrl),
			expectedRequestBody: `{"aura_instance":{"password":"dfjglhssdopfrow","username":"neo4j"},"features":{"subgraph":false},"name":"my-data-api-1","security":{"authentication_providers":[{"enabled":true,"name":"provider-1","type":"jwks","url":"https://test.com/.well-known/jwks.json"}]},"type_definitions":"dHlwZSBNb3ZpZSB7CiAgdGl0bGU6IFN0cmluZwp9Cg=="}`,
			expectedResponse: `{
	"data": {
		"authentication_providers": [
			{
				"enabled": true,
				"id": "5170b65b-1ea6-4d59-8df6-7fd02b77fc75",
				"name": "provider-1",
				"type": "jwks",
				"url": "https://test.com/.well-known/jwks.json"
			}
		],
		"id": "2f49c2b3",
		"name": "my-data-api-1",
		"status": "creating",
		"url": "https://2f49c2b3.28be6e4d8d3e8360197cb6c1fa1d25d1.graphql.neo4j-dev.io/graphql"
	}
}`,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			helper := testutils.NewAuraTestHelper(t)
			defer helper.Close()

			helper.SetConfigValue("aura.beta-enabled", "true")

			mockHandler := helper.NewRequestHandlerMock(fmt.Sprintf("/v1/instances/%s/data-apis/graphql", instanceId), http.StatusOK, tt.mockResponse)

			helper.ExecuteCommand(tt.executeCommand)

			mockHandler.AssertCalledTimes(1)
			mockHandler.AssertCalledWithMethod(http.MethodPost)
			mockHandler.AssertCalledWithBody(tt.expectedRequestBody)

			helper.AssertOut(tt.expectedResponse)
		})
	}
}

package graphql_test

import (
	"errors"
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
	tmpFile, err := os.CreateTemp("", "invalid-file")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name()) // clean up the file

	// Change permissions to make the file unreadable
	if err := os.Chmod(tmpFile.Name(), 0200); err != nil {
		t.Fatalf("failed to chmod file: %v", err)
	}

	tests := map[string]struct {
		flagValue     string
		expectedValue string
		expectedError error
	}{"correct path to file": {
		flagValue:     "../../../test/assets/typeDefs.graphql",
		expectedValue: "dHlwZSBNb3ZpZSB7CiAgdGl0bGU6IFN0cmluZwp9Cg==",
		expectedError: nil,
	}, "invalid path": {
		flagValue:     "../../test/assets/typeDefs.graphql",
		expectedValue: "",
		expectedError: errors.New("type definitions file '../../test/assets/typeDefs.graphql' does not exist"),
	}, "empty file": {
		flagValue:     "../../../test/assets/empty.graphql",
		expectedValue: "",
		expectedError: errors.New("read type definitions file is empty"),
	}, "unreadable file": {
		flagValue:     tmpFile.Name(),
		expectedValue: "",
		expectedError: errors.New("reading type definitions file failed with error"),
	}}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			val, err := graphql.ResolveTypeDefsFileFlagValue(test.flagValue)
			if test.expectedError != nil {
				assert.ErrorContains(t, err, test.expectedError.Error())
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
		mockResponse     string
		executeCommand   string
		expectedResponse string
		isJsonResponse   bool
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
		executeCommand: fmt.Sprintf("data-api graphql create --instance-id %s --instance-username %s --instance-password %s --name %s --type-definitions %s --security-auth-provider-name %s --security-auth-provider-type %s --feature-subgraph-enabled true", instanceId, instanceUsername, instancePassword, name, typeDefs, secAuthProviderName, secAuthProviderTypeApiKey),
		expectedResponse: `{
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
		isJsonResponse: true,
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
		executeCommand: fmt.Sprintf("data-api graphql create --output table --instance-id %s --instance-username %s --instance-password %s --name %s --type-definitions %s --security-auth-provider-name %s --security-auth-provider-type %s --feature-subgraph-enabled true", instanceId, instanceUsername, instancePassword, name, typeDefs, secAuthProviderName, secAuthProviderTypeApiKey),
		expectedResponse: `
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
		isJsonResponse: false,
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
			executeCommand: fmt.Sprintf("data-api graphql create --instance-id %s --instance-username %s --instance-password %s --name %s --type-definitions %s --security-auth-provider-name %s --security-auth-provider-type %s --security-auth-provider-url %s --feature-subgraph-enabled false", instanceId, instanceUsername, instancePassword, name, typeDefs, secAuthProviderName, secAuthProviderTypeJwks, secAuthProviderUrl),
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
			isJsonResponse: true,
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
			executeCommand: fmt.Sprintf("data-api graphql create --instance-id %s --instance-username %s --instance-password %s --name %s --type-definitions-file %s --security-auth-provider-name %s --security-auth-provider-type %s --security-auth-provider-url %s --feature-subgraph-enabled false", instanceId, instanceUsername, instancePassword, name, typeDefsFile, secAuthProviderName, secAuthProviderTypeJwks, secAuthProviderUrl),
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
			isJsonResponse: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if name != "providing type defs as file" {
				return
			}

			helper := testutils.NewAuraTestHelper(t)
			defer helper.Close()

			helper.SetConfigValue("aura.beta-enabled", "true")

			mockHandler := helper.NewRequestHandlerMock(fmt.Sprintf("/v1/instances/%s/data-apis/graphql", instanceId), http.StatusOK, tt.mockResponse)

			helper.ExecuteCommand(tt.executeCommand)

			mockHandler.AssertCalledTimes(1)
			mockHandler.AssertCalledWithMethod(http.MethodPost)

			if tt.isJsonResponse {
				helper.AssertOutJson(tt.expectedResponse)
			} else {
				helper.AssertOut(tt.expectedResponse)
			}
		})
	}
}

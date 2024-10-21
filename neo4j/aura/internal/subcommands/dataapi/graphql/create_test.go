package graphql_test

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/neo4j/cli/neo4j/aura/internal/subcommands/dataapi/graphql"
	"github.com/neo4j/cli/neo4j/aura/internal/test/testutils"
	"github.com/stretchr/testify/assert"
)

func TestResolveTypeDefsFlagValue(t *testing.T) {
	tests := map[string]struct {
		flagValue     string
		expectedValue string
		expectedError error
	}{"correct base64 string": {
		flagValue:     "TXkgc3RyaW5n",
		expectedValue: "TXkgc3RyaW5n",
		expectedError: nil,
	}, "correct path to file": {
		flagValue:     "../../../test/testutils/typeDefs.graphql",
		expectedValue: "dHlwZSBNb3ZpZSB7CiAgdGl0bGU6IFN0cmluZwp9Cg==",
		expectedError: nil,
	}, "invalid base 64 string": {
		flagValue:     "sdf",
		expectedValue: "",
		expectedError: errors.New("type definitions file 'sdf' does not exist"),
	}, "invalid path": {
		flagValue:     "../../test/testutils/typeDefs.graphql",
		expectedValue: "",
		expectedError: errors.New("type definitions file '../../test/testutils/typeDefs.graphql' does not exist"),
	}}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			val, err := graphql.ResolveTypeDefsFlagValue(test.flagValue)
			if test.expectedError != nil {
				assert.EqualError(t, err, test.expectedError.Error())
			} else {
				assert.Equal(t, test.expectedValue, val)
			}
		})
	}
}

func TestCreateGraphQLDataApisMissingFlags(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	helper.SetConfigValue("aura.beta-enabled", "true")

	instanceId := "2f49c2b3"

	helper.ExecuteCommand(fmt.Sprintf("data-api graphql create --instance-id %s", instanceId))

	helper.AssertErr("Error: required flag(s) \"instance-password\", \"instance-username\", \"name\", \"security-auth-provider-name\", \"security-auth-provider-type\", \"type-definitions\" not set")
}

func TestCreateGraphQLDataApisMissingFlagInstancePassword(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	helper.SetConfigValue("aura.beta-enabled", "true")

	instanceId := "2f49c2b3"
	instanceUsername := "neo4j"
	name := "my-data-api-1"
	typeDefs := "dHlwZSBBY3RvciB7CiAgbmFtZTogU3RyaW5nCiAgbW92aWVzOiBbTW92aWUhXSEgQHJlbGF0aW9uc2hpcCh0eXBlOiAiQUNURURfSU4iLCBkaXJlY3Rpb246IE9VVCkKfQoKdHlwZSBNb3ZpZSB7CiAgdGl0bGU6IFN0cmluZwogIGFjdG9yczogW0FjdG9yIV0hIEByZWxhdGlvbnNoaXAodHlwZTogIkFDVEVEX0lOIiwgZGlyZWN0aW9uOiBJTikKfQ=="
	secAuthProviderName := "provider-1"
	secAuthProviderType := "api-key"
	helper.ExecuteCommand(fmt.Sprintf("data-api graphql create --instance-id %s --instance-username %s --name %s --type-definitions %s --security-auth-provider-name %s --security-auth-provider-type %s", instanceId, instanceUsername, name, typeDefs, secAuthProviderName, secAuthProviderType))

	helper.AssertErr("Error: required flag(s) \"instance-password\" not set")
}

func TestCreateGraphQLDataApisMissingFlagInstanceUsername(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	helper.SetConfigValue("aura.beta-enabled", "true")

	instanceId := "2f49c2b3"
	instancePassword := "dlkshglsjsdfsd"
	name := "my-data-api-1"
	typeDefs := "dHlwZSBBY3RvciB7CiAgbmFtZTogU3RyaW5nCiAgbW92aWVzOiBbTW92aWUhXSEgQHJlbGF0aW9uc2hpcCh0eXBlOiAiQUNURURfSU4iLCBkaXJlY3Rpb246IE9VVCkKfQoKdHlwZSBNb3ZpZSB7CiAgdGl0bGU6IFN0cmluZwogIGFjdG9yczogW0FjdG9yIV0hIEByZWxhdGlvbnNoaXAodHlwZTogIkFDVEVEX0lOIiwgZGlyZWN0aW9uOiBJTikKfQ=="
	secAuthProviderName := "provider-1"
	secAuthProviderType := "api-key"
	helper.ExecuteCommand(fmt.Sprintf("data-api graphql create --instance-id %s --instance-password %s --name %s --type-definitions %s --security-auth-provider-name %s --security-auth-provider-type %s", instanceId, instancePassword, name, typeDefs, secAuthProviderName, secAuthProviderType))

	helper.AssertErr("Error: required flag(s) \"instance-username\" not set")
}

func TestCreateGraphQLDataApisMissingFlagName(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	helper.SetConfigValue("aura.beta-enabled", "true")

	instanceId := "2f49c2b3"
	instanceUsername := "neo4j"
	instancePassword := "dlkshglsjsdfsd"
	typeDefs := "dHlwZSBBY3RvciB7CiAgbmFtZTogU3RyaW5nCiAgbW92aWVzOiBbTW92aWUhXSEgQHJlbGF0aW9uc2hpcCh0eXBlOiAiQUNURURfSU4iLCBkaXJlY3Rpb246IE9VVCkKfQoKdHlwZSBNb3ZpZSB7CiAgdGl0bGU6IFN0cmluZwogIGFjdG9yczogW0FjdG9yIV0hIEByZWxhdGlvbnNoaXAodHlwZTogIkFDVEVEX0lOIiwgZGlyZWN0aW9uOiBJTikKfQ=="
	secAuthProviderName := "provider-1"
	secAuthProviderType := "api-key"
	helper.ExecuteCommand(fmt.Sprintf("data-api graphql create --instance-id %s --instance-username %s --instance-password %s --type-definitions %s --security-auth-provider-name %s --security-auth-provider-type %s", instanceId, instanceUsername, instancePassword, typeDefs, secAuthProviderName, secAuthProviderType))

	helper.AssertErr("Error: required flag(s) \"name\" not set")
}

func TestCreateGraphQLDataApisMissingFlagTypeDefs(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	helper.SetConfigValue("aura.beta-enabled", "true")

	instanceId := "2f49c2b3"
	instanceUsername := "neo4j"
	instancePassword := "dlkshglsjsdfsd"
	name := "my-data-api-1"
	secAuthProviderName := "provider-1"
	secAuthProviderType := "api-key"
	helper.ExecuteCommand(fmt.Sprintf("data-api graphql create --instance-id %s --instance-username %s --instance-password %s --name %s --security-auth-provider-name %s --security-auth-provider-type %s", instanceId, instanceUsername, instancePassword, name, secAuthProviderName, secAuthProviderType))

	helper.AssertErr("Error: required flag(s) \"type-definitions\" not set")
}

func TestCreateGraphQLDataApisMissingFlagAuthProviderName(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	helper.SetConfigValue("aura.beta-enabled", "true")

	instanceId := "2f49c2b3"
	instanceUsername := "neo4j"
	instancePassword := "dlkshglsjsdfsd"
	typeDefs := "dHlwZSBBY3RvciB7CiAgbmFtZTogU3RyaW5nCiAgbW92aWVzOiBbTW92aWUhXSEgQHJlbGF0aW9uc2hpcCh0eXBlOiAiQUNURURfSU4iLCBkaXJlY3Rpb246IE9VVCkKfQoKdHlwZSBNb3ZpZSB7CiAgdGl0bGU6IFN0cmluZwogIGFjdG9yczogW0FjdG9yIV0hIEByZWxhdGlvbnNoaXAodHlwZTogIkFDVEVEX0lOIiwgZGlyZWN0aW9uOiBJTikKfQ=="
	name := "my-data-api-1"
	secAuthProviderType := "api-key"
	helper.ExecuteCommand(fmt.Sprintf("data-api graphql create --instance-id %s --instance-username %s --instance-password %s --name %s --type-definitions %s --security-auth-provider-type %s", instanceId, instanceUsername, instancePassword, name, typeDefs, secAuthProviderType))

	helper.AssertErr("Error: required flag(s) \"security-auth-provider-name\" not set")
}

func TestCreateGraphQLDataApisMissingFlagAuthProviderType(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	helper.SetConfigValue("aura.beta-enabled", "true")

	instanceId := "2f49c2b3"
	instanceUsername := "neo4j"
	instancePassword := "dlkshglsjsdfsd"
	typeDefs := "dHlwZSBBY3RvciB7CiAgbmFtZTogU3RyaW5nCiAgbW92aWVzOiBbTW92aWUhXSEgQHJlbGF0aW9uc2hpcCh0eXBlOiAiQUNURURfSU4iLCBkaXJlY3Rpb246IE9VVCkKfQoKdHlwZSBNb3ZpZSB7CiAgdGl0bGU6IFN0cmluZwogIGFjdG9yczogW0FjdG9yIV0hIEByZWxhdGlvbnNoaXAodHlwZTogIkFDVEVEX0lOIiwgZGlyZWN0aW9uOiBJTikKfQ=="
	name := "my-data-api-1"
	secAuthProviderName := "provider-1"
	helper.ExecuteCommand(fmt.Sprintf("data-api graphql create --instance-id %s --instance-username %s --instance-password %s --name %s --type-definitions %s --security-auth-provider-name %s", instanceId, instanceUsername, instancePassword, name, typeDefs, secAuthProviderName))

	helper.AssertErr("Error: required flag(s) \"security-auth-provider-type\" not set")
}

func TestCreateGraphQLDataApisTypeDefsNotBase64(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	helper.SetConfigValue("aura.beta-enabled", "true")

	instanceId := "2f49c2b3"
	instanceUsername := "neo4j"
	instancePassword := "dlkshglsjsdfsd"
	typeDefs := "blabla"
	name := "my-data-api-1"
	secAuthProviderName := "provider-1"
	secAuthProviderType := "api-key"
	helper.ExecuteCommand(fmt.Sprintf("data-api graphql create --instance-id %s --instance-username %s --instance-password %s --name %s --type-definitions %s --security-auth-provider-name %s --security-auth-provider-type %s", instanceId, instanceUsername, instancePassword, name, typeDefs, secAuthProviderName, secAuthProviderType))

	helper.AssertErr("Error: provided type definitions are not valid base64")
}

// func TestCreateGraphQLDataApisTypeDefsAsFile(t *testing.T) {
// 	helper := testutils.NewAuraTestHelper(t)
// 	defer helper.Close()

// 	helper.SetConfigValue("aura.beta-enabled", "true")

// 	instanceId := "2f49c2b3"
// 	instanceUsername := "neo4j"
// 	instancePassword := "dlkshglsjsdfsd"
// 	typeDefs := "../../../test/testutils/typeDefs.graphql"
// 	name := "my-data-api-1"
// 	secAuthProviderName := "provider-1"
// 	secAuthProviderType := "api-key"
// 	helper.ExecuteCommand(fmt.Sprintf("data-api graphql create --instance-id %s --instance-username %s --instance-password %s --name %s --type-definitions %s --security-auth-provider-name %s --security-auth-provider-type %s", instanceId, instanceUsername, instancePassword, name, typeDefs, secAuthProviderName, secAuthProviderType))

// 	helper.AssertErr("Error: se64")
// }

func TestCreateGraphQLDataApisWrongAuthProviderType(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	helper.SetConfigValue("aura.beta-enabled", "true")

	instanceId := "2f49c2b3"
	instanceUsername := "neo4j"
	instancePassword := "dlkshglsjsdfsd"
	typeDefs := "blabla"
	name := "my-data-api-1"
	secAuthProviderName := "provider-1"
	secAuthProviderType := "non-existing-type"
	helper.ExecuteCommand(fmt.Sprintf("data-api graphql create --instance-id %s --instance-username %s --instance-password %s --name %s --type-definitions %s --security-auth-provider-name %s --security-auth-provider-type %s", instanceId, instanceUsername, instancePassword, name, typeDefs, secAuthProviderName, secAuthProviderType))

	helper.AssertErr("Error: invalid security auth provider type, got 'non-existing-type', expect 'api-key' or 'jwks'")
}

func TestCreateGraphQLDataApisJWKSAuthProviderMissingUrl(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	helper.SetConfigValue("aura.beta-enabled", "true")

	instanceId := "2f49c2b3"
	instanceUsername := "neo4j"
	instancePassword := "dlkshglsjsdfsd"
	typeDefs := "blabla"
	name := "my-data-api-1"
	secAuthProviderName := "provider-1"
	secAuthProviderType := "jwks"
	helper.ExecuteCommand(fmt.Sprintf("data-api graphql create --instance-id %s --instance-username %s --instance-password %s --name %s --type-definitions %s --security-auth-provider-name %s --security-auth-provider-type %s", instanceId, instanceUsername, instancePassword, name, typeDefs, secAuthProviderName, secAuthProviderType))

	helper.AssertErr("Error: required flag(s) \"security-auth-provider-url\" not set")
}

func TestCreateGraphQLDataApisOneAuthProviderWithApiKey(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	helper.SetConfigValue("aura.beta-enabled", "true")

	instanceId := "2f49c2b3"
	instanceUsername := "neo4j"
	instancePassword := "lxbckvpsdbfgsbsdfgbsdf"
	name := "my-data-api-1"
	typeDefs := "dHlwZSBBY3RvciB7CiAgbmFtZTogU3RyaW5nCiAgbW92aWVzOiBbTW92aWUhXSEgQHJlbGF0aW9uc2hpcCh0eXBlOiAiQUNURURfSU4iLCBkaXJlY3Rpb246IE9VVCkKfQoKdHlwZSBNb3ZpZSB7CiAgdGl0bGU6IFN0cmluZwogIGFjdG9yczogW0FjdG9yIV0hIEByZWxhdGlvbnNoaXAodHlwZTogIkFDVEVEX0lOIiwgZGlyZWN0aW9uOiBJTikKfQ=="
	secAuthProviderName := "provider-1"
	secAuthProviderType := "api-key"
	mockHandler := helper.NewRequestHandlerMock(fmt.Sprintf("/v1/instances/%s/data-apis/graphql", instanceId), http.StatusOK, `{
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
	}`)

	helper.ExecuteCommand(fmt.Sprintf("data-api graphql create --instance-id %s --instance-username %s --instance-password %s --name %s --type-definitions %s --security-auth-provider-name %s --security-auth-provider-type %s --feature-subgraph-enabled true", instanceId, instanceUsername, instancePassword, name, typeDefs, secAuthProviderName, secAuthProviderType))

	mockHandler.AssertCalledTimes(1)
	mockHandler.AssertCalledWithMethod(http.MethodPost)

	helper.AssertOutJson(`{
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
	}`)
}

func TestCreateGraphQLDataApisOneAuthProviderWithApiKeyOutputAsTable(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	helper.SetConfigValue("aura.beta-enabled", "true")

	instanceId := "2f49c2b3"
	instanceUsername := "neo4j"
	instancePassword := "lxbckvpsdbfgsbsdfgbsdf"
	name := "my-data-api-1"
	typeDefs := "dHlwZSBBY3RvciB7CiAgbmFtZTogU3RyaW5nCiAgbW92aWVzOiBbTW92aWUhXSEgQHJlbGF0aW9uc2hpcCh0eXBlOiAiQUNURURfSU4iLCBkaXJlY3Rpb246IE9VVCkKfQoKdHlwZSBNb3ZpZSB7CiAgdGl0bGU6IFN0cmluZwogIGFjdG9yczogW0FjdG9yIV0hIEByZWxhdGlvbnNoaXAodHlwZTogIkFDVEVEX0lOIiwgZGlyZWN0aW9uOiBJTikKfQ=="
	secAuthProviderName := "provider-1"
	secAuthProviderType := "api-key"
	mockHandler := helper.NewRequestHandlerMock(fmt.Sprintf("/v1/instances/%s/data-apis/graphql", instanceId), http.StatusOK, `{
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
	}`)

	helper.ExecuteCommand(fmt.Sprintf("data-api graphql create --output table --instance-id %s --instance-username %s --instance-password %s --name %s --type-definitions %s --security-auth-provider-name %s --security-auth-provider-type %s --feature-subgraph-enabled true", instanceId, instanceUsername, instancePassword, name, typeDefs, secAuthProviderName, secAuthProviderType))

	mockHandler.AssertCalledTimes(1)
	mockHandler.AssertCalledWithMethod(http.MethodPost)

	helper.AssertOut(`
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
	`)
}

func TestCreateGraphQLDataApisOneAuthProviderWithJwks(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	helper.SetConfigValue("aura.beta-enabled", "true")

	instanceId := "2f49c2b3"
	instanceUsername := "neo4j"
	instancePassword := "lxbckvpsdbfgsbsdfgbsdf"
	name := "my-data-api-1"
	typeDefs := "dHlwZSBBY3RvciB7CiAgbmFtZTogU3RyaW5nCiAgbW92aWVzOiBbTW92aWUhXSEgQHJlbGF0aW9uc2hpcCh0eXBlOiAiQUNURURfSU4iLCBkaXJlY3Rpb246IE9VVCkKfQoKdHlwZSBNb3ZpZSB7CiAgdGl0bGU6IFN0cmluZwogIGFjdG9yczogW0FjdG9yIV0hIEByZWxhdGlvbnNoaXAodHlwZTogIkFDVEVEX0lOIiwgZGlyZWN0aW9uOiBJTikKfQ=="
	secAuthProviderName := "provider-1"
	secAuthProviderType := "jwks"
	secAuthProviderUrl := "https://test.com/.well-known/jwks.json"
	mockHandler := helper.NewRequestHandlerMock(fmt.Sprintf("/v1/instances/%s/data-apis/graphql", instanceId), http.StatusOK, `{
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
	}`)

	helper.ExecuteCommand(fmt.Sprintf("data-api graphql create --instance-id %s --instance-username %s --instance-password %s --name %s --type-definitions %s --security-auth-provider-name %s --security-auth-provider-type %s --security-auth-provider-url %s --feature-subgraph-enabled false", instanceId, instanceUsername, instancePassword, name, typeDefs, secAuthProviderName, secAuthProviderType, secAuthProviderUrl))

	mockHandler.AssertCalledTimes(1)
	mockHandler.AssertCalledWithMethod(http.MethodPost)

	helper.AssertOutJson(`{
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
	}`)
}

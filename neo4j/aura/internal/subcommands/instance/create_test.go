package instance_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/neo4j/cli/common/clicfg"
	"github.com/neo4j/cli/common/clictx"
	"github.com/neo4j/cli/neo4j/aura"
	"github.com/neo4j/cli/neo4j/aura/internal/test/testutils"
	"github.com/neo4j/cli/test/utils/testfs"
	"github.com/stretchr/testify/assert"
)

func TestCreateFreeInstance(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	mockHandler := helper.NewRequestHandlerMock("/v1/instances", http.StatusOK, `{
			"data": {
				"id": "db1d1234",
				"connection_url": "YOUR_CONNECTION_URL",
				"username": "neo4j",
				"password": "letMeIn123!",
				"tenant_id": "YOUR_TENANT_ID",
				"cloud_provider": "gcp",
				"region": "europe-west1",
				"type": "free-db",
				"name": "Instance01"
			}
		}`)

	helper.ExecuteCommand("instance create --region europe-west1 --name Instance01 --type free-db --tenant-id YOUR_TENANT_ID --cloud-provider gcp")

	mockHandler.AssertCalledTimes(1)
	mockHandler.AssertCalledWithMethod(http.MethodPost)
	mockHandler.AssertCalledWithBody(`{"cloud_provider":"gcp","memory":"1GB","name":"Instance01","region":"europe-west1","tenant_id":"YOUR_TENANT_ID","type":"free-db","version":"5"}`)

	helper.AssertOutJson(`{
		"data": {
			"id": "db1d1234",
			"connection_url": "YOUR_CONNECTION_URL",
			"username": "neo4j",
			"password": "letMeIn123!",
			"tenant_id": "YOUR_TENANT_ID",
			"cloud_provider": "gcp",
			"region": "europe-west1",
			"type": "free-db",
			"name": "Instance01"
		}
	}`)
}

func TestCreateProfessionalInstance(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	mockHandler := helper.NewRequestHandlerMock("/v1/instances", http.StatusOK, `{
			"data": {
				"id": "db1d1234",
				"connection_url": "YOUR_CONNECTION_URL",
				"username": "neo4j",
				"password": "letMeIn123!",
				"tenant_id": "YOUR_TENANT_ID",
				"cloud_provider": "gcp",
				"region": "europe-west1",
				"type": "professional-db",
				"name": "Instance01"
			}
		}`)

	helper.ExecuteCommand("instance create --region europe-west1 --name Instance01 --type professional-db --tenant-id YOUR_TENANT_ID --cloud-provider gcp --memory 4GB")

	mockHandler.AssertCalledTimes(1)
	mockHandler.AssertCalledWithMethod(http.MethodPost)
	mockHandler.AssertCalledWithBody(`{"cloud_provider":"gcp","memory":"4GB","name":"Instance01","region":"europe-west1","tenant_id":"YOUR_TENANT_ID","type":"professional-db","version":"5"}`)

	helper.AssertOutJson(`{
		"data": {
			"id": "db1d1234",
			"connection_url": "YOUR_CONNECTION_URL",
			"username": "neo4j",
			"password": "letMeIn123!",
			"tenant_id": "YOUR_TENANT_ID",
			"cloud_provider": "gcp",
			"region": "europe-west1",
			"type": "professional-db",
			"name": "Instance01"
		}
	}
	`)
}

func TestCreateProfessionalInstanceNoMemory(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	mockHandler := helper.NewRequestHandlerMock("/v1/instances", http.StatusOK, "")

	helper.ExecuteCommand("instance create --region europe-west1 --name Instance01 --type professional-db --tenant-id YOUR_TENANT_ID --cloud-provider gcp")

	mockHandler.AssertCalledTimes(0)

	helper.AssertErr(`Error: required flag(s) "memory" not set
`)
	helper.AssertOut(`Usage:
  aura instance create [flags]

Flags:
      --cloud-provider string            The cloud provider hosting the instance.
      --customer-managed-key-id string   An optional customer managed key to be used for instance creation.
  -h, --help                             help for create
      --memory string                    The size of the instance memory in GB.
      --name string                      The name of the instance (any UTF-8 characters with no trailing or leading whitespace).
      --region string                    The region where the instance is hosted.
      --tenant-id string                 
      --type string                      The type of the instance.
      --version string                   The Neo4j version of the instance. (default "5")

Global Flags:
      --auth-url string   
      --base-url string   
      --output string

`)
}

func TestCreateProfessionalInstanceNoTenant(t *testing.T) {
	assert := assert.New(t)

	mux := http.NewServeMux()

	var authCounter = 0
	mux.HandleFunc("/oauth/token", func(res http.ResponseWriter, req *http.Request) {
		authCounter++
	})

	var postCounter = 0
	mux.HandleFunc("/v1/instances", func(res http.ResponseWriter, req *http.Request) {
		postCounter++
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	cmd := aura.NewCmd()
	b := bytes.NewBufferString("")
	cmd.SetOut(b)
	cmd.SetErr(b)
	cmd.SetArgs([]string{"instance", "create", "--auth-url", fmt.Sprintf("%s/oauth/token", server.URL), "--base-url", fmt.Sprintf("%s/v1", server.URL), "--region", "europe-west1", "--name", "Instance01", "--type", "professional-db", "--cloud-provider", "gcp", "--memory", "1GB"})

	fs, err := testfs.GetTestFs(`{
				"aura": {
			"credentials": [{
				"name": "test-cred",
				"access-token": "dsa",
				"token-expiry": 123
			}],
			"default-credential": "test-cred"
		}
	}`)
	assert.Nil(err)

	cfg, err := clicfg.NewConfig(fs)
	assert.Nil(err)

	ctx, err := clictx.NewContext(context.Background(), cfg, "test")
	assert.Nil(err)

	err = cmd.ExecuteContext(ctx)
	assert.ErrorContains(err, `required flag(s) "tenant-id" not set`)

	assert.Equal(0, authCounter)
	assert.Equal(0, postCounter)

	out, err := io.ReadAll(b)
	assert.Nil(err)

	assert.Equal(`Error: required flag(s) "tenant-id" not set
Usage:
  aura instance create [flags]

Flags:
      --cloud-provider string            The cloud provider hosting the instance.
      --customer-managed-key-id string   An optional customer managed key to be used for instance creation.
  -h, --help                             help for create
      --memory string                    The size of the instance memory in GB.
      --name string                      The name of the instance (any UTF-8 characters with no trailing or leading whitespace).
      --region string                    The region where the instance is hosted.
      --tenant-id string                 
      --type string                      The type of the instance.
      --version string                   The Neo4j version of the instance. (default "5")

Global Flags:
      --auth-url string   
      --base-url string   
      --output string

`, string(out))
}

func TestCreateInstanceError(t *testing.T) {
	testCases := []struct {
		statusCode    int
		expectedError string
		returnBody    string
	}{
		{
			statusCode:    http.StatusBadRequest,
			expectedError: "Error: [You must provide billing details in the Aura Console before creating an instance]",
			returnBody: `{
				"errors": [
					{
					"message": "You must provide billing details in the Aura Console before creating an instance",
					"reason": "missing-billing-details"
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

			mockHandler := helper.NewRequestHandlerMock("/v1/instances", testCase.statusCode, testCase.returnBody)

			helper.ExecuteCommand("instance create --region europe-west1 --name Instance01 --type professional-db --tenant-id YOUR_TENANT_ID --cloud-provider gcp --memory 4GB")

			mockHandler.AssertCalledTimes(1)
			mockHandler.AssertCalledWithMethod(http.MethodPost)

			helper.AssertOut("")
			helper.AssertErr(testCase.expectedError)
		})
	}
}

func TestInstanceWithCmkId(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	mockHandler := helper.NewRequestHandlerMock("/v1/instances", http.StatusOK, `{
			"data": {
				"id": "db1d1234",
				"connection_url": "YOUR_CONNECTION_URL",
				"username": "neo4j",
				"password": "letMeIn123!",
				"tenant_id": "YOUR_TENANT_ID",
				"cloud_provider": "gcp",
				"region": "europe-west1",
				"type": "enterprise-db",
				"name": "Instance01"
			}
		}`)

	helper.ExecuteCommand("instance create --region europe-west1 --name Instance01 --type enterprise-db --tenant-id YOUR_TENANT_ID --cloud-provider gcp --memory 16GB --customer-managed-key-id UUID_OF_YOUR_KEY")

	mockHandler.AssertCalledTimes(1)
	mockHandler.AssertCalledWithMethod(http.MethodPost)
	mockHandler.AssertCalledWithBody(`{"cloud_provider":"gcp","memory":"16GB","name":"Instance01","region":"europe-west1","tenant_id":"YOUR_TENANT_ID","type":"enterprise-db","version":"5","customer_managed_key_id":"UUID_OF_YOUR_KEY"}`)

	helper.AssertOutJson(`{
		"data": {
			"id": "db1d1234",
			"connection_url": "YOUR_CONNECTION_URL",
			"username": "neo4j",
			"password": "letMeIn123!",
			"tenant_id": "YOUR_TENANT_ID",
			"cloud_provider": "gcp",
			"region": "europe-west1",
			"type": "enterprise-db",
			"name": "Instance01"
		}
	}`)
}

func TestCreateFreeInstanceWithConfigTenantId(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	helper.SetConfig(`{
		"aura": {
		"default-tenant": "YOUR_TENANT_ID",
		"credentials": [{
			"name": "test-cred",
			"access-token": "dsa",
			"token-expiry": 123
		}],
		"default-credential": "test-cred"
		}
	}`)

	mockHandler := helper.NewRequestHandlerMock("/v1/instances", http.StatusOK, `{
			"data": {
				"id": "db1d1234",
				"connection_url": "YOUR_CONNECTION_URL",
				"username": "neo4j",
				"password": "letMeIn123!",
				"tenant_id": "YOUR_TENANT_ID",
				"cloud_provider": "gcp",
				"region": "europe-west1",
				"type": "free-db",
				"name": "Instance01"
			}
		}`)

	helper.ExecuteCommand("instance create --region europe-west1 --name Instance01 --type free-db --cloud-provider gcp")

	mockHandler.AssertCalledTimes(1)
	mockHandler.AssertCalledWithMethod(http.MethodPost)
	mockHandler.AssertCalledWithBody(`{"cloud_provider":"gcp","memory":"1GB","name":"Instance01","region":"europe-west1","tenant_id":"YOUR_TENANT_ID","type":"free-db","version":"5"}`)

	helper.AssertOutJson(`{
		"data": {
			"id": "db1d1234",
			"connection_url": "YOUR_CONNECTION_URL",
			"username": "neo4j",
			"password": "letMeIn123!",
			"tenant_id": "YOUR_TENANT_ID",
			"cloud_provider": "gcp",
			"region": "europe-west1",
			"type": "free-db",
			"name": "Instance01"
		}
	}`)
}

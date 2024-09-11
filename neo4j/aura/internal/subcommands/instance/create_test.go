package instance_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/neo4j/cli/neo4j/aura/internal/test/testutils"
)

func TestCreateFreeInstance(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	mockHandler := helper.NewRequestHandlerMock("/v1/instances", http.StatusAccepted, `{
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

	mockHandler := helper.NewRequestHandlerMock("/v1/instances", http.StatusAccepted, `{
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
      --await                            Waits until created instance is ready.
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
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	mockHandler := helper.NewRequestHandlerMock("/v1/instances", http.StatusOK, "")

	helper.ExecuteCommand("instance create --region europe-west1 --name Instance01 --type professional-db --memory 1GB --cloud-provider gcp")

	mockHandler.AssertCalledTimes(0)

	helper.AssertErr(`Error: required flag(s) "tenant-id" not set
`)
	helper.AssertOut(`Usage:
  aura instance create [flags]

Flags:
      --await                            Waits until created instance is ready.
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

	mockHandler := helper.NewRequestHandlerMock("/v1/instances", http.StatusAccepted, `{
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

	helper.SetConfigValue("aura.default-tenant", "YOUR_TENANT_ID")

	mockHandler := helper.NewRequestHandlerMock("/v1/instances", http.StatusAccepted, `{
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

func TestCreateFreeInstanceWithAwait(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	createMock := helper.NewRequestHandlerMock("POST /v1/instances", http.StatusAccepted, `{
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

	getMock := helper.NewRequestHandlerMock("GET /v1/instances/db1d1234", http.StatusOK, `{
			"data": {
				"id": "db1d1234",
				"status": "creating"
			}
		}`).AddResponse(http.StatusOK, `{
			"data": {
				"id": "db1d1234",
				"status": "ready"
			}
		}`)

	helper.ExecuteCommand("instance create --region europe-west1 --name Instance01 --type free-db --tenant-id YOUR_TENANT_ID --cloud-provider gcp --await")

	createMock.AssertCalledTimes(1)
	createMock.AssertCalledWithMethod(http.MethodPost)
	createMock.AssertCalledWithBody(`{"cloud_provider":"gcp","memory":"1GB","name":"Instance01","region":"europe-west1","tenant_id":"YOUR_TENANT_ID","type":"free-db","version":"5"}`)

	getMock.AssertCalledTimes(2)
	getMock.AssertCalledWithMethod(http.MethodGet)

	helper.AssertOut(`
{
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
}
Waiting for instance to be ready...
Instance Status: ready
	`)
}

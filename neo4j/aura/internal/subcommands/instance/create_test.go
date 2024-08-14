package instance_test

import (
	"net/http"
	"testing"

	"github.com/neo4j/cli/neo4j/aura/internal/test/testutils"
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

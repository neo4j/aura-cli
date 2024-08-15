package customermanagedkey_test

import (
	"net/http"
	"testing"

	"github.com/neo4j/cli/neo4j/aura/internal/test/testutils"
)

func TestCreateCustomerManagedKeys(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	mockHandler := helper.NewRequestHandlerMock("/v1/customer-managed-keys", http.StatusAccepted, `{
		"data": {
		  "id": "8c764aed-8eb3-4a1c-92f6-e4ef0c7a6ed9",
		  "name": "Instance01",
		  "created": "2024-01-31T14:06:57Z",
		  "cloud_provider": "aws",
		  "key_id": "arn:aws:kms:us-east-1:123456789:key/11111-a222-1212-x789-1212f1212f",
		  "region": "us-east-1",
		  "type": "enterprise-db",
		  "tenant_id": "dontpanic",
		  "status": "pending"
		}
	  }`)

	helper.ExecuteCommand(`customer-managed-key create --region us-west-2 --name "Production Key" --type enterprise-db --tenant-id dontpanic --cloud-provider aws --key-id arn:aws:kms:us-west-2:111122223333:key/1234abcd-12ab-34cd-56ef-1234567890ab`)

	mockHandler.AssertCalledTimes(1)
	mockHandler.AssertCalledWithMethod(http.MethodPost)
	mockHandler.AssertCalledWithBody(`{
		"key_id": "arn:aws:kms:us-west-2:111122223333:key/1234abcd-12ab-34cd-56ef-1234567890ab",
		"name": "Production Key",
		"cloud_provider": "aws",
		"instance_type": "enterprise-db",
		"region": "us-west-2",
		"tenant_id": "dontpanic"
	}`)

	helper.AssertErr("")
	helper.AssertOutJson(`{
		"data": {
		  "id": "8c764aed-8eb3-4a1c-92f6-e4ef0c7a6ed9",
		  "name": "Instance01",
		  "created": "2024-01-31T14:06:57Z",
		  "cloud_provider": "aws",
		  "key_id": "arn:aws:kms:us-east-1:123456789:key/11111-a222-1212-x789-1212f1212f",
		  "region": "us-east-1",
		  "type": "enterprise-db",
		  "tenant_id": "dontpanic",
		  "status": "pending"
		}
	  }
		`)
}

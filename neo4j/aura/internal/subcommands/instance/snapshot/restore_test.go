package snapshot_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/neo4j/cli/neo4j/aura/internal/test/testutils"
)

func TestRestoreSnapshot(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()
	instanceId := "2f49c2b3"
	snapshotId := "afdb4e9d-6ba6-4d45-b951-f82843dcbca6"
	mockHandler := helper.NewRequestHandlerMock(fmt.Sprintf("/v1/instances/%s/snapshots/%s/restore", instanceId, snapshotId), http.StatusAccepted, `{
		"data": {
		  "id": "2f49c2b3",
		  "name": "Production",
		  "status": "restoring",
		  "connection_url": "YOUR_CONNECTION_URL",
		  "tenant_id": "YOUR_TENANT_ID",
		  "cloud_provider": "gcp",
		  "memory": "8GB",
		  "region": "europe-west1",
		  "type": "enterprise-db"
		}
	  }`)

	helper.ExecuteCommand(fmt.Sprintf("instance snapshot restore --instance-id %s %s", instanceId, snapshotId))

	mockHandler.AssertCalledTimes(1)
	mockHandler.AssertCalledWithMethod(http.MethodPost)

	helper.AssertOutJson(`{
		"data": {
			"id": "2f49c2b3",
			"name": "Production",
			"status": "restoring",
			"connection_url": "YOUR_CONNECTION_URL",
			"tenant_id": "YOUR_TENANT_ID",
			"cloud_provider": "gcp",
			"memory": "8GB",
			"region": "europe-west1",
			"type": "enterprise-db"
		  }
	}`)
}

package snapshot_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/neo4j/cli/neo4j-cli/aura/internal/test/testutils"
)

func TestListSnapshot(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()
	instanceId := "2f49c2b3"
	mockHandler := helper.NewRequestHandlerMock(fmt.Sprintf("/v1/instances/%s/snapshots", instanceId), http.StatusOK, `{
		"data": [
			{
				"exportable": true,
				"instance_id": "7261d20a",
				"profile": "AdHoc",
				"snapshot_id": "afdb4e9d-6ba6-4d45-b951-f82843dcbca6",
				"status": "Completed",
				"timestamp": "2024-09-12T13:51:45Z"
			}
		]	
		}`)

	helper.ExecuteCommand(fmt.Sprintf("instance snapshot list --instance-id %s", instanceId))

	mockHandler.AssertCalledTimes(1)
	mockHandler.AssertCalledWithMethod(http.MethodGet)

	helper.AssertOutJson(`{
		"data": [
			{
				"exportable": true,
				"instance_id": "7261d20a",
				"profile": "AdHoc",
				"snapshot_id": "afdb4e9d-6ba6-4d45-b951-f82843dcbca6",
				"status": "Completed",
				"timestamp": "2024-09-12T13:51:45Z"
			}
		]
	}
	`)
}

func TestListSnapshotWithDate(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()
	instanceId := "2f49c2b3"
	mockHandler := helper.NewRequestHandlerMock(fmt.Sprintf("/v1/instances/%s/snapshots", instanceId), http.StatusOK, `{
		"data": [
			{
				"exportable": true,
				"instance_id": "7261d20a",
				"profile": "AdHoc",
				"snapshot_id": "afdb4e9d-6ba6-4d45-b951-f82843dcbca6",
				"status": "Completed",
				"timestamp": "2024-09-12T13:51:45Z"
			}
		]	
		}`)

	helper.ExecuteCommand(fmt.Sprintf("instance snapshot list --instance-id %s --date 2024-02-13", instanceId))

	mockHandler.AssertCalledTimes(1)
	mockHandler.AssertCalledWithMethod(http.MethodGet)
	mockHandler.AssertCalledWithQueryParam("date", "2024-02-13")

	helper.AssertOutJson(`{
		"data": [
			{
				"exportable": true,
				"instance_id": "7261d20a",
				"profile": "AdHoc",
				"snapshot_id": "afdb4e9d-6ba6-4d45-b951-f82843dcbca6",
				"status": "Completed",
				"timestamp": "2024-09-12T13:51:45Z"
			}
		]
	}
	`)
}

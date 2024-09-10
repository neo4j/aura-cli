package snapshot_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/neo4j/cli/neo4j/aura/internal/test/testutils"
)

func TestCreateSnapshot(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()
	instanceId := "2f49c2b3"
	mockHandler := helper.NewRequestHandlerMock(fmt.Sprintf("/v1/instances/%s/snapshots", instanceId), http.StatusAccepted, `{
		"data": {
		  "snapshot_id": "snap123"
		}
	  }`)

	helper.ExecuteCommand(fmt.Sprintf("instance snapshot create --instance-id %s", instanceId))

	mockHandler.AssertCalledTimes(1)
	mockHandler.AssertCalledWithMethod(http.MethodPost)

	helper.AssertOutJson(`{
		"data": {
		  "snapshot_id": "snap123"
		}
	  }`)
}

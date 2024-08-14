package customermanagedkey_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/neo4j/cli/neo4j/aura/internal/test/testutils"
)

func TestDeleteCustomerManagedKey(t *testing.T) {
	for _, command := range []string{"customer-managed-key", "cmk"} {
		helper := testutils.NewAuraTestHelper(t)
		defer helper.Close()

		cmkId := "8c764aed-8eb3-4a1c-92f6-e4ef0c7a6ed9"

		mockHandler := helper.NewRequestHandlerMock(fmt.Sprintf("/v1/customer-managed-keys/%s", cmkId), http.StatusNoContent, "")

		helper.ExecuteCommand(fmt.Sprintf("%s delete %s", command, cmkId))

		mockHandler.AssertCalledTimes(1)
		mockHandler.AssertCalledWithMethod(http.MethodDelete)

		helper.AssertOut("Operation Successful\n")
	}
}

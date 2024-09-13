package credential_test

import (
	"testing"

	"github.com/neo4j/cli/neo4j/aura/internal/test/testutils"
)

func TestRemoveCredential(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	helper.SetConfigValue("aura.credentials", []map[string]string{{"name": "test", "client-id": "testclientid", "client-secret": "testclientsecret"}})

	helper.ExecuteCommand("credential remove test")

	helper.AssertConfigValue("aura.credentials", "[]")
}

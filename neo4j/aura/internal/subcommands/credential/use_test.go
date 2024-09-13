package credential_test

import (
	"testing"

	"github.com/neo4j/cli/neo4j/aura/internal/test/testutils"
)

func TestUseCredential(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	helper.SetConfigValue("aura.credentials", []map[string]string{{"name": "test", "client-id": "testclientid", "client-secret": "testclientsecret"}})

	helper.ExecuteCommand("credential use test")

	helper.AssertConfigValue("aura.default-credential", "test")
}

func TestUseCredentialIfDoesNotExist(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	helper.ExecuteCommand("credential use test")

	helper.AssertErr("Error: could not find credential with name test")
}

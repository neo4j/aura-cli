package credential_test

import (
	"testing"

	"github.com/neo4j/cli/neo4j-cli/aura/internal/test/testutils"
)

func TestUseCredential(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	helper.SetCredentialsValue("aura.credentials", []map[string]string{{"name": "test", "client-id": "testclientid", "client-secret": "testclientsecret"}})

	helper.ExecuteCommand("credential use test")

	helper.AssertCredentialsValue("aura.default-credential", "test")
}

func TestUseCredentialIfDoesNotExist(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	helper.ExecuteCommand("credential use test")

	helper.AssertErr("Error: could not find credential with name test")
}

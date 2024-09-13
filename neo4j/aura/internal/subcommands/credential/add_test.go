package credential_test

import (
	"testing"

	"github.com/neo4j/cli/neo4j/aura/internal/test/testutils"
)

func TestAddCredential(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	helper.SetConfigValue("aura.credentials", []map[string]string{})

	helper.ExecuteCommand("credential add --name test --client-id testclientid --client-secret testclientsecret")

	helper.AssertConfigValue("aura.credentials", `[{"name":"test","client-id":"testclientid","client-secret":"testclientsecret","access-token":"","token-expiry":0}]`)
}

func TestAddCredentialIfAlreadyExists(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	helper.SetConfigValue("aura.credentials", []map[string]string{{"name": "test", "client-id": "testclientid", "client-secret": "testclientsecret"}})

	helper.ExecuteCommand("credential add --name test --client-id testclientid --client-secret testclientsecret")

	helper.AssertErr("Error: already have credential with name test")
}

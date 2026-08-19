// Copyright (c) "Neo4j"
// Neo4j Sweden AB [http://neo4j.com]

package credential_test

import (
	"testing"

	"github.com/neo4j/cli/neo4j-cli/aura/internal/test/testutils"
)

func TestAddFirstCredential(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	helper.SetCredentialsValue("aura.credentials", []map[string]string{})

	helper.ExecuteCommand("credential add --name test --client-id testclientid --client-secret testclientsecret")

	// Verify the values are present (order-independent check)
	helper.AssertCredentialsValue("aura.credentials.0.name", "test")
	helper.AssertCredentialsValue("aura.credentials.0.client-id", "testclientid")
	helper.AssertCredentialsValue("aura.credentials.0.client-secret", "testclientsecret")
	helper.AssertCredentialsValue("aura.credentials.0.access-token", "")
	helper.AssertCredentialsValue("aura.credentials.0.token-expiry", "0")
	helper.AssertCredentialsValue("aura.default-credential", "test")
}

func TestAddCredentialIfAlreadyExists(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	helper.SetCredentialsValue("aura.credentials", []map[string]string{{"name": "test", "client-id": "testclientid", "client-secret": "testclientsecret"}})

	helper.ExecuteCommand("credential add --name test --client-id testclientid --client-secret testclientsecret")

	helper.AssertErr("Error: already have credential with name test")
}
func TestAddAditionalCredentials(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	helper.SetCredentialsValue("aura.credentials", []map[string]string{{"name": "test", "client-id": "testclientid", "client-secret": "testclientsecret"}})
	helper.SetCredentialsValue("aura.default-credential", "test")

	helper.ExecuteCommand("credential add --name test-new --client-id testclientid2 --client-secret testclientsecret2")

	// Check that both credentials exist with correct values
	helper.AssertCredentialsValue("aura.credentials.0.name", "test")
	helper.AssertCredentialsValue("aura.credentials.0.client-id", "testclientid")
	helper.AssertCredentialsValue("aura.credentials.0.client-secret", "testclientsecret")
	helper.AssertCredentialsValue("aura.credentials.1.name", "test-new")
	helper.AssertCredentialsValue("aura.credentials.1.client-id", "testclientid2")
	helper.AssertCredentialsValue("aura.credentials.1.client-secret", "testclientsecret2")
	helper.AssertCredentialsValue("aura.default-credential", "test")
}

// Copyright (c) "Neo4j"
// Neo4j Sweden AB [http://neo4j.com]

package credential_test

import (
	"encoding/json"
	"testing"

	"github.com/neo4j/cli/neo4j-cli/aura/internal/test/testutils"
	"github.com/stretchr/testify/assert"
)

func TestCredentialListMasksClientSecret(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	secretValue := "S3cr3t-Value-12345"
	helper.SetCredentialsValue("aura.credentials", []map[string]string{{
		"name":          "demo",
		"client-id":     "abc",
		"client-secret": secretValue,
	}})

	helper.ExecuteCommand("credential list")

	out := helper.PrintOut()

	assert.NotContains(t, out, secretValue,
		"credential list should not print the real client secret")

	assert.Contains(t, out, "****",
		"credential list should contain the mask marker for secrets")

	var credentials []map[string]interface{}
	err := json.Unmarshal([]byte(out), &credentials)
	assert.NoError(t, err, "output should be valid JSON")

	assert.Len(t, credentials, 1, "should have one credential")
	assert.Equal(t, "demo", credentials[0]["name"], "credential name should match")
	assert.Equal(t, "****", credentials[0]["client-secret"], "client-secret should be masked in JSON")
}

func TestCredentialListFileStillContainsRealSecret(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	secretValue := "RealSecretValue-98765"
	helper.SetCredentialsValue("aura.credentials", []map[string]string{{
		"name":          "prod",
		"client-id":     "xyz",
		"client-secret": secretValue,
	}})

	helper.ExecuteCommand("credential list")

	helper.AssertCredentialsValue("aura.credentials.0.client-secret", secretValue)
}

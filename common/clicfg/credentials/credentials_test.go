// Copyright (c) "Neo4j"
// Neo4j Sweden AB [http://neo4j.com]

package credentials

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/neo4j/cli/common/redact"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOnDiskJSONStructure(t *testing.T) {
	fs := afero.NewMemMapFs()
	c := NewCredentials(fs, "/test")

	err := c.Aura.Add("test-cred", "test-client-id", "test-secret")
	assert.NoError(t, err, "failed to add credential")

	err = c.Aura.SetDefault("test-cred")
	assert.NoError(t, err, "failed to set default")

	data, _ := afero.ReadFile(fs, "/test/neo4j/cli/credentials.json")
	var parsed map[string]interface{}
	err = json.Unmarshal(data, &parsed)
	assert.NoError(t, err, "failed to parse saved JSON")

	auraData, ok := parsed["aura"].(map[string]interface{})
	require.True(t, ok, "aura field is not a map")

	credentialsData, ok := auraData["credentials"].([]interface{})
	require.True(t, ok, "credentials field is not an array")

	require.Equal(t, 1, len(credentialsData), "expected 1 credential")

	credentialMap, ok := credentialsData[0].(map[string]interface{})
	require.True(t, ok, "credential is not a map")

	credType := reflect.TypeOf((*AuraCredential)(nil)).Elem()
	for i := 0; i < credType.NumField(); i++ {
		field := credType.Field(i)
		if tag, ok := field.Tag.Lookup("json"); ok {
			_, present := credentialMap[tag]
			assert.True(t, present, "field '%s' (JSON tag '%s') missing from on-disk JSON", field.Name, tag)
		}
	}

	secret, ok := credentialMap["client-secret"].(string)
	require.True(t, ok, "client-secret field is not a string")
	assert.Equal(t, "test-secret", secret, "expected secret in saved file")

	for field, value := range credentialMap {
		if str, ok := value.(string); ok {
			assert.NotEqual(t, "****", str, "field '%s' was persisted as the mask instead of its real value", field)
		}
	}
}

// A redact.Secret field on AuraCredential without a shadow field on
// auraCredentialOnDisk would silently persist as the mask instead of its real value.
func TestSecretFieldsHaveOnDiskShadow(t *testing.T) {
	secretType := reflect.TypeOf(redact.Secret{})
	credType := reflect.TypeOf(AuraCredential{})

	onDiskType := reflect.TypeOf(auraCredentialOnDisk{})
	shadowKindByTag := map[string]reflect.Kind{}
	for i := 0; i < onDiskType.NumField(); i++ {
		field := onDiskType.Field(i)
		if field.Anonymous {
			continue
		}
		if tag, ok := field.Tag.Lookup("json"); ok {
			shadowKindByTag[tag] = field.Type.Kind()
		}
	}

	for i := 0; i < credType.NumField(); i++ {
		field := credType.Field(i)
		if field.Type != secretType {
			continue
		}

		tag, ok := field.Tag.Lookup("json")
		require.True(t, ok, "secret field '%s' has no json tag", field.Name)

		kind, shadowed := shadowKindByTag[tag]
		if assert.True(t, shadowed, "secret field '%s' (json tag '%s') has no shadow field on auraCredentialOnDisk - it will be persisted as the mask instead of the real value", field.Name, tag) {
			assert.Equal(t, reflect.String, kind, "shadow field for '%s' should be a plain string", field.Name)
		}
	}
}

func TestRoundTripSaveAndLoad(t *testing.T) {
	fs := afero.NewMemMapFs()

	c1 := NewCredentials(fs, "/test")
	c1.Aura.Add("cred1", "id1", "secret1")
	c1.Aura.SetDefault("cred1")

	c2 := NewCredentials(fs, "/test")

	assert.Equal(t, 1, len(c2.Aura.Credentials), "expected 1 credential after load")

	cred := c2.Aura.Credentials[0]
	assert.Equal(t, "cred1", cred.Name, "expected name to match")
	assert.Equal(t, "id1", cred.ClientId, "expected client-id to match")
	assert.Equal(t, "secret1", cred.ClientSecret.Reveal(), "expected secret to match")
	assert.Equal(t, "cred1", c2.Aura.DefaultCredential, "expected default credential to match")

	assert.Equal(t, "****", cred.ClientSecret.String(), "expected masked secret")
}

func TestLoadPartialCredentialObjectDoesNotPanic(t *testing.T) {
	fs := afero.NewMemMapFs()
	afero.WriteFile(fs, "/test/neo4j/cli/credentials.json", []byte(`{"aura":{"credentials":[{"client-secret":"s"}]}}`), 0o644)

	c := NewCredentials(fs, "/test")

	assert.Equal(t, 1, len(c.Aura.Credentials), "expected 1 credential after load")
	assert.Equal(t, "s", c.Aura.Credentials[0].ClientSecret.Reveal(), "expected secret to match")
}

func TestJSONMarshalingPreservesSecretMasking(t *testing.T) {
	fs := afero.NewMemMapFs()
	c := NewCredentials(fs, "/test")
	c.Aura.Add("test-cred", "test-id", "test-secret")

	data, err := json.Marshal(c.Aura.Credentials)
	assert.NoError(t, err, "failed to marshal credentials")

	var unmarshaled []map[string]interface{}
	err = json.Unmarshal(data, &unmarshaled)
	assert.NoError(t, err, "failed to unmarshal credentials")

	assert.Equal(t, 1, len(unmarshaled), "expected 1 credential")

	secret := unmarshaled[0]["client-secret"]
	assert.Equal(t, "****", secret, "expected masked secret in JSON output")
}

func TestConcurrentSaveDoesNotLoseAnEarlierWrite(t *testing.T) {
	fs := afero.NewMemMapFs()
	const configPrefix = "/config"

	// Two CLI processes starting at (roughly) the same time both load the
	// same, initially-empty store.
	processA := NewCredentials(fs, configPrefix)
	processB := NewCredentials(fs, configPrefix)

	// Process A finishes first: `credential add --name first ...`.
	assert.NoError(t, processA.Aura.Add("first", "client-a", "secret-a"))

	// Process B finishes second — e.g. a routine token refresh triggering a
	// save as a side effect, or a second `credential add`. It never re-read
	// the file, so its in-memory view still thinks the store is empty.
	assert.NoError(t, processB.Aura.Add("second", "client-b", "secret-b"))

	// A fresh read of the file is what the next command would see.
	onDisk := NewCredentials(fs, configPrefix)
	names := make([]string, 0, len(onDisk.Aura.List()))
	for _, c := range onDisk.Aura.List() {
		names = append(names, c.Name)
	}

	assert.ElementsMatch(t, []string{"first", "second"}, names,
		"a save from one process must not silently erase a credential a concurrently-running process already saved")
}

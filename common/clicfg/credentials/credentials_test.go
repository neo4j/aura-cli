// Copyright (c) "Neo4j"
// Neo4j Sweden AB [http://neo4j.com]

package credentials

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/spf13/afero"
)

func TestOnDiskJSONStructure(t *testing.T) {
	fs := afero.NewMemMapFs()
	c := NewCredentials(fs, "/test")

	err := c.Aura.Add("test-cred", "test-client-id", "test-secret")
	if err != nil {
		t.Fatalf("failed to add credential: %v", err)
	}

	err = c.Aura.SetDefault("test-cred")
	if err != nil {
		t.Fatalf("failed to set default: %v", err)
	}

	data, _ := afero.ReadFile(fs, "/test/neo4j/cli/credentials.json")
	var parsed map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("failed to parse saved JSON: %v", err)
	}

	auraData, ok := parsed["aura"].(map[string]interface{})
	if !ok {
		t.Fatal("aura field is not a map")
	}

	credentialsData, ok := auraData["credentials"].([]interface{})
	if !ok {
		t.Fatal("credentials field is not an array")
	}

	if len(credentialsData) != 1 {
		t.Fatalf("expected 1 credential, got %d", len(credentialsData))
	}

	credentialMap, ok := credentialsData[0].(map[string]interface{})
	if !ok {
		t.Fatal("credential is not a map")
	}

	credType := reflect.TypeOf((*AuraCredential)(nil)).Elem()
	for i := 0; i < credType.NumField(); i++ {
		field := credType.Field(i)
		if tag, ok := field.Tag.Lookup("json"); ok {
			if _, present := credentialMap[tag]; !present {
				t.Errorf("field '%s' (JSON tag '%s') missing from on-disk JSON", field.Name, tag)
			}
		}
	}

	if secret, ok := credentialMap["client-secret"].(string); ok {
		if secret != "test-secret" {
			t.Errorf("expected secret 'test-secret' in saved file, got '%s'", secret)
		}
	} else {
		t.Fatal("client-secret field is not a string")
	}
}

func TestRoundTripSaveAndLoad(t *testing.T) {
	fs := afero.NewMemMapFs()

	c1 := NewCredentials(fs, "/test")
	c1.Aura.Add("cred1", "id1", "secret1")
	c1.Aura.SetDefault("cred1")

	c2 := NewCredentials(fs, "/test")

	if len(c2.Aura.Credentials) != 1 {
		t.Fatalf("expected 1 credential after load, got %d", len(c2.Aura.Credentials))
	}

	cred := c2.Aura.Credentials[0]
	if cred.Name != "cred1" {
		t.Errorf("expected name 'cred1', got '%s'", cred.Name)
	}
	if cred.ClientId != "id1" {
		t.Errorf("expected client-id 'id1', got '%s'", cred.ClientId)
	}
	if cred.ClientSecret.Reveal() != "secret1" {
		t.Errorf("expected secret 'secret1', got '%s'", cred.ClientSecret.Reveal())
	}
	if c2.Aura.DefaultCredential != "cred1" {
		t.Errorf("expected default 'cred1', got '%s'", c2.Aura.DefaultCredential)
	}

	if cred.ClientSecret.String() != "****" {
		t.Errorf("expected masked secret '****', got '%s'", cred.ClientSecret.String())
	}
}

func TestLoadPartialCredentialObjectDoesNotPanic(t *testing.T) {
	fs := afero.NewMemMapFs()
	afero.WriteFile(fs, "/test/neo4j/cli/credentials.json", []byte(`{"aura":{"credentials":[{"client-secret":"s"}]}}`), 0o644)

	c := NewCredentials(fs, "/test")

	if len(c.Aura.Credentials) != 1 {
		t.Fatalf("expected 1 credential after load, got %d", len(c.Aura.Credentials))
	}
	if c.Aura.Credentials[0].ClientSecret.Reveal() != "s" {
		t.Errorf("expected secret 's', got '%s'", c.Aura.Credentials[0].ClientSecret.Reveal())
	}
}

func TestJSONMarshalingPreservesSecretMasking(t *testing.T) {
	fs := afero.NewMemMapFs()
	c := NewCredentials(fs, "/test")
	c.Aura.Add("test-cred", "test-id", "test-secret")

	data, err := json.Marshal(c.Aura.Credentials)
	if err != nil {
		t.Fatalf("failed to marshal credentials: %v", err)
	}

	var unmarshaled []map[string]interface{}
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("failed to unmarshal credentials: %v", err)
	}

	if len(unmarshaled) != 1 {
		t.Fatalf("expected 1 credential, got %d", len(unmarshaled))
	}

	secret := unmarshaled[0]["client-secret"]
	if secret != "****" {
		t.Errorf("expected masked secret '****' in JSON output, got '%v'", secret)
	}
}

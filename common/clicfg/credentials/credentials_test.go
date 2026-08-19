// Copyright (c) "Neo4j"
// Neo4j Sweden AB [http://neo4j.com]

package credentials

import (
	"encoding/json"
	"testing"

	"github.com/spf13/afero"
)

// TestOnDiskJSONStructure verifies that the on-disk JSON uses a typed struct representation,
// not a hand-built map. This ensures future fields added to AuraCredential are preserved.
func TestOnDiskJSONStructure(t *testing.T) {
	fs := afero.NewMemMapFs()
	c := NewCredentials(fs, "/test")

	// Add a credential
	err := c.Aura.Add("test-cred", "test-client-id", "test-secret")
	if err != nil {
		t.Fatalf("failed to add credential: %v", err)
	}

	// Set the default
	err = c.Aura.SetDefault("test-cred")
	if err != nil {
		t.Fatalf("failed to set default: %v", err)
	}

	// Read the saved JSON
	data, _ := afero.ReadFile(fs, "/test/neo4j/cli/credentials.json")
	var parsed map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("failed to parse saved JSON: %v", err)
	}

	// Verify the structure uses a typed representation
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

	// Verify all expected fields are present
	expectedFields := map[string]bool{
		"name":            true,
		"client-id":       true,
		"client-secret":   true,
		"access-token":    true,
		"token-expiry":    true,
	}

	for field := range credentialMap {
		if !expectedFields[field] {
			t.Errorf("unexpected field in saved credential: %s", field)
		}
	}

	for field := range expectedFields {
		if _, ok := credentialMap[field]; !ok {
			t.Errorf("expected field not found in saved credential: %s", field)
		}
	}

	// Verify the secret is revealed (unmasked) in the saved file
	if secret, ok := credentialMap["client-secret"].(string); ok {
		if secret != "test-secret" {
			t.Errorf("expected secret 'test-secret' in saved file, got '%s'", secret)
		}
	} else {
		t.Fatal("client-secret field is not a string")
	}
}

// TestRoundTripSaveAndLoad verifies that credentials survive a save/load cycle
// with the secret properly masked in display but unmasked in storage.
func TestRoundTripSaveAndLoad(t *testing.T) {
	fs := afero.NewMemMapFs()

	// Create and save credentials
	c1 := NewCredentials(fs, "/test")
	c1.Aura.Add("cred1", "id1", "secret1")
	c1.Aura.SetDefault("cred1")

	// Load credentials from the same file
	c2 := NewCredentials(fs, "/test")

	// Verify the loaded credential matches
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

	// Verify the secret is masked when used as redact.Secret
	if cred.ClientSecret.String() != "****" {
		t.Errorf("expected masked secret '****', got '%s'", cred.ClientSecret.String())
	}
}

// TestJSONMarshalingPreservesSecretMasking verifies that when AuraCredentials
// are marshaled to JSON (e.g., for display), the secret is masked.
func TestJSONMarshalingPreservesSecretMasking(t *testing.T) {
	fs := afero.NewMemMapFs()
	c := NewCredentials(fs, "/test")
	c.Aura.Add("test-cred", "test-id", "test-secret")

	// Marshal the credentials for display
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

	// Verify the secret is masked in the JSON output
	secret := unmarshaled[0]["client-secret"]
	if secret != "****" {
		t.Errorf("expected masked secret '****' in JSON output, got '%v'", secret)
	}
}

// TestConversionMethods verifies the conversion methods work correctly.
func TestConversionMethods(t *testing.T) {
	// Create an on-disk representation
	onDisk := auraCredentialsOnDisk{
		DefaultCredential: "default",
		Credentials: []auraCredentialOnDisk{
			{
				Name:         "cred1",
				ClientId:     "id1",
				ClientSecret: "secret1",
				AccessToken:  "token1",
				TokenExpiry:  1234567890,
			},
		},
	}

	// Convert to in-memory format
	auraCredentials := onDisk.toAuraCredentials(func() {})

	if auraCredentials.DefaultCredential != "default" {
		t.Errorf("expected default 'default', got '%s'", auraCredentials.DefaultCredential)
	}
	if len(auraCredentials.Credentials) != 1 {
		t.Fatalf("expected 1 credential, got %d", len(auraCredentials.Credentials))
	}

	cred := auraCredentials.Credentials[0]
	if cred.Name != "cred1" {
		t.Errorf("expected name 'cred1', got '%s'", cred.Name)
	}
	if cred.ClientSecret.Reveal() != "secret1" {
		t.Errorf("expected secret 'secret1', got '%s'", cred.ClientSecret.Reveal())
	}

	// Convert back to on-disk format
	backToDisk := auraCredentials.toOnDisk()

	if backToDisk.DefaultCredential != "default" {
		t.Errorf("expected default 'default' after round-trip, got '%s'", backToDisk.DefaultCredential)
	}
	if len(backToDisk.Credentials) != 1 {
		t.Fatalf("expected 1 credential after round-trip, got %d", len(backToDisk.Credentials))
	}

	credDisk := backToDisk.Credentials[0]
	if credDisk.ClientSecret != "secret1" {
		t.Errorf("expected secret 'secret1' after round-trip, got '%s'", credDisk.ClientSecret)
	}
}

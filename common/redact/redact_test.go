// Copyright (c) "Neo4j"
// Neo4j Sweden AB [http://neo4j.com]

package redact

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSecretString(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		wantMask  bool
		wantValue string
	}{
		{
			name:     "short secret",
			value:    "abc123",
			wantMask: true,
		},
		{
			name:     "long secret",
			value:    "very-long-secret-value-with-many-characters",
			wantMask: true,
		},
		{
			name:     "empty secret",
			value:    "",
			wantMask: true,
		},
		{
			name:     "secret with special chars",
			value:    "secret!@#$%^&*()",
			wantMask: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSecret(tt.value)
			got := s.String()
			assert.Equal(t, mask, got, "String()")
			if len(tt.value) > 0 {
				assert.NotEqual(t, tt.value, got, "String() leaked the actual value")
			}
		})
	}
}

func TestSecretMarshalJSON(t *testing.T) {
	tests := []struct {
		name       string
		value      string
		wantMask   bool
		wantString string
	}{
		{
			name:       "short secret",
			value:      "secret123",
			wantMask:   true,
			wantString: `"****"`,
		},
		{
			name:       "long secret",
			value:      "very-long-secret-that-should-be-masked",
			wantMask:   true,
			wantString: `"****"`,
		},
		{
			name:       "empty secret",
			value:      "",
			wantMask:   true,
			wantString: `"****"`,
		},
		{
			name:       "secret with quotes",
			value:      `secret"with"quotes`,
			wantMask:   true,
			wantString: `"****"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSecret(tt.value)
			got, err := s.MarshalJSON()
			require.NoError(t, err, "MarshalJSON()")
			gotString := string(got)
			assert.Equal(t, tt.wantString, gotString, "MarshalJSON()")
			assert.NotEqual(t, `"`+tt.value+`"`, gotString, "MarshalJSON() leaked the actual value")
		})
	}
}

func TestSecretReveal(t *testing.T) {
	tests := []struct {
		name  string
		value string
	}{
		{
			name:  "short secret",
			value: "abc123",
		},
		{
			name:  "long secret",
			value: "very-long-secret-with-many-characters",
		},
		{
			name:  "empty secret",
			value: "",
		},
		{
			name:  "secret with special chars",
			value: "secret!@#$%^&*()",
		},
		{
			name:  "secret with spaces",
			value: "secret with spaces",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSecret(tt.value)
			got := s.Reveal()
			assert.Equal(t, tt.value, got, "Reveal()")
		})
	}
}

func TestSecretInJSON(t *testing.T) {
	tests := []struct {
		name         string
		secretValue  string
		expectMasked bool
	}{
		{
			name:         "secret in struct",
			secretValue:  "my-secret-key",
			expectMasked: true,
		},
		{
			name:         "another secret",
			secretValue:  "another-value-123",
			expectMasked: true,
		},
	}

	type TestStruct struct {
		Secret Secret `json:"secret"`
		Safe   string `json:"safe"`
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := TestStruct{
				Secret: NewSecret(tt.secretValue),
				Safe:   "visible",
			}
			data, err := json.Marshal(ts)
			require.NoError(t, err, "json.Marshal")

			var result map[string]interface{}
			err = json.Unmarshal(data, &result)
			require.NoError(t, err, "json.Unmarshal")

			secretVal, ok := result["secret"].(string)
			require.True(t, ok, "secret field is not a string")
			assert.Equal(t, mask, secretVal, "JSON marshaled secret")
			assert.NotEqual(t, tt.secretValue, secretVal, "JSON marshaled secret leaked actual value")

			safeVal, ok := result["safe"].(string)
			require.True(t, ok, "safe field is not a string")
			assert.Equal(t, "visible", safeVal, "Safe field")
		})
	}
}

func TestMaskArgs(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		want     []string
		testName string
	}{
		{
			name:     "empty args",
			args:     []string{},
			want:     []string{},
			testName: "empty arguments",
		},
		{
			name:     "no flags",
			args:     []string{"command", "subcommand"},
			want:     []string{"command", "subcommand"},
			testName: "arguments with no flags",
		},
		{
			name:     "safe flag with value",
			args:     []string{"--name", "my-instance"},
			want:     []string{"--name", "my-instance"},
			testName: "safe flag value not masked",
		},
		{
			name:     "unsafe flag with value",
			args:     []string{"--client-secret", "super-secret-123"},
			want:     []string{"--client-secret", "****"},
			testName: "unsafe flag value masked",
		},
		{
			name:     "multiple flags mixed",
			args:     []string{"--name", "my-app", "--client-secret", "secret123", "--output", "json"},
			want:     []string{"--name", "my-app", "--client-secret", "****", "--output", "json"},
			testName: "multiple flags with mix of safe and unsafe",
		},
		{
			name:     "instance-password unsafe",
			args:     []string{"--instance-username", "admin", "--instance-password", "pass123"},
			want:     []string{"--instance-username", "admin", "--instance-password", "****"},
			testName: "instance password masked, username visible",
		},
		{
			name:     "tenant-id safe",
			args:     []string{"--tenant-id", "tenant-abc-123"},
			want:     []string{"--tenant-id", "tenant-abc-123"},
			testName: "tenant ID not masked",
		},
		{
			name:     "db-id safe",
			args:     []string{"--db-id", "db-12345"},
			want:     []string{"--db-id", "db-12345"},
			testName: "database ID not masked",
		},
		{
			name:     "dbid safe",
			args:     []string{"--dbid", "db-12345"},
			want:     []string{"--dbid", "db-12345"},
			testName: "agent database ID not masked",
		},
		{
			name:     "organization-id shorthand safe",
			args:     []string{"-o", "org-12345"},
			want:     []string{"-o", "org-12345"},
			testName: "organization ID shorthand not masked",
		},
		{
			name:     "project-id shorthand safe",
			args:     []string{"-p", "project-12345"},
			want:     []string{"-p", "project-12345"},
			testName: "project ID shorthand not masked",
		},
		{
			name:     "deployment-id shorthand safe",
			args:     []string{"-d", "deployment-12345"},
			want:     []string{"-d", "deployment-12345"},
			testName: "deployment ID shorthand not masked",
		},
		{
			name:     "server-id shorthand safe",
			args:     []string{"-s", "server-12345"},
			want:     []string{"-s", "server-12345"},
			testName: "server ID shorthand not masked",
		},
		{
			name:     "flag with equals syntax safe",
			args:     []string{"--name=my-instance"},
			want:     []string{"--name=my-instance"},
			testName: "safe flag with equals syntax",
		},
		{
			name:     "flag with equals syntax unsafe",
			args:     []string{"--client-secret=secret123"},
			want:     []string{"--client-secret=****"},
			testName: "unsafe flag with equals syntax",
		},
		{
			name:     "trailing flag without value",
			args:     []string{"--await"},
			want:     []string{"--await"},
			testName: "boolean flag without value",
		},
		{
			name:     "multiple unsafe flags",
			args:     []string{"--client-secret", "s1", "--password", "s2", "--instance-password", "s3"},
			want:     []string{"--client-secret", "****", "--password", "****", "--instance-password", "****"},
			testName: "multiple different unsafe flags",
		},
		{
			name:     "safe flag followed by unsafe",
			args:     []string{"--name", "safe-value", "--client-secret", "unsafe-value"},
			want:     []string{"--name", "safe-value", "--client-secret", "****"},
			testName: "safe and unsafe flags in sequence",
		},
		{
			name:     "special characters in values",
			args:     []string{"--client-secret", "p@ssw0rd!#$%", "--name", "my-app-v1"},
			want:     []string{"--client-secret", "****", "--name", "my-app-v1"},
			testName: "special characters in flag values",
		},
		{
			name:     "graph-analytics-plugin safe",
			args:     []string{"--graph-analytics-plugin", "true"},
			want:     []string{"--graph-analytics-plugin", "true"},
			testName: "graph-analytics-plugin is a safe flag, value left untouched",
		},
		{
			name:     "version safe",
			args:     []string{"--version", "5"},
			want:     []string{"--version", "5"},
			testName: "version flag not masked",
		},
		{
			name:     "region safe",
			args:     []string{"--region", "us-east-1"},
			want:     []string{"--region", "us-east-1"},
			testName: "region flag not masked",
		},
		{
			name:     "memory safe",
			args:     []string{"--memory", "64GB"},
			want:     []string{"--memory", "64GB"},
			testName: "memory flag not masked",
		},
		{
			name:     "cloud-provider safe",
			args:     []string{"--cloud-provider", "aws"},
			want:     []string{"--cloud-provider", "aws"},
			testName: "cloud-provider not masked",
		},
		{
			name:     "auth-url safe",
			args:     []string{"--auth-url", "https://auth.example.com"},
			want:     []string{"--auth-url", "https://auth.example.com"},
			testName: "auth URL not masked",
		},
		{
			name:     "base-url safe",
			args:     []string{"--base-url", "https://api.example.com"},
			want:     []string{"--base-url", "https://api.example.com"},
			testName: "base URL not masked",
		},
		{
			name:     "output safe",
			args:     []string{"--output", "json"},
			want:     []string{"--output", "json"},
			testName: "output format not masked",
		},
		{
			name:     "data-api-id safe",
			args:     []string{"--data-api-id", "api-xyz-789"},
			want:     []string{"--data-api-id", "api-xyz-789"},
			testName: "data API ID not masked",
		},
		{
			name:     "customer-managed-key-id safe",
			args:     []string{"--customer-managed-key-id", "key-abc-123"},
			want:     []string{"--customer-managed-key-id", "key-abc-123"},
			testName: "customer managed key ID not masked",
		},
		{
			name:     "unsafe flag with dash-prefixed secret",
			args:     []string{"--client-secret", "-generated-secret-starting-with-dash"},
			want:     []string{"--client-secret", "****"},
			testName: "unsafe flag value starting with dash is masked",
		},
		{
			name:     "unsafe boolean flag followed by safe named flag",
			args:     []string{"--enabled", "--name", "my-api"},
			want:     []string{"--enabled", "--name", "my-api"},
			testName: "boolean flag doesn't consume next flag as its value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			got := MaskArgs(tt.args)
			require.Equal(t, len(tt.want), len(got), "MaskArgs() element count")

			for i, gotArg := range got {
				assert.Equal(t, tt.want[i], gotArg, "MaskArgs()[%d]", i)
			}
		})
	}
}

func TestMaskArgsMasksSecrets(t *testing.T) {
	tests := []struct {
		name       string
		flagName   string
		shouldMask bool
	}{
		{"client-secret", "client-secret", true},
		{"password", "password", true},
		{"instance-password", "instance-password", true},
		{"instance-username", "instance-username", false},
		{"name", "name", false},
		{"instance-id", "instance-id", false},
		{"tenant-id", "tenant-id", false},
		{"output", "output", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := []string{"--" + tt.flagName, "test-value"}
			result := MaskArgs(args)

			require.Len(t, result, 2, "MaskArgs() element count")

			if tt.shouldMask {
				assert.Equal(t, mask, result[1], "flag %s", tt.flagName)
			} else {
				assert.NotEqual(t, mask, result[1], "flag %s: got masked value, should not be masked", tt.flagName)
				assert.Equal(t, "test-value", result[1], "flag %s", tt.flagName)
			}
		})
	}
}

func TestMaskArgsPreservesStructure(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{
			name: "complex command with multiple flags",
			args: []string{
				"instance", "create",
				"--name", "my-instance",
				"--version", "5",
				"--region", "us-east-1",
				"--memory", "64GB",
				"--cloud-provider", "aws",
				"--client-secret", "supersecret",
				"--tenant-id", "tenant-123",
				"--output", "json",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MaskArgs(tt.args)

			require.Len(t, result, len(tt.args), "MaskArgs() changed argument count")

			for i, arg := range tt.args {
				if strings.HasPrefix(arg, "--") {
					assert.Equal(t, arg, result[i], "flag name at index %d changed", i)
				}
			}
		})
	}
}

// Copyright (c) "Neo4j"
// Neo4j Sweden AB [http://neo4j.com]

package redact

import (
	"encoding/json"
	"strings"
)

const mask = "****"

// Secret is a string wrapper that masks its value when printed or marshaled.
// The only way to retrieve the real value is via Reveal().
type Secret struct {
	value string
}

// NewSecret creates a new Secret with the given value.
func NewSecret(value string) Secret {
	return Secret{value: value}
}

// String returns the masked representation.
func (s Secret) String() string {
	return mask
}

// MarshalJSON returns the masked representation as JSON.
func (s Secret) MarshalJSON() ([]byte, error) {
	return json.Marshal(mask)
}

// Reveal returns the actual secret value.
func (s Secret) Reveal() string {
	return s.value
}

// safeFlags is a map of flag names that are safe to echo verbatim in diagnostic output.
// Any flag not in this map is considered unsafe and will be masked by Args().
var safeFlags = map[string]bool{
	// Identifiers and names
	"name":                     true,
	"instance-id":              true,
	"tenant-id":                true,
	"customer-managed-key-id":  true,
	"db-id":                    true,
	"data-api-id":              true,
	"type":                     true,
	"deployment-id":            true,
	"key-id":                   true,
	"server-id":                true,
	"source-instance-id":       true,
	"source-snapshot-id":       true,
	"description":              true,
	"project-id":               true,
	"organization-id":          true,
	"import-model-id":          true,

	// Output and format options
	"output":                   true,

	// Configuration endpoints (not secrets themselves)
	"auth-url":                 true,
	"base-url":                 true,

	// Feature flags and non-secret toggles
	"await":                    true,
	"graph-analytics-plugin":   true,

	// Instance configuration (non-secret metadata)
	"version":                  true,
	"region":                   true,
	"memory":                   true,
	"cloud-provider":           true,

	// Usernames (not secret by themselves, but paired with passwords)
	"instance-username":        true,

	// Data/Time related
	"date":                     true,
	"ttl":                      true,

	// Data import
	"import-type":              true,
}

// Args takes a slice of command-line arguments (as os.Args[1:] is shaped)
// and returns a copy with values following unsafe flags replaced by the mask.
// Allow-listed flag names and their values are left untouched.
func Args(args []string) []string {
	if len(args) == 0 {
		return args
	}

	result := make([]string, 0, len(args))

	for i := 0; i < len(args); i++ {
		arg := args[i]
		result = append(result, arg)

		// Check if this is a flag (starts with -)
		if strings.HasPrefix(arg, "-") {
			// Strip leading dashes and extract flag name
			flagName := strings.TrimLeft(arg, "-")

			// Handle flags with = (e.g., --flag=value)
			if strings.Contains(flagName, "=") {
				parts := strings.SplitN(flagName, "=", 2)
				flagName = parts[0]
				value := parts[1]

				// Remove the old arg and replace with masked version if needed
				result[len(result)-1] = arg[:strings.Index(arg, "=")+1] + maskIfUnsafe(flagName, value)
			} else if i+1 < len(args) {
				// For unsafe flags, mask the next argument regardless of whether it starts with dash.
				// For safe flags, only consume the next argument if it doesn't look like a flag.
				isSafeFlag := safeFlags[flagName]
				nextLooksLikeFlag := strings.HasPrefix(args[i+1], "-")
				if !isSafeFlag || !nextLooksLikeFlag {
					i++
					value := args[i]
					if isSafeFlag {
						result = append(result, value)
					} else {
						result = append(result, mask)
					}
				}
			}
		}
	}

	return result
}

func maskIfUnsafe(flagName, value string) string {
	if safeFlags[flagName] {
		return value
	}
	return mask
}

// Copyright (c) "Neo4j"
// Neo4j Sweden AB [http://neo4j.com]

package redact

import (
	"encoding/json"
	"strings"
)

const mask = "****"

type Secret struct {
	value string
}

func NewSecret(value string) Secret {
	return Secret{value: value}
}

func (s Secret) String() string {
	return mask
}

func (s Secret) MarshalJSON() ([]byte, error) {
	return json.Marshal(mask)
}

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

// booleanFlags is a set of flag names that don't take a value (boolean flags).
// These flags should never have their following argument consumed as a value.
var booleanFlags = map[string]bool{
	"await":                    true,
	"enabled":                  true,
	"disabled":                 true,
	"is-private":               true,
	"is-mcp-enabled":           true,
	"vector-optimized":         true,
	"graph-analytics-plugin":   true,
	"progress":                 true,
	"show-progress":            true,
}

func Args(args []string) []string {
	if len(args) == 0 {
		return args
	}

	result := make([]string, 0, len(args))

	for i := 0; i < len(args); i++ {
		arg := args[i]
		result = append(result, arg)

		if strings.HasPrefix(arg, "-") {
			flagName := strings.TrimLeft(arg, "-")

			if strings.Contains(flagName, "=") {
				parts := strings.SplitN(flagName, "=", 2)
				flagName = parts[0]
				value := parts[1]

				result[len(result)-1] = arg[:strings.Index(arg, "=")+1] + maskIfUnsafe(flagName, value)
			} else if i+1 < len(args) && !booleanFlags[flagName] {
				// Boolean flags don't take values, so never consume the next argument.
				// For other unsafe flags, mask the next argument.
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

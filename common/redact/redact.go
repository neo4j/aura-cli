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

var safeFlags = map[string]bool{
	// Identifiers and names
	"name":                    true,
	"instance-id":             true,
	"tenant-id":               true,
	"customer-managed-key-id": true,
	"db-id":                   true,
	"data-api-id":             true,
	"type":                    true,
	"deployment-id":           true,
	"key-id":                  true,
	"server-id":               true,
	"source-instance-id":      true,
	"source-snapshot-id":      true,
	"description":             true,
	"project-id":              true,
	"organization-id":         true,
	"import-model-id":         true,

	// Output and format options
	"output": true,

	// Configuration endpoints (not secrets themselves)
	"auth-url": true,
	"base-url": true,

	// Feature flags and non-secret toggles
	"await":                  true,
	"graph-analytics-plugin": true,

	// Instance configuration (non-secret metadata)
	"version":        true,
	"region":         true,
	"memory":         true,
	"cloud-provider": true,

	// Usernames (not secret by themselves, but paired with passwords)
	"instance-username": true,

	// Data/Time related
	"date": true,
	"ttl":  true,

	// Data import
	"import-type": true,
}

var booleanFlags = map[string]bool{
	"await":                  true,
	"enabled":                true,
	"disabled":               true,
	"is-private":             true,
	"is-mcp-enabled":         true,
	"vector-optimized":       true,
	"graph-analytics-plugin": true,
	"progress":               true,
	"help":                   true,
}

func parseFlagFromArg(arg string) (flagName string, inlineValue string, hasInlineValue bool) {
	flagName = strings.TrimLeft(arg, "-")
	if strings.Contains(flagName, "=") {
		parts := strings.SplitN(flagName, "=", 2)
		return parts[0], parts[1], true
	}
	return flagName, "", false
}

// maskArg handles the arg at the front of remaining, returning its masked output tokens
// and how many of remaining's leading elements they account for (1, or 2 if a following
// value was consumed).
func maskArg(remaining []string) (output []string, consumed int) {
	arg := remaining[0]

	if !strings.HasPrefix(arg, "-") {
		return []string{arg}, 1
	}

	flagName, inlineValue, hasInlineValue := parseFlagFromArg(arg)

	if hasInlineValue {
		masked := arg[:strings.Index(arg, "=")+1] + maskIfUnsafe(flagName, inlineValue)
		return []string{masked}, 1
	}

	if booleanFlags[flagName] {
		return []string{arg}, 1
	}

	if len(remaining) < 2 {
		return []string{arg}, 1
	}

	isSafeFlag := safeFlags[flagName]
	nextLooksLikeFlag := strings.HasPrefix(remaining[1], "-")
	if isSafeFlag && nextLooksLikeFlag {
		return []string{arg}, 1
	}

	value := remaining[1]
	if isSafeFlag {
		return []string{arg, value}, 2
	}
	return []string{arg, mask}, 2
}

func MaskArgs(args []string) []string {
	if len(args) == 0 {
		return args
	}

	result := make([]string, 0, len(args))

	remaining := args
	for len(remaining) > 0 {
		output, consumed := maskArg(remaining)
		result = append(result, output...)
		remaining = remaining[consumed:]
	}

	return result
}

func maskIfUnsafe(flagName, value string) string {
	if safeFlags[flagName] {
		return value
	}
	return mask
}

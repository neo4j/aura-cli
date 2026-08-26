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
	"dbid":                    true,

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

// ShorthandResolver maps a flag shorthand (e.g. "p") to its long name (e.g.
// "project-id") for the command actually being invoked.
type ShorthandResolver func(flagName string) (longName string, ok bool)

// parseFlagFromArg also recognises pflag's attached-shorthand form
// ("-xvalue", no "=") - without it, a secret passed that way would be left
// on the token entirely unmasked.
func parseFlagFromArg(arg string, resolve ShorthandResolver) (flagName string, inlineValue string, hasInlineValue bool) {
	trimmed := strings.TrimLeft(arg, "-")

	if !strings.HasPrefix(arg, "--") && len(trimmed) > 1 && trimmed[1] != '=' {
		flagName, inlineValue, hasInlineValue = trimmed[:1], trimmed[1:], true
	} else {
		flagName = trimmed
	}

	if strings.Contains(flagName, "=") {
		parts := strings.SplitN(flagName, "=", 2)
		flagName, inlineValue, hasInlineValue = parts[0], parts[1], true
	}

	if resolve != nil {
		if longName, ok := resolve(flagName); ok {
			flagName = longName
		}
	}

	return flagName, inlineValue, hasInlineValue
}

func maskArg(remainingArgs []string, resolve ShorthandResolver) (output []string, consumed int) {
	arg := remainingArgs[0]

	if arg == "--" || !strings.HasPrefix(arg, "-") {
		return []string{arg}, 1
	}

	flagName, inlineValue, hasInlineValue := parseFlagFromArg(arg, resolve)

	if hasInlineValue {
		masked := arg[:len(arg)-len(inlineValue)] + maskIfUnsafe(flagName, inlineValue)
		return []string{masked}, 1
	}

	if booleanFlags[flagName] {
		return []string{arg}, 1
	}

	if len(remainingArgs) < 2 {
		return []string{arg}, 1
	}

	isSafeFlag := safeFlags[flagName]
	nextLooksLikeFlag := strings.HasPrefix(remainingArgs[1], "-")
	if isSafeFlag && nextLooksLikeFlag {
		return []string{arg}, 1
	}

	// Consume the next token as the value even if it looks like a flag:
	// generated secrets can themselves start with a dash.
	value := remainingArgs[1]
	if isSafeFlag {
		return []string{arg, value}, 2
	}
	return []string{arg, mask}, 2
}

func MaskArgs(args []string) []string {
	return MaskArgsWithShorthandResolver(args, nil)
}

// MaskArgsWithShorthandResolver is MaskArgs, but resolves shorthands via resolve first.
func MaskArgsWithShorthandResolver(args []string, resolve ShorthandResolver) []string {
	if len(args) == 0 {
		return args
	}

	result := make([]string, 0, len(args))

	remaindingArgs := args
	for len(remaindingArgs) > 0 {
		output, consumed := maskArg(remaindingArgs, resolve)
		result = append(result, output...)
		remaindingArgs = remaindingArgs[consumed:]
	}

	return result
}

func maskIfUnsafe(flagName, value string) string {
	if safeFlags[flagName] {
		return value
	}
	return mask
}

// Copyright (c) "Neo4j"
// Neo4j Sweden AB [http://neo4j.com]

package query

// truncateArrays recursively walks value, replacing any []any whose length
// exceeds max with an empty slice `[]any{}`. Maps (`map[string]any`) and
// slices (`[]any`) are descended into; all other types are returned
// unchanged. When max <= 0, the input is returned as-is (no truncation
// applied).
//
// The empty-slice signal is type-safe (downstream JSON consumers see an
// array, not a string) and unambiguously marks the value as elided.
//
// The function is pure: it does not mutate its input. Slices and maps that
// require modification are reallocated; values that pass through untouched
// are returned by reference.
func truncateArrays(value any, max int) any {
	if max <= 0 {
		return value
	}
	return truncate(value, max)
}

func truncate(value any, max int) any {
	switch v := value.(type) {
	case []any:
		if len(v) > max {
			return []any{}
		}
		out := make([]any, len(v))
		for i, item := range v {
			out[i] = truncate(item, max)
		}
		return out
	case map[string]any:
		out := make(map[string]any, len(v))
		for k, item := range v {
			out[k] = truncate(item, max)
		}
		return out
	default:
		return v
	}
}

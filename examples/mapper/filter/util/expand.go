package util

import (
	"strings"
)

func Expand(value map[string]any) map[string]any {
	return ExpandPrefixed(value, "")
}

func ExpandPrefixed(value map[string]any, prefix string) map[string]any {
	m := make(map[string]any)
	ExpandPrefixedToResult(value, prefix, m)
	return m
}

func ExpandPrefixedToResult(value map[string]any, prefix string, result map[string]any) {
	if prefix != "" {
		prefix += "."
	}
	for k, val := range value {
		if !strings.HasPrefix(k, prefix) {
			continue
		}

		key := k[len(prefix):]
		idx := strings.Index(key, ".")
		if idx != -1 {
			key = key[:idx]
		}
		// It is possible for the map to contain conflicts:
		// {"a.b": 1, "a": 2}
		// What should the result be? We overwrite the less-specific key.
		// {"a.b": 1, "a": 2} -> {"a.b": 1, "a": 2}
		if _, ok := result[key]; ok && idx == -1 {
			continue
		}
		if idx == -1 {
			result[key] = val
			continue
		}

		// It contains a period, so it is a more complex structure
		result[key] = ExpandPrefixed(value, k[:len(prefix)+len(key)])
	}
}

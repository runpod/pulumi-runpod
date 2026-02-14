package runpod

import (
	"fmt"
	"strings"
)

// EnvMapToGQL converts a map of env vars to the GraphQL input format.
// Returns an empty slice (not nil) when no env vars — many GQL mutations require non-null.
func EnvMapToGQL(env map[string]string) []EnvironmentVariable {
	result := make([]EnvironmentVariable, 0, len(env))
	for k, v := range env {
		result = append(result, EnvironmentVariable{Key: k, Value: v})
	}
	return result
}

// EnvSliceToMap converts the "KEY=VALUE" string slice from API responses to a map.
func EnvSliceToMap(env []string) map[string]string {
	if len(env) == 0 {
		return nil
	}
	result := make(map[string]string, len(env))
	for _, e := range env {
		parts := strings.SplitN(e, "=", 2)
		if len(parts) == 2 {
			result[parts[0]] = parts[1]
		}
	}
	return result
}

// EnvGQLToMap converts EnvironmentVariable slice to a map.
func EnvGQLToMap(env []EnvironmentVariable) map[string]string {
	if len(env) == 0 {
		return nil
	}
	result := make(map[string]string, len(env))
	for _, e := range env {
		result[e.Key] = e.Value
	}
	return result
}

// StringPtr returns a pointer to a string, or nil if empty.
func StringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// IntPtr returns a pointer to an int, or nil if zero.
func IntPtr(i int) *int {
	if i == 0 {
		return nil
	}
	return &i
}

// BoolPtr returns a pointer to a bool.
func BoolPtr(b bool) *bool {
	return &b
}

// PtrString returns the string value of a pointer, or empty string if nil.
func PtrString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// FormatError creates a formatted error with resource context.
func FormatError(operation, resource, id string, err error) error {
	if id != "" {
		return fmt.Errorf("%s %s %q: %w", operation, resource, id, err)
	}
	return fmt.Errorf("%s %s: %w", operation, resource, err)
}

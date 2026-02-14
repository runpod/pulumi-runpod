package runpod

import (
	"fmt"
	"strings"
)

// EnvMapToGQL converts a map of env vars to genqlient's EnvironmentVariableInput slice.
// Returns an empty slice (not nil) when no env vars — many GQL mutations require non-null.
func EnvMapToGQL(env map[string]string) []*EnvironmentVariableInput {
	result := make([]*EnvironmentVariableInput, 0, len(env))
	for k, v := range env {
		result = append(result, &EnvironmentVariableInput{Key: k, Value: v})
	}
	return result
}

// EnvSliceToMap converts the "KEY=VALUE" string pointer slice from API responses to a map.
func EnvSliceToMap(env []*string) map[string]string {
	if len(env) == 0 {
		return nil
	}
	result := make(map[string]string, len(env))
	for _, e := range env {
		if e == nil {
			continue
		}
		parts := strings.SplitN(*e, "=", 2)
		if len(parts) == 2 {
			result[parts[0]] = parts[1]
		}
	}
	return result
}

// EnvGQLToMap converts TemplateResponseEnvEnvironmentVariable slice to a map.
func EnvGQLToMap(env []*TemplateResponseEnvEnvironmentVariable) map[string]string {
	if len(env) == 0 {
		return nil
	}
	result := make(map[string]string, len(env))
	for _, e := range env {
		if e == nil || e.Key == nil {
			continue
		}
		val := ""
		if e.Value != nil {
			val = *e.Value
		}
		result[*e.Key] = val
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

// PtrInt returns the int value of a pointer, or 0 if nil.
func PtrInt(i *int) int {
	if i == nil {
		return 0
	}
	return *i
}

// PtrBool returns the bool value of a pointer, or false if nil.
func PtrBool(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}

// PtrFloat64 returns the float64 value of a pointer, or 0 if nil.
func PtrFloat64(f *float64) float64 {
	if f == nil {
		return 0
	}
	return *f
}

// FormatError creates a formatted error with resource context.
func FormatError(operation, resource, id string, err error) error {
	if id != "" {
		return fmt.Errorf("%s %s %q: %w", operation, resource, id, err)
	}
	return fmt.Errorf("%s %s: %w", operation, resource, err)
}

// StringPtrSlice converts a []string to []*string.
func StringPtrSlice(ss []string) []*string {
	result := make([]*string, len(ss))
	for i := range ss {
		result[i] = &ss[i]
	}
	return result
}

// DerefStringSlice converts []*string to []string, skipping nils.
func DerefStringSlice(ss []*string) []string {
	if len(ss) == 0 {
		return nil
	}
	result := make([]string, 0, len(ss))
	for _, s := range ss {
		if s != nil {
			result = append(result, *s)
		}
	}
	return result
}

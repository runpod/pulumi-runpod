//go:build ignore

// This script applies post-genqlient fixes to generated.go:
//
//  1. Adds omitempty to scalar/struct pointer fields so nil is omitted from
//     JSON (not sent as null). RunPod's API distinguishes null from absent
//     for fields like "repo" — sending null triggers validation errors.
//
//  2. Keeps "env" fields WITHOUT omitempty. RunPod mutations require
//     env:[EnvironmentVariableInput]! (non-null array), so we must always
//     send [] rather than omitting the field entirely.
package main

import (
	"os"
	"regexp"
	"strings"
)

func main() {
	path := "generated.go"
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	// Step 1: add omitempty to bare pointer fields (*T `json:"name"`)
	// Only match lines where the field type starts with a literal * (not []*).
	// We use a line-by-line approach to avoid matching slice-of-pointer types.
	lines := strings.Split(string(data), "\n")
	rePointerField := regexp.MustCompile(
		`^(\s+\w+\s+)\*\S+(\s+` + "`" + `json:"([^"]+)")` + "`",
	)
	for i, line := range lines {
		if rePointerField.MatchString(line) && !strings.Contains(line, "omitempty") {
			lines[i] = strings.Replace(line, `json:"`, `json:"`, 1) // no-op anchor
			lines[i] = regexp.MustCompile(`(json:"[^"]+)"` + "`").
				ReplaceAllString(lines[i], `${1},omitempty"`+"`")
		}
	}
	result := strings.Join(lines, "\n")

	// Step 2: strip omitempty back off any "env" fields — the API requires
	// non-null env arrays, so we must send [] not omit the field.
	reEnvOmit := regexp.MustCompile(`(json:"env),omitempty"`)
	result = reEnvOmit.ReplaceAllString(result, `${1}"`)

	if err := os.WriteFile(path, []byte(result), 0o644); err != nil {
		panic(err)
	}
}

//go:build ignore

// This script adds omitempty to all pointer-type JSON tags in genqlient's
// generated.go. Without this, nil pointer fields get serialized as null
// which some GraphQL APIs misinterpret (e.g. RunPod treats repo:null as
// "create a repo-based endpoint").
package main

import (
	"os"
	"regexp"
)

func main() {
	path := "generated.go"
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	// Match json tags on pointer fields that don't already have omitempty.
	// Pattern: a pointer type (*Something) followed by a json tag without omitempty.
	re := regexp.MustCompile(`(\*\S+\s+` + "`" + `json:"[^"]+)"` + "`")
	result := re.ReplaceAllFunc(data, func(match []byte) []byte {
		s := string(match)
		// Don't double-add
		if regexp.MustCompile(`omitempty`).MatchString(s) {
			return match
		}
		// Insert ,omitempty before the closing quote of the json tag
		reInner := regexp.MustCompile(`(json:"[^"]+)"`)
		return []byte(reInner.ReplaceAllString(s, `${1},omitempty"`))
	})

	if err := os.WriteFile(path, result, 0o644); err != nil {
		panic(err)
	}
}

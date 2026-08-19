// Copyright (c) "Neo4j"
// Neo4j Sweden AB [http://neo4j.com]

package redact

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"testing"
)

// TestArgScopeRegression ensures that raw os.Args[1:] reads are confined to
// the sanctioned capture points in neo4j-cli/aura/cmd/main.go and neo4j-cli/main.go (added via task-003)
// and never leak elsewhere in the codebase. This guards against regressions where
// unredacted argument slices reach panic/error output paths.
func TestArgScopeRegression(t *testing.T) {
	// Sanctioned locations where os.Args[1:] (or os.Args slicing) is permitted.
	// Task-003 introduced capture points in both main() entrypoints.
	sanctionedLocations := map[string]bool{
		// This test file itself (guard definition)
		"common/redact/argscope_test.go": true,
		// Sanctioned capture points (task-003)
		"neo4j-cli/aura/cmd/main.go": true,
		"neo4j-cli/main.go":           true,
	}

	// Pattern to find raw os.Args reads: os.Args[...] where ... is any slice/index expression
	// This catches: os.Args[1:], os.Args[0], os.Args[1:], etc.
	osArgsPattern := regexp.MustCompile(`\bos\.Args\s*\[`)

	repoRoot := findRepoRoot(t)
	violations := []string{}

	// Walk the repository and examine all .go files
	err := filepath.Walk(repoRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and non-.go files
		if info.IsDir() {
			// Skip hidden directories (.git, .claude, etc.)
			if strings.HasPrefix(info.Name(), ".") {
				return filepath.SkipDir
			}
			return nil
		}

		if !strings.HasSuffix(path, ".go") {
			return nil
		}

		// Get the relative path from repo root for checking against sanctioned locations
		relPath, err := filepath.Rel(repoRoot, path)
		if err != nil {
			return err
		}

		// Normalize path separators to forward slashes for consistency
		relPath = filepath.ToSlash(relPath)

		// Check if this file is in a sanctioned location
		if sanctionedLocations[relPath] {
			return nil
		}

		// Examine the file line by line for os.Args[ patterns
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		lineNum := 0
		for scanner.Scan() {
			lineNum++
			line := scanner.Text()
			trimmedLine := strings.TrimSpace(line)

			// Skip comment-only lines
			if strings.HasPrefix(trimmedLine, "//") {
				continue
			}

			// Check for os.Args[ pattern (catches all slice/index operations)
			if osArgsPattern.MatchString(line) {
				// Found a raw os.Args read outside a sanctioned location
				violations = append(violations,
					fmt.Sprintf("%s:%d: raw os.Args read: %s",
						relPath, lineNum, strings.TrimSpace(line)))
			}
		}

		if err := scanner.Err(); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		t.Fatalf("Failed to walk repository: %v", err)
	}

	// If violations were found, fail with a detailed report
	if len(violations) > 0 {
		t.Errorf("Found raw os.Args reads outside sanctioned locations:\n%s",
			strings.Join(violations, "\n"))
	}
}

// findRepoRoot finds the root of the aura-cli repository by walking up from the test file.
func findRepoRoot(t *testing.T) string {
	// Get the directory of this test file
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Unable to determine test file location")
	}

	current := filepath.Dir(filename)

	// Walk up until we find a go.mod file
	for {
		gomod := filepath.Join(current, "go.mod")
		if _, err := os.Stat(gomod); err == nil {
			// Found the go.mod, this should be our repo root
			// But we want the parent if we're in a subdirectory with go.mod
			// For aura-cli, the root should have go.mod at the top level
			// Let's verify by checking for typical repo markers
			if fileExists(filepath.Join(current, ".git")) {
				return current
			}
		}

		parent := filepath.Dir(current)
		if parent == current {
			// Reached filesystem root without finding .git
			t.Fatalf("Could not find repository root")
		}
		current = parent
	}
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

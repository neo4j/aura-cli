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

func TestArgScopeRegression(t *testing.T) {
	sanctionedLocations := map[string]bool{
		"common/redact/argscope_test.go": true,
		"neo4j-cli/aura/cmd/main.go":     true,
		"neo4j-cli/main.go":              true,
	}

	osArgsPattern := regexp.MustCompile(`\bos\.Args\b`)

	repoRoot := findRepoRoot(t)
	violations := []string{}

	err := filepath.Walk(repoRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if strings.HasPrefix(info.Name(), ".") {
				return filepath.SkipDir
			}
			return nil
		}

		if !strings.HasSuffix(path, ".go") {
			return nil
		}

		relPath, err := filepath.Rel(repoRoot, path)
		if err != nil {
			return err
		}

		relPath = filepath.ToSlash(relPath)

		if sanctionedLocations[relPath] {
			return nil
		}

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

			if strings.HasPrefix(trimmedLine, "//") {
				continue
			}

			codePart := line
			if idx := strings.Index(line, "//"); idx >= 0 {
				codePart = line[:idx]
			}

			if osArgsPattern.MatchString(codePart) {
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

	if len(violations) > 0 {
		t.Errorf("Found raw os.Args reads outside sanctioned locations:\n%s",
			strings.Join(violations, "\n"))
	}
}

func findRepoRoot(t *testing.T) string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Unable to determine test file location")
	}

	current := filepath.Dir(filename)
	for {
		if _, err := os.Stat(filepath.Join(current, ".git")); err == nil {
			return current
		}
		parent := filepath.Dir(current)
		if parent == current {
			t.Fatalf("Could not find repository root")
		}
		current = parent
	}
}

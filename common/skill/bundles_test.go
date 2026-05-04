// Copyright (c) "Neo4j"
// Neo4j Sweden AB [http://neo4j.com]

package skill_test

import (
	"bytes"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCommittedBundlesObeyAnthropicSkillRules walks every committed
// `<bin>/internal/skill/bundle/SKILL.md` under the repo root and asserts the
// Anthropic skill-format constraints from the PRD:
//
//   - frontmatter `name`: ≤64 chars and matches `[a-z0-9-]+`.
//   - frontmatter `description`: non-empty, ≤1024 chars.
//   - frontmatter `version` line present and equals the literal
//     `{{VERSION}}` placeholder on disk (substitution happens at install
//     time, not generation).
//   - SKILL.md body ≤500 lines.
//
// The walk discovers bundles automatically so adding a new binary's bundle
// (e.g. a future `cypher-cli`) gets covered without test edits. Failure
// messages identify the offending bundle path and which rule was breached.
func TestCommittedBundlesObeyAnthropicSkillRules(t *testing.T) {
	repoRoot := repoRoot(t)

	bundles := findCommittedBundles(t, repoRoot)
	require.NotEmpty(t, bundles, "expected at least one committed SKILL.md under <bin>/internal/skill/bundle/")

	for _, path := range bundles {
		path := path
		rel, _ := filepath.Rel(repoRoot, path)
		t.Run(rel, func(t *testing.T) {
			data, err := os.ReadFile(path)
			require.NoErrorf(t, err, "%s: read", rel)

			fm := parseFrontmatter(t, rel, data)

			// REQ-F-007: name ≤64 chars, [a-z0-9-]+ charset.
			name := fm["name"]
			assert.NotEmptyf(t, name, "%s: frontmatter `name` is empty", rel)
			assert.LessOrEqualf(t, len(name), 64,
				"%s: frontmatter `name` (%q) exceeds 64 chars", rel, name)
			assert.Regexpf(t, regexp.MustCompile(`^[a-z0-9-]+$`), name,
				"%s: frontmatter `name` (%q) violates [a-z0-9-]+ charset", rel, name)

			// REQ-F-007: description non-empty, ≤1024 chars.
			desc := fm["description"]
			assert.NotEmptyf(t, desc, "%s: frontmatter `description` is empty", rel)
			assert.LessOrEqualf(t, len(desc), 1024,
				"%s: frontmatter `description` is %d chars (>1024)", rel, len(desc))

			// PRD: on-disk `version` equals literal `{{VERSION}}` placeholder.
			version, ok := fm["version"]
			assert.Truef(t, ok, "%s: frontmatter is missing `version` line", rel)
			assert.Equalf(t, "{{VERSION}}", version,
				"%s: on-disk `version` must be the literal `{{VERSION}}` placeholder (got %q) — substitution happens at install time, not generation",
				rel, version)

			// REQ-F-006 / REQ-T-008: SKILL.md body ≤500 lines.
			lineCount := bytes.Count(data, []byte("\n"))
			if len(data) > 0 && data[len(data)-1] != '\n' {
				lineCount++ // trailing line without terminator still counts
			}
			assert.LessOrEqualf(t, lineCount, 500,
				"%s: SKILL.md is %d lines (>500)", rel, lineCount)
		})
	}
}

// repoRoot returns the absolute path to the repository root, resolved via
// runtime.Caller(0). This file lives at <repo>/common/skill/bundles_test.go,
// so two parent jumps from its directory land at the repo root. Mirrors the
// pattern used in <bin>/internal/skill/gen/main.go.
func repoRoot(t *testing.T) string {
	t.Helper()
	_, thisFile, _, ok := runtime.Caller(0)
	require.True(t, ok, "runtime.Caller(0) failed")
	// thisFile = .../common/skill/bundles_test.go
	// parent of skill/ is common/, parent of common/ is the repo root.
	return filepath.Dir(filepath.Dir(filepath.Dir(thisFile)))
}

// findCommittedBundles walks the repo and returns every
// `<bin>/internal/skill/bundle/SKILL.md` it finds. Using filepath.Walk
// (rather than a hardcoded slice) means new binaries pick up the gate
// automatically. Skips `.git`, `node_modules`, and other heavy trees.
func findCommittedBundles(t *testing.T, repoRoot string) []string {
	t.Helper()
	var out []string
	err := filepath.Walk(repoRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			base := info.Name()
			// Prune dirs that can't contain bundles and would slow the walk.
			if base == ".git" || base == "node_modules" || base == "bin" || base == ".changes" {
				return filepath.SkipDir
			}
			return nil
		}
		if info.Name() != "SKILL.md" {
			return nil
		}
		// Match `<bin>/internal/skill/bundle/SKILL.md`. Use forward-slash
		// comparison so the suffix check works on Windows.
		slash := filepath.ToSlash(path)
		if strings.HasSuffix(slash, "/internal/skill/bundle/SKILL.md") {
			out = append(out, path)
		}
		return nil
	})
	require.NoError(t, err)
	return out
}

// parseFrontmatter extracts simple `key: value` lines between the opening
// and closing `---` fences of a SKILL.md file. Single-line values only —
// matches the generator's output and avoids a YAML dependency. Fails the
// test when the frontmatter block is missing or malformed.
func parseFrontmatter(t *testing.T, rel string, data []byte) map[string]string {
	t.Helper()
	lines := strings.Split(string(data), "\n")
	require.GreaterOrEqualf(t, len(lines), 2, "%s: file too short for frontmatter", rel)
	require.Equalf(t, "---", strings.TrimRight(lines[0], "\r"),
		"%s: missing opening `---` frontmatter fence", rel)

	out := map[string]string{}
	closed := false
	for i := 1; i < len(lines); i++ {
		line := strings.TrimRight(lines[i], "\r")
		if line == "---" {
			closed = true
			break
		}
		idx := strings.IndexByte(line, ':')
		if idx <= 0 {
			continue
		}
		key := strings.TrimSpace(line[:idx])
		value := strings.TrimSpace(line[idx+1:])
		out[key] = value
	}
	require.Truef(t, closed, "%s: missing closing `---` frontmatter fence", rel)
	return out
}

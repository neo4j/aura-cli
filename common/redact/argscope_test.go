// Copyright (c) "Neo4j"
// Neo4j Sweden AB [http://neo4j.com]

package redact

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestArgScopeRegression(t *testing.T) {
	sanctionedLocations := map[string]bool{
		"common/redact/argscope_test.go": true,
		"neo4j-cli/aura/cmd/main.go":     true,
		"neo4j-cli/main.go":              true,
	}

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

		fset := token.NewFileSet()
		file, err := parser.ParseFile(fset, path, nil, 0)
		if err != nil {
			return fmt.Errorf("failed to parse %s: %w", relPath, err)
		}

		ast.Inspect(file, func(n ast.Node) bool {
			sel, ok := n.(*ast.SelectorExpr)
			if !ok {
				return true
			}

			pkgIdent, ok := sel.X.(*ast.Ident)
			if !ok || pkgIdent.Name != "os" || sel.Sel.Name != "Args" {
				return true
			}

			pos := fset.Position(sel.Pos())
			violations = append(violations,
				fmt.Sprintf("%s:%d: raw os.Args read", relPath, pos.Line))
			return true
		})

		return nil
	})

	assert.NoError(t, err, "failed to walk repository")

	assert.Empty(t, violations, "found raw os.Args reads outside sanctioned locations:\n%s",
		strings.Join(violations, "\n"))
}

func findRepoRoot(t *testing.T) string {
	_, filename, _, ok := runtime.Caller(0)
	require.True(t, ok, "unable to determine test file location")

	current := filepath.Dir(filename)
	for {
		if _, err := os.Stat(filepath.Join(current, ".git")); err == nil {
			return current
		}
		parent := filepath.Dir(current)
		require.NotEqual(t, parent, current, "could not find repository root")
		current = parent
	}
}

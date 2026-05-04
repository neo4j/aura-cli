// Copyright (c) "Neo4j"
// Neo4j Sweden AB [http://neo4j.com]

//go:build !windows

package query

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

// TestHTTPS_Smoke is an env-gated wrapper around scripts/test-https.sh. It
// boots a real Neo4j 5 container with HTTPS enabled and verifies the
// `--insecure` flag end-to-end. Skipped by default so `go test ./...` is
// unaffected; opt in with NEO4J_HTTPS_TEST=1.
//
// Requires: docker, openssl, curl on PATH; ports 7473/7474 free.
//
// Build constraint: Unix-only (the script is bash and assumes POSIX tools).
func TestHTTPS_Smoke(t *testing.T) {
	if os.Getenv("NEO4J_HTTPS_TEST") != "1" {
		t.Skip("set NEO4J_HTTPS_TEST=1 to run; needs docker + scripts/test-https.sh")
	}

	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatalf("runtime.Caller(0) failed")
	}
	// thisFile = <repo>/neo4j-cli/query/query_https_smoke_test.go
	// repoRoot = parent of neo4j-cli/
	repoRoot := filepath.Join(filepath.Dir(thisFile), "..", "..")
	scriptPath := filepath.Join(repoRoot, "scripts", "test-https.sh")

	if _, err := os.Stat(scriptPath); err != nil {
		t.Fatalf("script not found at %s: %v", scriptPath, err)
	}

	cmd := exec.Command("bash", scriptPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("scripts/test-https.sh failed: %v", err)
	}
}

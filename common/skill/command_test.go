// Copyright (c) "Neo4j"
// Neo4j Sweden AB [http://neo4j.com]

package skill_test

import (
	"bytes"
	"encoding/json"
	"io/fs"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/neo4j/cli/common/clicfg"
	"github.com/neo4j/cli/common/skill"
	"github.com/neo4j/cli/test/utils/testfs"
)

const testSkillName = "neo4j-cli"

// fixtureBundle mirrors the shape used by installer_test.go.
func fixtureBundle() fs.FS {
	return fstest.MapFS{
		"SKILL.md": {Data: []byte(`---
name: neo4j-cli
description: test desc
version: {{VERSION}}
---

# neo4j-cli
body
`)},
		"references/aura.md": {Data: []byte("# aura\n")},
	}
}

// fixture is the per-test wiring: a `skill` cobra command, captured
// stdout/stderr buffers, and the in-memory afero.Fs carried by the config.
type fixture struct {
	cmd    *cobra.Command
	stdout *bytes.Buffer
	stderr *bytes.Buffer
	fs     afero.Fs
}

// newFixture wires up a cfg + skill command, populating agent detect dirs
// for every name in agentNames. `output` becomes the `aura.output` config
// value ("default", "json", or "table"). The test HOME is set to homeDir.
func newFixture(t *testing.T, homeDir, output string, agentNames ...string) *fixture {
	t.Helper()
	t.Setenv("HOME", homeDir)
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(homeDir, "xdg"))

	cfgJSON := `{"aura":{"output":"` + output + `"}}`
	memFs, err := testfs.GetTestFs(cfgJSON, "{}")
	require.NoError(t, err)

	for _, n := range agentNames {
		a := skill.FindAgent(n)
		require.NotNil(t, a, "fixture agent %q must exist", n)
		dp, ok := a.DetectPath()
		require.True(t, ok)
		require.NoError(t, memFs.MkdirAll(dp, 0755))
	}

	cfg := clicfg.NewConfig(memFs, "1.7.0")
	cmd := skill.NewCmd(cfg, fixtureBundle(), testSkillName)

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)

	return &fixture{cmd: cmd, stdout: stdout, stderr: stderr, fs: memFs}
}

// exec runs the skill command with the given argv split.
func (f *fixture) exec(t *testing.T, args ...string) error {
	t.Helper()
	f.cmd.SetArgs(args)
	return f.cmd.Execute()
}

// resetBuffers lets a single fixture run multiple commands.
func (f *fixture) resetBuffers() {
	f.stdout.Reset()
	f.stderr.Reset()
}

// =========================================================================
// install
// =========================================================================

func TestInstallCmd_Table(t *testing.T) {
	f := newFixture(t, "/home/alice", "default", "claude-code", "cursor")

	require.NoError(t, f.exec(t, "install"))

	out := f.stdout.String()
	// Both agents listed with the "installed" action.
	assert.Contains(t, out, "claude-code")
	assert.Contains(t, out, "cursor")
	assert.Contains(t, out, "installed")
	// Verify on-disk side-effect (sanity check the underlying installer ran).
	cc := skill.FindAgent("claude-code")
	sp, _ := cc.SkillsPath()
	exists, _ := afero.Exists(f.fs, filepath.Join(sp, testSkillName, "SKILL.md"))
	assert.True(t, exists)
}

func TestInstallCmd_JSON(t *testing.T) {
	f := newFixture(t, "/home/alice", "json", "claude-code")

	require.NoError(t, f.exec(t, "install"))

	var rows []map[string]any
	require.NoError(t, json.Unmarshal(f.stdout.Bytes(), &rows))
	require.Len(t, rows, 1)
	assert.Equal(t, "claude-code", rows[0]["agent"])
	assert.Equal(t, "Claude Code", rows[0]["display_name"])
	assert.Equal(t, "installed", rows[0]["action"])
	assert.Contains(t, rows[0]["skills_path"].(string), testSkillName)
}

func TestInstallCmd_SingleAgentArg(t *testing.T) {
	f := newFixture(t, "/home/alice", "default", "claude-code", "cursor")

	require.NoError(t, f.exec(t, "install", "Claude-Code")) // case-insensitive

	out := f.stdout.String()
	assert.Contains(t, out, "claude-code")
	assert.NotContains(t, out, "cursor", "single-agent install must not touch cursor")
}

func TestInstallCmd_UnknownAgent(t *testing.T) {
	f := newFixture(t, "/home/alice", "default", "claude-code")

	err := f.exec(t, "install", "vscode")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unknown agent")
	assert.Contains(t, err.Error(), "valid agents:")
}

func TestInstallCmd_NoAgentsDetected(t *testing.T) {
	f := newFixture(t, "/home/alice", "default") // no agents detected

	err := f.exec(t, "install")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no supported agents detected")
}

// =========================================================================
// remove
// =========================================================================

func TestRemoveCmd_Table(t *testing.T) {
	f := newFixture(t, "/home/alice", "default", "claude-code")

	require.NoError(t, f.exec(t, "install"))
	f.resetBuffers()

	require.NoError(t, f.exec(t, "remove"))
	out := f.stdout.String()
	assert.Contains(t, out, "claude-code")
	assert.Contains(t, out, "removed")

	// Verify on-disk: install dir gone.
	cc := skill.FindAgent("claude-code")
	sp, _ := cc.SkillsPath()
	exists, _ := afero.DirExists(f.fs, filepath.Join(sp, testSkillName))
	assert.False(t, exists)
}

func TestRemoveCmd_JSON(t *testing.T) {
	f := newFixture(t, "/home/alice", "json", "claude-code")
	require.NoError(t, f.exec(t, "install"))
	f.resetBuffers()

	require.NoError(t, f.exec(t, "remove"))
	var rows []map[string]any
	require.NoError(t, json.Unmarshal(f.stdout.Bytes(), &rows))
	require.Len(t, rows, 1)
	assert.Equal(t, "removed", rows[0]["action"])
}

func TestRemoveCmd_Idempotent(t *testing.T) {
	f := newFixture(t, "/home/alice", "default", "claude-code")
	// Never installed.
	require.NoError(t, f.exec(t, "remove", "claude-code"))
	f.resetBuffers()
	// Second run still succeeds.
	require.NoError(t, f.exec(t, "remove", "claude-code"))
}

func TestRemoveCmd_UnknownAgent(t *testing.T) {
	f := newFixture(t, "/home/alice", "default", "claude-code")
	err := f.exec(t, "remove", "vscode")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unknown agent")
}

// =========================================================================
// list
// =========================================================================

func TestListCmd_Table(t *testing.T) {
	f := newFixture(t, "/home/alice", "default", "claude-code")
	require.NoError(t, f.exec(t, "install", "claude-code"))
	f.resetBuffers()

	require.NoError(t, f.exec(t, "list"))
	out := f.stdout.String()
	lower := strings.ToLower(out)
	assert.Contains(t, out, "claude-code")
	// Header columns present (case-insensitive — go-pretty upper-cases).
	assert.Contains(t, lower, "detected")
	assert.Contains(t, lower, "installed")
	// claude-code shows installed=yes; another agent (e.g. cursor) shows no.
	assert.Contains(t, out, "yes")
	assert.Contains(t, out, "no")
	// Installed version threaded through.
	assert.Contains(t, out, "1.7.0")
}

func TestListCmd_JSON(t *testing.T) {
	f := newFixture(t, "/home/alice", "json", "claude-code")
	require.NoError(t, f.exec(t, "install", "claude-code"))
	f.resetBuffers()

	require.NoError(t, f.exec(t, "list"))
	var rows []map[string]any
	require.NoError(t, json.Unmarshal(f.stdout.Bytes(), &rows))
	// Catalog length — should include all 10 agents.
	assert.Len(t, rows, 10)

	// claude-code entry should have detected/installed = true and version.
	var cc map[string]any
	for _, r := range rows {
		if r["agent"] == "claude-code" {
			cc = r
			break
		}
	}
	require.NotNil(t, cc)
	assert.Equal(t, true, cc["detected"])
	assert.Equal(t, true, cc["installed"])
	assert.Equal(t, "1.7.0", cc["installed_version"])
}

// =========================================================================
// check
// =========================================================================

func TestCheckCmd_NoDrift(t *testing.T) {
	f := newFixture(t, "/home/alice", "default", "claude-code")
	require.NoError(t, f.exec(t, "install", "claude-code"))
	f.resetBuffers()

	require.NoError(t, f.exec(t, "check"))
	out := f.stdout.String()
	assert.Contains(t, out, "ok")
	assert.Contains(t, out, "1.7.0")
}

func TestCheckCmd_Drift(t *testing.T) {
	f := newFixture(t, "/home/alice", "default", "claude-code")
	require.NoError(t, f.exec(t, "install", "claude-code"))

	// Mutate the installed SKILL.md to a stale version.
	cc := skill.FindAgent("claude-code")
	sp, _ := cc.SkillsPath()
	skillFile := filepath.Join(sp, testSkillName, "SKILL.md")
	require.NoError(t, afero.WriteFile(f.fs, skillFile, []byte("---\nversion: 0.1.0\n---\n"), 0600))
	f.resetBuffers()

	err := f.exec(t, "check")
	require.Error(t, err, "check must exit non-zero on drift")
	assert.Contains(t, err.Error(), "drift")

	out := f.stdout.String()
	assert.Contains(t, out, "drift")
	assert.Contains(t, out, "0.1.0")
}

func TestCheckCmd_JSON(t *testing.T) {
	f := newFixture(t, "/home/alice", "json", "claude-code")
	require.NoError(t, f.exec(t, "install", "claude-code"))
	f.resetBuffers()

	require.NoError(t, f.exec(t, "check"))
	var rows []map[string]any
	require.NoError(t, json.Unmarshal(f.stdout.Bytes(), &rows))
	require.Len(t, rows, 1)
	assert.Equal(t, "ok", rows[0]["status"])
	assert.Equal(t, "1.7.0", rows[0]["installed_version"])
	assert.Equal(t, "1.7.0", rows[0]["current_version"])
}

func TestCheckCmd_NoneInstalled(t *testing.T) {
	f := newFixture(t, "/home/alice", "default", "claude-code")
	// Never installed — Check returns no rows, drift=false, exits 0.
	require.NoError(t, f.exec(t, "check"))
	out := f.stdout.String()
	assert.Contains(t, strings.ToLower(out), "no installed skills")
}

// =========================================================================
// invalid output flag
// =========================================================================

func TestInvalidOutputFlag(t *testing.T) {
	f := newFixture(t, "/home/alice", "default", "claude-code")
	err := f.exec(t, "list", "--output", "yaml")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid output value")
}

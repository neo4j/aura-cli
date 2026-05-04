// Copyright (c) "Neo4j"
// Neo4j Sweden AB [http://neo4j.com]

package skill

import (
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAGENTSCatalog(t *testing.T) {
	// Catalog must match the Rust reference (10 agents, in this order).
	// Locking order means stable `skill list` output across releases.
	expected := []string{
		"claude-code", "cursor", "windsurf", "copilot", "gemini-cli",
		"cline", "codex", "pi", "opencode", "junie",
	}
	require.Len(t, AGENTS, len(expected))
	for i, want := range expected {
		assert.Equal(t, want, AGENTS[i].Name, "agent at index %d", i)
		assert.NotEmpty(t, AGENTS[i].DisplayName)
		assert.NotEmpty(t, AGENTS[i].DetectDir)
		assert.NotEmpty(t, AGENTS[i].SkillsDir)
	}
}

func TestExpandPath(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		home     string
		xdg      string
		xdgUnset bool
		want     string
		wantOk   bool
	}{
		{
			name:   "bare tilde expands to HOME",
			path:   "~",
			home:   "/home/alice",
			want:   "/home/alice",
			wantOk: true,
		},
		{
			name:   "tilde slash prefix expands to HOME/sub",
			path:   "~/.claude",
			home:   "/home/alice",
			want:   filepath.Join("/home/alice", ".claude"),
			wantOk: true,
		},
		{
			name:   "tilde slash deeper path",
			path:   "~/.codeium/windsurf",
			home:   "/home/alice",
			want:   filepath.Join("/home/alice", ".codeium/windsurf"),
			wantOk: true,
		},
		{
			name:   "bare tilde with no HOME returns ok=false",
			path:   "~",
			home:   "",
			wantOk: false,
		},
		{
			name:   "tilde slash with no HOME returns ok=false",
			path:   "~/.claude",
			home:   "",
			wantOk: false,
		},
		{
			name:   "XDG set substitutes verbatim",
			path:   "$XDG_CONFIG_HOME/opencode",
			home:   "/home/alice",
			xdg:    "/home/alice/xdg",
			want:   "/home/alice/xdg/opencode",
			wantOk: true,
		},
		{
			name:     "XDG unset falls back to HOME/.config",
			path:     "$XDG_CONFIG_HOME/opencode",
			home:     "/home/alice",
			xdgUnset: true,
			want:     filepath.Join("/home/alice", ".config", "opencode"),
			wantOk:   true,
		},
		{
			name:     "XDG empty string falls back to HOME/.config",
			path:     "$XDG_CONFIG_HOME/opencode",
			home:     "/home/alice",
			xdg:      "",
			xdgUnset: false,
			want:     filepath.Join("/home/alice", ".config", "opencode"),
			wantOk:   true,
		},
		{
			name:     "XDG fallback with no HOME returns ok=false",
			path:     "$XDG_CONFIG_HOME/opencode",
			home:     "",
			xdgUnset: true,
			wantOk:   false,
		},
		{
			name:   "absolute path unchanged",
			path:   "/etc/agent/skills",
			home:   "/home/alice",
			want:   "/etc/agent/skills",
			wantOk: true,
		},
		{
			name:   "relative path unchanged",
			path:   "relative/path",
			home:   "/home/alice",
			want:   "relative/path",
			wantOk: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("HOME", tc.home)
			// Setting XDG_CONFIG_HOME="" is equivalent to unset for our
			// purposes — expandPath treats empty the same as missing.
			if tc.xdgUnset {
				t.Setenv("XDG_CONFIG_HOME", "")
			} else {
				t.Setenv("XDG_CONFIG_HOME", tc.xdg)
			}

			got, ok := expandPath(tc.path)
			assert.Equal(t, tc.wantOk, ok)
			if tc.wantOk {
				assert.Equal(t, tc.want, got)
			}
		})
	}
}

func TestFindAgent(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantName string // empty means nil
	}{
		{name: "exact lowercase match", input: "claude-code", wantName: "claude-code"},
		{name: "mixed case match", input: "Claude-Code", wantName: "claude-code"},
		{name: "uppercase match", input: "CURSOR", wantName: "cursor"},
		{name: "unknown returns nil", input: "vscode", wantName: ""},
		{name: "empty returns nil", input: "", wantName: ""},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := FindAgent(tc.input)
			if tc.wantName == "" {
				assert.Nil(t, got)
				return
			}
			require.NotNil(t, got)
			assert.Equal(t, tc.wantName, got.Name)
		})
	}
}

func TestDetectAgents(t *testing.T) {
	t.Setenv("HOME", "/home/alice")
	t.Setenv("XDG_CONFIG_HOME", "/home/alice/xdg")

	fs := afero.NewMemMapFs()
	// Fake two agent install markers.
	require.NoError(t, fs.MkdirAll("/home/alice/.claude", 0755))
	require.NoError(t, fs.MkdirAll("/home/alice/xdg/opencode", 0755))

	got := DetectAgents(fs)
	require.Len(t, got, 2)
	assert.Equal(t, "claude-code", got[0].Name)
	assert.Equal(t, "opencode", got[1].Name)
}

func TestDetectAgentsEmpty(t *testing.T) {
	t.Setenv("HOME", "/home/alice")
	t.Setenv("XDG_CONFIG_HOME", "/home/alice/xdg")

	fs := afero.NewMemMapFs()
	got := DetectAgents(fs)
	assert.Empty(t, got)
}

func TestDetectAgentsHomeUnset(t *testing.T) {
	t.Setenv("HOME", "")
	t.Setenv("XDG_CONFIG_HOME", "")

	fs := afero.NewMemMapFs()
	got := DetectAgents(fs)
	// All AGENTS reference HOME (directly or via XDG fallback) — unresolvable
	// paths are skipped, so the result is empty.
	assert.Empty(t, got)
}

func TestDetectAgentsIgnoresFile(t *testing.T) {
	// DetectDir must be a directory, not a file. afero.DirExists returns
	// false for files, so a file at the marker path doesn't count as
	// detected.
	t.Setenv("HOME", "/home/alice")
	fs := afero.NewMemMapFs()
	require.NoError(t, afero.WriteFile(fs, "/home/alice/.claude", []byte("not a dir"), 0644))

	got := DetectAgents(fs)
	assert.Empty(t, got)
}

func TestAgentDetectAndSkillsPath(t *testing.T) {
	t.Setenv("HOME", "/home/alice")
	t.Setenv("XDG_CONFIG_HOME", "/home/alice/xdg")

	a := FindAgent("claude-code")
	require.NotNil(t, a)

	dp, ok := a.DetectPath()
	require.True(t, ok)
	assert.Equal(t, filepath.Join("/home/alice", ".claude"), dp)

	sp, ok := a.SkillsPath()
	require.True(t, ok)
	assert.Equal(t, filepath.Join("/home/alice", ".claude/skills"), sp)
}

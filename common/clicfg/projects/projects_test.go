// Copyright (c) "Neo4j"
// Neo4j Sweden AB [http://neo4j.com]

package projects

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestProjects(t *testing.T) *AuraConfigProjects {
	t.Helper()
	fs := afero.NewMemMapFs()
	require.NoError(t, afero.WriteFile(fs, "/test", []byte(`{}`), 0o600))
	return NewAuraConfigProjects(fs, "/test")
}

func TestRemoveDefaultProjectClearsDefault(t *testing.T) {
	p := newTestProjects(t)

	require.NoError(t, p.Add("prod", "o1", "p1"))
	require.NoError(t, p.Add("staging", "o1", "p2"))
	require.NoError(t, p.Add("dev", "o1", "p3"))
	_, err := p.SetDefault("prod")
	require.NoError(t, err)

	require.NoError(t, p.Remove("prod"))

	def, err := p.Default()
	require.NoError(t, err)
	assert.Equal(t, &AuraProject{}, def, "expected no default project after removing the default")
}

func TestRemoveNonDefaultProjectLeavesDefaultUnchanged(t *testing.T) {
	p := newTestProjects(t)

	require.NoError(t, p.Add("prod", "o1", "p1"))
	require.NoError(t, p.Add("staging", "o1", "p2"))
	_, err := p.SetDefault("prod")
	require.NoError(t, err)

	require.NoError(t, p.Remove("staging"))

	def, err := p.Default()
	require.NoError(t, err)
	assert.Equal(t, &AuraProject{OrganizationId: "o1", ProjectId: "p1"}, def)
}

func TestRemoveLastProjectClearsDefault(t *testing.T) {
	p := newTestProjects(t)

	require.NoError(t, p.Add("prod", "o1", "p1"))
	require.NoError(t, p.Remove("prod"))

	def, err := p.Default()
	require.NoError(t, err)
	assert.Equal(t, &AuraProject{}, def)
}

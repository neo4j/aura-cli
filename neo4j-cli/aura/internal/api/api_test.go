// Copyright (c) "Neo4j"
// Neo4j Sweden AB [http://neo4j.com]

package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidatePathRejectsDotDotSegments(t *testing.T) {
	cases := []string{
		"/tenants/../../../../evil",
		"../evil",
		"/tenants/..",
		"tenants/../evil",
	}

	for _, path := range cases {
		err := validatePath(path)
		assert.NotNil(t, err, "expected path %q to be rejected", path)
	}
}

func TestValidatePathAllowsOrdinaryPaths(t *testing.T) {
	cases := []string{
		"/tenants/6981ace7-efe8-4f5c-b7c5-267b5162ce91",
		"/tenants/6981ace7-efe8-4f5c-b7c5-267b5162ce91/metrics-integration",
		"/instances",
		"/instances/some..id",
	}

	for _, path := range cases {
		err := validatePath(path)
		assert.Nil(t, err, "expected path %q to be allowed", path)
	}
}

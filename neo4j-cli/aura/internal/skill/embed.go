// Copyright (c) "Neo4j"
// Neo4j Sweden AB [http://neo4j.com]

// Package skill embeds the generated aura-cli agent-skill bundle so the
// `aura-cli skill install` command can write it onto the host filesystem.
//
// Regenerate the bundle with `go generate ./neo4j-cli/aura/internal/skill/...`
// (or `make generate`). The committed bundle/ tree is the source of truth
// shipped in the binary; CI runs `make generate-check` to fail on drift.
//
// Edit description.txt and additions.md (sibling files) to change the
// frontmatter description and the SKILL.md "Gotchas" section. The
// remainder of SKILL.md and every references/<sub>.md is derived from
// the cobra tree built by aura.NewStandaloneCmd.
package skill

import "embed"

//go:generate go run ./gen

//go:embed bundle
var Bundle embed.FS

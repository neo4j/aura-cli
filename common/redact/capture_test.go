// Copyright (c) "Neo4j"
// Neo4j Sweden AB [http://neo4j.com]

package redact

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

// commandTreeWithClashingShorthand registers "-p" on two subcommands for
// different long flags: one safe ("project-id"), one not ("password").
func commandTreeWithClashingShorthand() *cobra.Command {
	root := &cobra.Command{Use: "neo4j-cli"}

	safeCmd := &cobra.Command{Use: "list", Run: func(cmd *cobra.Command, args []string) {}}
	safeCmd.Flags().StringP("project-id", "p", "", "")

	unsafeCmd := &cobra.Command{Use: "login", Run: func(cmd *cobra.Command, args []string) {}}
	unsafeCmd.Flags().StringP("password", "p", "", "")

	root.AddCommand(safeCmd, unsafeCmd)
	return root
}

func TestCaptureArgsResolvesShorthandPerCommand(t *testing.T) {
	root := commandTreeWithClashingShorthand()

	CaptureArgs(root, []string{"list", "-p", "project-12345"})
	assert.Equal(t, []string{"list", "-p", "project-12345"}, CapturedArgs(),
		"-p resolves to the safe project-id flag on the list command")

	CaptureArgs(root, []string{"login", "-p", "super-secret"})
	assert.Equal(t, []string{"login", "-p", mask}, CapturedArgs(),
		"the same -p resolves to the unsafe password flag on the login command, so it's still masked")
}

func TestCaptureArgsWithNoMatchingCommandDefaultsToMasked(t *testing.T) {
	root := commandTreeWithClashingShorthand()

	CaptureArgs(root, []string{"not-a-real-command", "-p", "some-value"})
	assert.Equal(t, []string{"not-a-real-command", "-p", mask}, CapturedArgs(),
		"args that don't resolve to a known command fall back to masking")
}

func TestCaptureArgsWithNilRootCmdDefaultsToMasked(t *testing.T) {
	CaptureArgs(nil, []string{"-p", "some-value"})
	assert.Equal(t, []string{"-p", mask}, CapturedArgs())
}

func TestCaptureArgsResolvesInheritedPersistentShorthand(t *testing.T) {
	root := &cobra.Command{Use: "neo4j-cli"}
	root.PersistentFlags().StringP("organization-id", "o", "", "")

	child := &cobra.Command{Use: "list", Run: func(cmd *cobra.Command, args []string) {}}
	root.AddCommand(child)

	CaptureArgs(root, []string{"list", "-o", "org-12345"})
	assert.Equal(t, []string{"list", "-o", "org-12345"}, CapturedArgs(),
		"a shorthand registered as a persistent flag on a parent command resolves on the child")
}

// Copyright (c) "Neo4j"
// Neo4j Sweden AB [http://neo4j.com]

package redact

import "github.com/spf13/cobra"

var capturedArgs []string

func CaptureArgs(rootCmd *cobra.Command, args []string) {
	capturedArgs = MaskArgsWithShorthandResolver(args, shorthandResolverFor(rootCmd, args))
}

func CapturedArgs() []string {
	return capturedArgs
}

func shorthandResolverFor(rootCmd *cobra.Command, args []string) ShorthandResolver {
	if rootCmd == nil {
		return nil
	}

	cmd, _, err := rootCmd.Find(args)
	if err != nil || cmd == nil {
		return nil
	}

	return func(flagName string) (string, bool) {
		if len(flagName) != 1 {
			return "", false
		}

		flag := cmd.Flags().ShorthandLookup(flagName)
		if flag == nil {
			return "", false
		}
		return flag.Name, true
	}
}

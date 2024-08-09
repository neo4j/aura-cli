package instance

import (
	"fmt"

	"github.com/neo4j/cli/neo4j/aura/internal/api"
	"github.com/spf13/cobra"
)

func NewUpdateCmd() *cobra.Command {
	var (
		memory string
		name   string
	)

	const (
		memoryFlag = "memory"
		nameFlag   = "name"
	)

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Updates an instance",
		Long: `This command allows you to rename and/or resize an Aura instance.

Resizing an instance is an asynchronous operation. The instance remains available throughout.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			body := map[string]any{}

			if memory != "" {
				body["memory"] = memory
			}

			if name != "" {
				body["name"] = name
			}

			return api.MakeRequest(cmd, "PATCH", fmt.Sprintf("/instances/%s", args[0]), body)
		},
	}

	cmd.Flags().StringVar(&memory, memoryFlag, "", "The size of the instance memory in GB.")

	cmd.Flags().StringVar(&name, nameFlag, "", "The name of the instance (any UTF-8 characters with no trailing or leading whitespace).")

	cmd.MarkFlagsOneRequired(memoryFlag, nameFlag)

	return cmd
}

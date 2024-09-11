package instance

import (
	"fmt"
	"net/http"

	"github.com/neo4j/cli/common/clicfg"
	"github.com/neo4j/cli/neo4j/aura/internal/api"
	"github.com/neo4j/cli/neo4j/aura/internal/output"
	"github.com/spf13/cobra"
)

func NewUpdateCmd(cfg *clicfg.Config) *cobra.Command {
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

			path := fmt.Sprintf("/instances/%s", args[0])

			cmd.SilenceUsage = true
			resBody, statusCode, err := api.MakeRequest(cfg, http.MethodPatch, path, body)
			if err != nil {
				return err
			}

			if statusCode == http.StatusAccepted || statusCode == http.StatusOK {
				err = output.PrintBody(cmd, cfg, resBody)
				if err != nil {
					return err
				}

			}
			return nil
		},
	}

	cmd.Flags().StringVar(&memory, memoryFlag, "", "The size of the instance memory in GB.")

	cmd.Flags().StringVar(&name, nameFlag, "", "The name of the instance (any UTF-8 characters with no trailing or leading whitespace).")

	cmd.MarkFlagsOneRequired(memoryFlag, nameFlag)

	return cmd
}

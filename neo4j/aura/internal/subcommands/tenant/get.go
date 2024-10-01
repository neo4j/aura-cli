package tenant

import (
	"fmt"
	"net/http"

	"github.com/neo4j/cli/common/clicfg"
	"github.com/neo4j/cli/neo4j/aura/internal/api"
	"github.com/neo4j/cli/neo4j/aura/internal/output"
	"github.com/spf13/cobra"
)

func NewGetCmd(cfg *clicfg.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Returns tenant details",
		Long:  "This subcommand returns details about a specific Aura Tenant.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := fmt.Sprintf("/tenants/%s", args[0])

			cmd.SilenceUsage = true
			resBody, statusCode, err := api.MakeRequest(cfg, http.MethodGet, path, nil)
			if err != nil {
				return err
			}

			if statusCode == http.StatusOK {
				err = output.PrintBody(cmd, cfg, resBody, []string{"id", "name"})
				if err != nil {
					return err
				}

			}

			return nil
		},
	}
}

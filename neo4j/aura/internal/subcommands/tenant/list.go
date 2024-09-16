package tenant

import (
	"net/http"

	"github.com/neo4j/cli/common/clicfg"
	"github.com/neo4j/cli/neo4j/aura/internal/api"
	"github.com/neo4j/cli/neo4j/aura/internal/output"
	"github.com/spf13/cobra"
)

func NewListCmd(cfg *clicfg.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "Returns a list of tenants",
		Long:  "This subcommand returns a list containing a summary of each of your Aura Tenants. To find out more about a specific Tenant, retrieve the details using the get subcommand.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			resBody, statusCode, err := api.MakeRequest(cfg, "/tenants", &api.RequestConfig{
				Method: http.MethodGet,
			})
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

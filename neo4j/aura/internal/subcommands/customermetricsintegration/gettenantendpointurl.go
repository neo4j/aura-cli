package customermetricsintegration

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"

	"github.com/neo4j/cli/common/clicfg"
	"github.com/neo4j/cli/neo4j/aura/internal/api"
	"github.com/neo4j/cli/neo4j/aura/internal/output"
)

func NewGetTenantEndpointUrlCmd(cfg *clicfg.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "get-tenant-endpoint-url",
		Short: "Returns tenant metric endpoint URL",
		Long:  "This subcommand returns the Prometheus metric endpoint URL for the specified tenant.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := fmt.Sprintf("/tenants/%s/metrics-integration", args[0])
			cmd.SilenceUsage = true
			resBody, statusCode, err := api.MakeRequest(cfg, http.MethodGet, path, nil)
			if err != nil {
				return err
			}

			if statusCode == http.StatusOK {
				fmt.Println(string(resBody))
				err = output.PrintBody(cmd, cfg, resBody, []string{"endpoint"})
				if err != nil {
					return err
				}

			}

			return nil
		},
	}
}

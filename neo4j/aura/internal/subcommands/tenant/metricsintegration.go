package tenant

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"

	"github.com/neo4j/cli/common/clicfg"
	"github.com/neo4j/cli/neo4j/aura/internal/api"
	"github.com/neo4j/cli/neo4j/aura/internal/output"
)

func NewMetricsIntegrationCmd(cfg *clicfg.Config) *cobra.Command {
	var tenantId string

	const tenantIdFlag = "tenant-id"

	cmd := &cobra.Command{
		Use:   "metrics-integration",
		Short: "Returns tenant metric endpoint URL",
		Long:  "This subcommand returns the Prometheus metric endpoint URL for the specified tenant.",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if cfg.Aura.DefaultTenant() == "" {
				cmd.MarkFlagRequired(tenantIdFlag)
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if tenantId == "" {
				tenantId = cfg.Aura.DefaultTenant()
			}

			path := fmt.Sprintf("/tenants/%s/metrics-integration", tenantId)
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

	cmd.Flags().StringVar(&tenantId, tenantIdFlag, "", "")

	return cmd
}

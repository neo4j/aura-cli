package instance

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
		Use:   "get",
		Short: "Returns instance details",
		Long:  "This endpoint returns details about a specific Aura Instance.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := fmt.Sprintf("/instances/%s", args[0])

			cmd.SilenceUsage = true
			resBody, statusCode, err := api.MakeRequest(cfg, http.MethodGet, path, nil)
			if err != nil {
				return err
			}

			if statusCode == http.StatusOK {
				err = output.PrintBody(cmd, cfg, resBody, []string{"id", "name", "tenant_id", "status", "connection_url", "cloud_provider", "region", "type", "memory", "storage", "customer_managed_key_id"})
				if err != nil {
					return err
				}

			}

			return nil
		},
	}
}

package customermanagedkey

import (
	"fmt"
	"net/http"

	"github.com/neo4j/cli/common/clicfg"
	"github.com/neo4j/cli/neo4j/aura/internal/api"
	"github.com/neo4j/cli/neo4j/aura/internal/output"
	"github.com/spf13/cobra"
)

func NewListCmd(cfg *clicfg.Config) *cobra.Command {
	var tenantId string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "Returns a list of customer managed keys",
		Long: `This subcommand returns a list containing a summary of each of your customer managed keys. To find out more about a specific key, retrieve the details using the get subcommand.

You can filter keys in a particular tenant using --tenant-id. If the tenant flag is not specified, this endpoint lists all keys a user has access to across all tenants.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			var path string

			if tenantId != "" {
				path = fmt.Sprintf("/customer-managed-keys?tenantId=%s", tenantId)
			} else {
				path = "/customer-managed-keys"
			}

			cmd.SilenceUsage = true
			resBody, statusCode, err := api.MakeRequest(cfg, http.MethodGet, path, nil)
			if err != nil {
				return err
			}

			if statusCode == http.StatusOK {
				output.PrintBody(cmd, cfg, resBody)

			}

			return nil
		},
	}

	cmd.Flags().StringVar(&tenantId, "tenant-id", "", "An optional Tenant ID to filter customer managed keys in a tenant")

	return cmd
}

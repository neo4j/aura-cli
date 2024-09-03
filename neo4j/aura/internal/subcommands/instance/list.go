package instance

import (
	"fmt"
	"net/http"

	"github.com/neo4j/cli/neo4j/aura/internal/api"
	"github.com/neo4j/cli/neo4j/aura/internal/output"
	"github.com/spf13/cobra"
)

func NewListCmd() *cobra.Command {
	var tenantId string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "Returns a list of instances",
		Long: `This subcommand returns a list containing a summary of each of your Aura instances. To find out more about a specific instance, retrieve the details using the get subcommand.

You can filter instances in a particular tenant using --tenant-id. If the tenant flag is not specified, this subcommand lists all instances a user has access to across all tenants.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var path string

			if tenantId != "" {
				path = fmt.Sprintf("/instances?tenantId=%s", tenantId)
			} else {
				path = "/instances"
			}

			resBody, statusCode, err := api.MakeRequest(cmd, http.MethodGet, path, nil)
			if err != nil {
				return err
			}

			if statusCode == http.StatusOK {
				err = output.PrintBody2(cmd, resBody, []string{"id", "name", "tenant_id", "cloud_provider"})
				if err != nil {
					return err
				}
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&tenantId, "tenant-id", "", "An optional Tenant ID to filter instances in a tenant")

	return cmd
}

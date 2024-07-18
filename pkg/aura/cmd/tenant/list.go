package tenant

import (
	"github.com/neo4j/cli/pkg/aura/api"
	"github.com/spf13/cobra"
)

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "Returns a list of tenants",
	Long:  "This subcommand returns a list containing a summary of each of your Aura Tenants. To find out more about a specific Tenant, retrieve the details using the get subcommand.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return api.MakeRequest(cmd, "GET", "/tenants", nil)
	},
}

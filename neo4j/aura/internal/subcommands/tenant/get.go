package tenant

import (
	"fmt"

	"github.com/neo4j/cli/neo4j/aura/internal/api"
	"github.com/spf13/cobra"
)

func NewGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get",
		Short: "Returns tenant details",
		Long:  "This subcommand returns details about a specific Aura Tenant.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return api.MakeRequest(cmd, "GET", fmt.Sprintf("/tenants/%s", args[0]), nil)
		},
	}
}

package instance

import (
	"fmt"

	"github.com/neo4j/cli/pkg/aura/api"
	"github.com/spf13/cobra"
)

func NewGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get",
		Short: "Returns instance details",
		Long:  "This endpoint returns details about a specific Aura Instance.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return api.MakeRequest(cmd, "GET", fmt.Sprintf("/instances/%s", args[0]), nil)
		},
	}
}

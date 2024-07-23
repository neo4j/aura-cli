package customermanagedkey

import (
	"fmt"

	"github.com/neo4j/cli/pkg/aura/api"
	"github.com/spf13/cobra"
)

func NewDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete",
		Short: "Deletes a customer managed key",
		Long: `Deletes a Customer Managed Key from Aura.

Note that you can only delete a Key if it is not being used by any instances, otherwise you will get an error with the reason field set to encryption-key-is-active.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return api.MakeRequest(cmd, "DELETE", fmt.Sprintf("/customer-managed-keys/%s", args[0]), nil)
		},
	}
}

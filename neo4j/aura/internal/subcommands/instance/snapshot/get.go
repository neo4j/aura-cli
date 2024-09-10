package snapshot

import (
	"fmt"
	"net/http"

	"github.com/neo4j/cli/neo4j/aura/internal/api"
	"github.com/neo4j/cli/neo4j/aura/internal/output"
	"github.com/spf13/cobra"
)

func NewGetCmd() *cobra.Command {
	var instanceId string
	// var date string

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Returns details of a snapshot",
		Long:  `This subcommand returns a list containing a summary of each snapshot of an Aura instance. To find out more about a specific snapshot, retrieve the details using the get subcommand.`,
		RunE: func(cmd *cobra.Command, args []string) error {

			path := fmt.Sprintf("/instances/%s/snapshots/%s", instanceId, args[0])

			resBody, statusCode, err := api.MakeRequest(cmd, http.MethodGet, path, nil)
			if err != nil {
				return err
			}

			if statusCode == http.StatusOK {
				err = output.PrintBody(cmd, resBody, []string{"snapshot_id", "instance_id", "profile", "status", "timestamp", "exportable"})
				if err != nil {
					return err
				}
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&instanceId, "instance-id", "", "The id of the instance to list its snapshots")
	cmd.MarkFlagRequired("instance-id")
	// cmd.Flags().StringVar(&tenantId, "tenant-id", "", "An optional Tenant ID to filter instances in a tenant")

	return cmd
}

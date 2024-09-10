package snapshot

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/neo4j/cli/neo4j/aura/internal/api"
	"github.com/neo4j/cli/neo4j/aura/internal/output"
	"github.com/spf13/cobra"
)

func NewCreateCmd() *cobra.Command {
	var instanceId string
	var await bool
	// var date string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a snapshot",
		// Long:  `This subcommand returns a list containing a summary of each snapshot of an Aura instance. To find out more about a specific snapshot, retrieve the details using the get subcommand.`,
		RunE: func(cmd *cobra.Command, args []string) error {

			path := fmt.Sprintf("/instances/%s/snapshots", instanceId)

			resBody, statusCode, err := api.MakeRequest(cmd, http.MethodPost, path, nil)
			if err != nil {
				return err
			}

			if statusCode == http.StatusAccepted {
				err = output.PrintBody(cmd, resBody, []string{"snapshot_id"})
				if err != nil {
					return err
				}

				if await {
					cmd.Println("Waiting for instance to be ready...")
					var response api.CreateSnapshotResponse
					if err := json.Unmarshal(resBody, &response); err != nil {
						return err
					}

					// Snapshot is not ready after pending
					pollResponse, err := api.PollSnapshot(cmd, instanceId, response.Data.SnapshotId, api.SnapshotStatusPending)
					if err != nil {
						return err
					}

					cmd.Println("Instance Status:", pollResponse.Data.Status)
				}
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&instanceId, "instance-id", "", "The id of the instance to list its snapshots")
	cmd.MarkFlagRequired("instance-id")
	cmd.Flags().BoolVar(&await, "await", false, "Waits until created snapshot is ready.")

	// cmd.Flags().StringVar(&tenantId, "tenant-id", "", "An optional Tenant ID to filter instances in a tenant")

	return cmd
}

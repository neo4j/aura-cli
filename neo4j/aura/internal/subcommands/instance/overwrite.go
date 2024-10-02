package instance

import (
	"fmt"
	"net/http"

	"github.com/neo4j/cli/common/clicfg"
	"github.com/neo4j/cli/neo4j/aura/internal/api"
	"github.com/neo4j/cli/neo4j/aura/internal/output"
	"github.com/spf13/cobra"
)

func NewOverwriteCmd(cfg *clicfg.Config) *cobra.Command {
	var sourceInstanceId string
	var sourceSnapshotId string
	var await bool

	cmd := &cobra.Command{
		Use:   "overwrite",
		Short: "Starts the process of overwriting the specified instance with data from the source instance provided",
		Long: `Starts the process of overwriting the specified instance with data from the source instance provided.
		The overwrite process mimics the 'Clone to existing' functionality of the Aura Console.

If only a source_instance_id is provided, a new snapshot of that instance is created and used for overwriting. Alternatively, you can specify an additional source_snapshot_id to use a specific snapshot of the source instance for overwriting.
		`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			instanceId := args[0]
			path := fmt.Sprintf("/instances/%s/overwrite", instanceId)

			cmd.SilenceUsage = true

			postBody := make(map[string]any)
			if sourceInstanceId == "" {
				sourceInstanceId = instanceId
			}

			postBody["source_instance_id"] = sourceInstanceId
			if sourceSnapshotId != "" {
				postBody["source_snapshot_id"] = sourceSnapshotId

			}

			resBody, statusCode, err := api.MakeRequest(cfg, path, &api.RequestConfig{
				Method:   http.MethodPost,
				PostBody: postBody,
			})
			if err != nil {
				return err
			}

			if statusCode == http.StatusAccepted {
				err = output.PrintBody(cmd, cfg, resBody, []string{"id", "name", "tenant_id", "status", "connection_url", "cloud_provider", "region", "type", "memory", "storage", "customer_managed_key_id"})
				if err != nil {
					return err
				}
			}

			if await {
				cmd.Println("Waiting for instance to be ready...")
				pollResponse, err := api.PollInstance(cfg, instanceId, api.InstanceStatusOverwriting)
				if err != nil {
					return err
				}

				cmd.Println("Instance Status:", pollResponse.Data.Status)
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&sourceInstanceId, "source-instance-id", "", "The id of the source instance")
	cmd.Flags().StringVar(&sourceSnapshotId, "source-snapshot-id", "", "The id of the snapshot to use")
	cmd.Flags().BoolVar(&await, "await", false, "Waits until created snapshot is ready.")

	return cmd
}

package snapshot

import (
	"fmt"
	"net/http"

	"github.com/neo4j/cli/common/clicfg"
	"github.com/neo4j/cli/neo4j/aura/internal/api"
	"github.com/neo4j/cli/neo4j/aura/internal/output"
	"github.com/spf13/cobra"
)

func NewGetCmd(cfg *clicfg.Config) *cobra.Command {
	var instanceId string
	var snapshotId string

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get details of a snapshot",
		Long:  `This endpoint returns details about a specific snapshot.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println(len(args))

			cmd.SilenceUsage = true
			path := fmt.Sprintf("/instances/%s/snapshots/%s", instanceId, snapshotId)

			resBody, statusCode, err := api.MakeRequest(cfg, path, &api.RequestConfig{
				Method: http.MethodGet,
			})
			if err != nil {
				return err
			}

			if statusCode == http.StatusOK {
				err = output.PrintBody(cmd, cfg, resBody, []string{"snapshot_id", "instance_id", "profile", "status", "timestamp", "exportable"})
				if err != nil {
					return err
				}
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&instanceId, "instance-id", "", "The id of the instance with the snapshot")
	cmd.Flags().StringVar(&snapshotId, "snapshot-id", "", "The id of the snaphost")
	cmd.MarkFlagRequired("instance-id")
	cmd.MarkFlagRequired("snapshot-id")

	return cmd
}

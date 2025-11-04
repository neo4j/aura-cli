package serverdatabase

import (
	"fmt"
	"log"
	"net/http"

	"github.com/neo4j/cli/common/clicfg"
	"github.com/neo4j/cli/neo4j-cli/aura/internal/api"
	"github.com/neo4j/cli/neo4j-cli/aura/internal/output"
	"github.com/spf13/cobra"
)

func NewListCmd(cfg *clicfg.Config) *cobra.Command {
	var (
		organizationId string
		projectId      string
		deploymentId   string
		serverId       string
	)

	const (
		organizationIdFlag = "organization-id"
		projectIdFlag      = "project-id"
		deploymentIdFlag   = "deployment-id"
		serverIdFlag       = "server-id"
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "Returns all deployments",
		Long:  "This endpoint returns all Fleet Manager deployments for the given project.",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := fmt.Sprintf("/organizations/%s/projects/%s/deployments/%s/servers/%s/databases", organizationId, projectId, deploymentId, serverId)

			cmd.SilenceUsage = true
			resBody, statusCode, err := api.MakeRequest(cfg, path, &api.RequestConfig{
				Method:  http.MethodGet,
				Version: api.AuraApiVersion2,
			})
			if err != nil {
				return err
			}

			if statusCode == http.StatusOK {
				fields := []string{
					"id",
					"deployment_id",
					"server_id",
					"name",
					"type",
					"current_status",
					"status_message",
					"graph_shards",
					"last_committed_txn",
					"last_seen",
					"property_shards",
					"replication_lag",
					"role",
					"writer",
				}
				output.PrintBody(cmd, cfg, resBody, fields)
			}

			return nil
		},
	}
	cmd.Flags().StringVar(&organizationId, organizationIdFlag, "", "(required) Organization ID")
	cmd.Flags().StringVar(&projectId, projectIdFlag, "", "(required) Project/tenant ID")
	cmd.Flags().StringVar(&deploymentId, deploymentIdFlag, "", "(required) Deployment ID")
	cmd.Flags().StringVar(&serverId, serverIdFlag, "", "(required) Server ID")

	err := cmd.MarkFlagRequired(organizationIdFlag)
	if err != nil {
		log.Fatal(err)
	}
	err = cmd.MarkFlagRequired(projectIdFlag)
	if err != nil {
		log.Fatal(err)
	}
	err = cmd.MarkFlagRequired(deploymentIdFlag)
	if err != nil {
		log.Fatal(err)
	}
	err = cmd.MarkFlagRequired(serverIdFlag)
	if err != nil {
		log.Fatal(err)
	}

	return cmd
}

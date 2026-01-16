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
		Short: "Returns deployment server databases.",
		Long:  "Returns databases for the given Fleet Manager deployment server.",
		Args:  cobra.ExactArgs(0),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			defaultSetting, err := cfg.Settings.Aura.GetDefault()
			if err != nil {
				log.Fatal(err)
			}
			if defaultSetting.OrganizationId == "" {
				err := cmd.MarkFlagRequired(organizationIdFlag)
				if err != nil {
					log.Fatal(err)
				}
			}

			if defaultSetting.ProjectId == "" {
				err := cmd.MarkFlagRequired(projectIdFlag)
				if err != nil {
					log.Fatal(err)
				}
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			defaultSetting, err := cfg.Settings.Aura.GetDefault()
			if err != nil {
				log.Fatal(err)
			}
			if organizationId == "" {
				organizationId = defaultSetting.OrganizationId
			}
			if projectId == "" {
				projectId = defaultSetting.ProjectId
			}
			path := fmt.Sprintf("/organizations/%s/projects/%s/fleet-manager/deployments/%s/servers/%s/databases", organizationId, projectId, deploymentId, serverId)

			cmd.SilenceUsage = true
			resBody, statusCode, err := api.MakeRequest(cfg, path, &api.RequestConfig{
				Method:  http.MethodGet,
				Version: api.AuraApiVersion2,
			})
			if err != nil {
				return err
			}

			if api.IsSuccessful(statusCode) {
				fields := []string{
					"name",
					"type",
					"current_status",
					"last_committed_txn",
					"last_seen",
					"replication_lag",
					"role",
					"writer",
				}
				output.PrintBody(cmd, cfg, resBody, fields)
			}

			return nil
		},
	}
	cmd.Flags().StringVarP(&organizationId, organizationIdFlag, "o", "", "(required) Organization ID")
	cmd.Flags().StringVarP(&projectId, projectIdFlag, "p", "", "(required) Project/tenant ID")
	cmd.Flags().StringVarP(&deploymentId, deploymentIdFlag, "d", "", "(required) Deployment ID")
	cmd.Flags().StringVarP(&serverId, serverIdFlag, "s", "", "(required) Server ID")

	err := cmd.MarkFlagRequired(deploymentIdFlag)
	if err != nil {
		log.Fatal(err)
	}
	err = cmd.MarkFlagRequired(serverIdFlag)
	if err != nil {
		log.Fatal(err)
	}

	return cmd
}

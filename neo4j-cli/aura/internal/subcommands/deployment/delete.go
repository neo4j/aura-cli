package deployment

import (
	"fmt"
	"log"
	"net/http"

	"github.com/neo4j/cli/common/clicfg"
	"github.com/neo4j/cli/neo4j-cli/aura/internal/api"
	"github.com/spf13/cobra"
)

func NewDeleteCmd(cfg *clicfg.Config) *cobra.Command {
	var (
		organizationId string
		projectId      string
	)

	const (
		organizationIdFlag = "organization-id"
		projectIdFlag      = "project-id"
	)

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete the given deployment",
		Long:  "Deletes the given Fleet Manager deployment. This will only delete the deployment from Fleet Manager without affecting the actual running database. It is advised to disable Fleet Management for the database using `call fleetManagement.disable()`",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if cfg.Aura.DefaultOrganization() == "" {
				err := cmd.MarkFlagRequired(organizationIdFlag)
				if err != nil {
					log.Fatal(err)
				}
			}

			if cfg.Aura.DefaultProject() == "" {
				err := cmd.MarkFlagRequired(projectIdFlag)
				if err != nil {
					log.Fatal(err)
				}
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if organizationId == "" {
				organizationId = cfg.Aura.DefaultOrganization()
			}
			if projectId == "" {
				projectId = cfg.Aura.DefaultProject()
			}
			deploymentId := args[0]
			path := fmt.Sprintf("/organizations/%s/projects/%s/fleet-manager/deployments/%s", organizationId, projectId, deploymentId)

			cmd.SilenceUsage = true
			_, statusCode, err := api.MakeRequest(cfg, path, &api.RequestConfig{
				Method:  http.MethodDelete,
				Version: api.AuraApiVersion2,
			})
			if err != nil {
				return err
			}

			if api.IsSuccessful(statusCode) {
				cmd.Println("Deployment deleted successfully", deploymentId)
			}

			return nil
		},
	}
	cmd.Flags().StringVar(&organizationId, organizationIdFlag, "", "(required) Organization ID")
	cmd.Flags().StringVar(&projectId, projectIdFlag, "", "(required) Project/tenant ID")

	return cmd
}

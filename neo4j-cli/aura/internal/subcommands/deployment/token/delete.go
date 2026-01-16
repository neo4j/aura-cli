package token

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
		deploymentId   string
	)

	const (
		organizationIdFlag = "organization-id"
		projectIdFlag      = "project-id"
		deploymentIdFlag   = "deployment-id"
	)

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete the deployment token",
		Long:  "Deletes the token for the given Fleet Manager deployment. After deleting the token, users should also disable Fleet Manager from the database using `call fleetManagement.disable();`",
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
			path := fmt.Sprintf("/organizations/%s/projects/%s/fleet-manager/deployments/%s/token", organizationId, projectId, deploymentId)

			cmd.SilenceUsage = true
			_, statusCode, err := api.MakeRequest(cfg, path, &api.RequestConfig{
				Method:  http.MethodDelete,
				Version: api.AuraApiVersion2,
			})
			if err != nil {
				return err
			}

			if api.IsSuccessful(statusCode) {
				cmd.Println("Deployment token deleted successfully for deployment", deploymentId)
			}

			return nil
		},
	}
	cmd.Flags().StringVar(&organizationId, organizationIdFlag, "", "(required) Organization ID")
	cmd.Flags().StringVar(&projectId, projectIdFlag, "", "(required) Project/tenant ID")
	cmd.Flags().StringVar(&deploymentId, deploymentIdFlag, "", "(required) Deployment ID")

	err := cmd.MarkFlagRequired(deploymentIdFlag)
	if err != nil {
		log.Fatal(err)
	}

	return cmd
}

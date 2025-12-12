package token

import (
	"fmt"
	"log"
	"net/http"

	"github.com/neo4j/cli/common/clicfg"
	"github.com/neo4j/cli/neo4j-cli/aura/internal/api"
	"github.com/neo4j/cli/neo4j-cli/aura/internal/output"
	"github.com/spf13/cobra"
)

func NewUpdateCmd(cfg *clicfg.Config) *cobra.Command {
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
		Use:   "update",
		Short: "Update the deployment token incase it needs to be rotated manually.",
		Long:  "Creates a new auto rotating Fleet Manager deployment token with a three month rotation interval. The token should be registered to the database again using `call fleetManagement.registerToken('$token');`",
		Args:  cobra.ExactArgs(0),
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
			path := fmt.Sprintf("/organizations/%s/projects/%s/fleet-manager/deployments/%s/token", organizationId, projectId, deploymentId)

			cmd.SilenceUsage = true
			resBody, statusCode, err := api.MakeRequest(cfg, path, &api.RequestConfig{
				Method:   http.MethodPatch,
				PostBody: map[string]any{},
				Version:  api.AuraApiVersion2,
			})
			if err != nil {
				return err
			}

			if api.IsSuccessful(statusCode) {
				output.PrintBody(cmd, cfg, resBody, []string{"token"})
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

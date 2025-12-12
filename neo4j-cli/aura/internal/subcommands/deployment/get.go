package deployment

import (
	"fmt"
	"net/http"

	"github.com/neo4j/cli/common/clicfg"
	"github.com/neo4j/cli/neo4j-cli/aura/internal/api"
	"github.com/neo4j/cli/neo4j-cli/aura/internal/output"
	"github.com/spf13/cobra"
)

func NewGetCmd(cfg *clicfg.Config) *cobra.Command {
	var (
		organizationId string
		projectId      string
	)

	const (
		organizationIdFlag = "organization-id"
		projectIdFlag      = "project-id"
	)

	cmd := &cobra.Command{
		Use:     "get <id>",
		Short:   "Returns deployment details",
		Long:    "Returns details about a specific Fleet Manager deployment.",
		Args:    cobra.ExactArgs(1),
		PreRunE: cfg.Aura.PreRunWithDefaultOrganizationAndProject(organizationId, projectId),
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
			resBody, statusCode, err := api.MakeRequest(cfg, path, &api.RequestConfig{
				Method:  http.MethodGet,
				Version: api.AuraApiVersion2,
			})
			if err != nil {
				return err
			}

			if api.IsSuccessful(statusCode) {
				fields := []string{
					"id",
					"name",
					"dbms:edition",
					"dbms:packaging",
					"token:expiry_time",
					"token:auto_rotated",
					"token:creation_time",
				}
				output.PrintBody(cmd, cfg, resBody, fields)
			}

			return nil
		},
	}
	cmd.Flags().StringVar(&organizationId, organizationIdFlag, "", "(required) Organization ID")
	cmd.Flags().StringVar(&projectId, projectIdFlag, "", "(required) Project/tenant ID")

	return cmd
}

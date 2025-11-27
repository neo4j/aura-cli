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

func NewCreateCmd(cfg *clicfg.Config) *cobra.Command {
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
		Use:   "create",
		Short: "Create a deployment token",
		Long:  "Create a new auto rototating Fleet Manager deployment token with a three month rotation interval. Register the deployment with Fleet Manager with the `call fleetManagement.registerToken('$token');` database procedure.",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := fmt.Sprintf("/organizations/%s/projects/%s/fleet-manager/deployments/%s/token", organizationId, projectId, deploymentId)

			cmd.SilenceUsage = true
			resBody, statusCode, err := api.MakeRequest(cfg, path, &api.RequestConfig{
				Method:   http.MethodPost,
				PostBody: map[string]any{},
				Version:  api.AuraApiVersion2,
			})
			if err != nil {
				return err
			}

			// NOTE: Deployment register should not return OK (200), it always returns 201, checking both just in case
			if statusCode == http.StatusCreated || statusCode == http.StatusOK {
				output.PrintBody(cmd, cfg, resBody, []string{"token"})
			}

			return nil
		},
	}
	cmd.Flags().StringVarP(&organizationId, organizationIdFlag, "o", "", "(required) Organization ID")
	cmd.Flags().StringVarP(&projectId, projectIdFlag, "p", "", "(required) Project/tenant ID")
	cmd.Flags().StringVarP(&deploymentId, deploymentIdFlag, "d", "", "(required) Deployment ID")

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

	return cmd
}

package deployment

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
		name           string
		connectionUrl  string
	)

	const (
		organizationIdFlag = "organization-id"
		projectIdFlag      = "project-id"
		nameFlag           = "name"
		connectionUrlFlag  = "connection-url"
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new deployment",
		Long:  "This endpoint creates a new Fleet Manager deployment.",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := fmt.Sprintf("/organizations/%s/projects/%s/deployments", organizationId, projectId)

			body := map[string]any{
				"name":           name,
				"connection_url": connectionUrl,
			}

			cmd.SilenceUsage = true
			resBody, statusCode, err := api.MakeRequest(cfg, path, &api.RequestConfig{
				Method:   http.MethodPost,
				PostBody: body,
				Version:  api.AuraApiVersion2,
			})
			if err != nil {
				return err
			}

			// NOTE: Deployment create should not return OK (200), it always returns 201, checking both just in case
			if statusCode == http.StatusCreated || statusCode == http.StatusOK {
				output.PrintBody(cmd, cfg, resBody, []string{"id"})
			}

			return nil
		},
	}
	cmd.Flags().StringVarP(&organizationId, organizationIdFlag, "o", "", "(required) Organization ID")
	cmd.Flags().StringVarP(&projectId, projectIdFlag, "p", "", "(required) Project/tenant ID")
	cmd.Flags().StringVarP(&name, nameFlag, "n", "", "(required) Deployment name")
	cmd.Flags().StringVarP(&connectionUrl, connectionUrlFlag, "c", "", "An optional connection URL for the deployment")

	err := cmd.MarkFlagRequired(organizationIdFlag)
	if err != nil {
		log.Fatal(err)
	}
	err = cmd.MarkFlagRequired(projectIdFlag)
	if err != nil {
		log.Fatal(err)
	}
	err = cmd.MarkFlagRequired(nameFlag)
	if err != nil {
		log.Fatal(err)
	}

	return cmd
}

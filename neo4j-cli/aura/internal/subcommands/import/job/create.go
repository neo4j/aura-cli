package job

import (
	"fmt"
	"github.com/neo4j/cli/common/clicfg"
	"github.com/neo4j/cli/neo4j-cli/aura/internal/api"
	"github.com/neo4j/cli/neo4j-cli/aura/internal/output"
	"github.com/spf13/cobra"
	"log"
	"net/http"
)

func NewCreateCmd(cfg *clicfg.Config) *cobra.Command {
	var (
		organizationId string
		projectId      string
		importModelId  string
		auraDbId       string
		user           string
		password       string
	)

	const (
		organizationIdFlag = "organization-id"
		projectIdFlag      = "project-id"
		importModelIdFlag  = "import-model-id"
		auraDbIdFlag       = "aura-db-id"
		userFlag           = "user"
		passwordFlag       = "password"
	)
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Allows you to create a new import job",
		RunE: func(cmd *cobra.Command, args []string) error {
			path := fmt.Sprintf("/organizations/%s/projects/%s/import/jobs", organizationId, projectId)

			responseBody, statusCode, err := api.MakeRequest(cfg, path, &api.RequestConfig{
				Method:  http.MethodPost,
				Version: api.AuraApiVersion2,
				PostBody: map[string]any{
					"importModelId": importModelId,
					"auraCredentials": map[string]any{
						"dbId":     auraDbId,
						"user":     user,
						"password": password,
					},
				},
			})
			if err != nil || statusCode != 201 {
				return err
			}
			output.PrintBody(cmd, cfg, responseBody, []string{"id"})
			return nil
		},
	}

	cmd.Flags().StringVar(&organizationId, organizationIdFlag, "", "Sets the organization ID the job belongs to")
	cmd.Flags().StringVar(&projectId, projectIdFlag, "", "Project/Tenant ID")
	cmd.Flags().StringVar(&importModelId, importModelIdFlag, "", "Import model ID, you can find it from your Aura Console")
	cmd.Flags().StringVar(&auraDbId, auraDbIdFlag, "", "Aura DB ID targeting for import data goes in")
	cmd.Flags().StringVar(&user, userFlag, "", "Username to use for authentication")
	cmd.Flags().StringVar(&password, passwordFlag, "", "Password to use for authentication")
	err := cmd.MarkFlagRequired(organizationIdFlag)
	if err != nil {
		log.Fatal(err)
	}
	err = cmd.MarkFlagRequired(projectIdFlag)
	if err != nil {
		log.Fatal(err)
	}
	err = cmd.MarkFlagRequired(importModelIdFlag)
	if err != nil {
		log.Fatal(err)
	}
	err = cmd.MarkFlagRequired(auraDbIdFlag)
	if err != nil {
		log.Fatal(err)
	}

	return cmd
}

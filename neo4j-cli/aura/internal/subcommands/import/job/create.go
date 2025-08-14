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
		projectId     string
		importModelId string
		auraDbId      string
	)

	const (
		projectIdFlag     = "project-id"
		importModelIdFlag = "import-model-id"
		auraDbIdFlag      = "aura-db-id"
	)
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Allows you to create a new import job",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if projectId == "" {
				return fmt.Errorf("projectId is required")
			}
			if importModelId == "" {
				return fmt.Errorf("importModelId is required")
			}
			if auraDbId == "" {
				return fmt.Errorf("auraDbId is required")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			path := fmt.Sprintf("/projects/%s/import/jobs", projectId)

			responseBody, statusCode, err := api.MakeRequest(cfg, path, &api.RequestConfig{
				Method:  http.MethodPost,
				Version: api.AuraApiVersion2,
				PostBody: map[string]any{
					"importModelId": importModelId,
					"auraCredentials": map[string]any{
						"dbId": auraDbId,
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

	cmd.Flags().StringVar(&projectId, projectIdFlag, "", "Project/Tenant ID")
	cmd.Flags().StringVar(&importModelId, importModelIdFlag, "", "Import model ID, you can find it from your Aura Console")
	cmd.Flags().StringVar(&auraDbId, auraDbIdFlag, "", "Aura DB ID targeting for import data goes in")
	err := cmd.MarkFlagRequired(projectIdFlag)
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

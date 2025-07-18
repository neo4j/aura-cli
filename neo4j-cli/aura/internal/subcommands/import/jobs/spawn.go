package jobs

import (
	"fmt"
	"github.com/neo4j/cli/common/clicfg"
	"github.com/neo4j/cli/neo4j-cli/aura/internal/api"
	"github.com/spf13/cobra"
	"log"
	"net/http"
)

func NewSpawnCmd(cfg *clicfg.Config) *cobra.Command {
	var (
		organizationId string
		projectId      string
		importModelId  string
		auraDbId       string
	)

	const (
		organizationIdFlag = "organization-id"
		projectIdFlag      = "project-id"
		importModelIdFlag  = "import-model-id"
		auraDbIdFlag       = "aura-db-id"
	)
	cmd := &cobra.Command{
		Use:   "spawn",
		Short: "Allows you to spawn your import jobs",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if organizationId == "" {
				return fmt.Errorf("organizationId is required")
			}
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
			path := fmt.Sprintf("/organizations/%s/projects/%s/import/jobs", organizationId, projectId)

			resBody, statusCode, err := api.MakeRequest(cfg, path, &api.RequestConfig{
				Method:  http.MethodPost,
				Version: api.AuraApiVersion2,
				PostBody: map[string]any{
					"importModelId": importModelId,
					"auraCredentials": map[string]any{
						"dbId": auraDbId,
					},
				},
			})
			log.Printf(fmt.Sprintf("Response body: %+v\n", string(resBody)))
			log.Printf(fmt.Sprintf("Response status code: %d\n", statusCode))
			if err != nil {
				return err
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&organizationId, organizationIdFlag, "", "Organization ID")
	cmd.Flags().StringVar(&projectId, projectIdFlag, "", "Project ID")
	cmd.Flags().StringVar(&importModelId, importModelIdFlag, "", "Import model id")
	cmd.Flags().StringVar(&auraDbId, auraDbIdFlag, "", "Aura DB ID")
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

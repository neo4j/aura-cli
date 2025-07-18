package jobs

import (
	"fmt"
	"github.com/neo4j/cli/common/clicfg"
	"github.com/neo4j/cli/neo4j-cli/aura/internal/api"
	"github.com/spf13/cobra"
	"log"
	"net/http"
)

func NewGetCmd(cfg *clicfg.Config) *cobra.Command {
	var (
		organizationId string
		projectId      string
		jobId          string
	)

	const (
		organizationIdFlag = "organization-id"
		projectIdFlag      = "project-id"
		jobIdFlag          = "job-id"
	)

	cmd := &cobra.Command{
		Use:   "get <id>",
		Short: "Get a job by id",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if organizationId == "" {
				return fmt.Errorf("organizationId is required")
			}
			if projectId == "" {
				return fmt.Errorf("projectId is required")
			}
			if jobId == "" {
				return fmt.Errorf("importModelId is required")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			path := fmt.Sprintf("/organizations/%s/projects/%s/import/jobs/%s", organizationId, projectId, jobId)

			resBody, statusCode, err := api.MakeRequest(cfg, path, &api.RequestConfig{
				Method: http.MethodGet,
				UseV2:  true,
			})
			log.Printf(fmt.Sprintf("Response body: %+v\n", string(resBody)))
			log.Printf(fmt.Sprintf("Response status code: %d\n", statusCode))
			if err != nil {
				log.Fatal(err)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&organizationId, organizationIdFlag, "", "Organization ID")
	cmd.Flags().StringVar(&projectId, projectIdFlag, "", "Project ID")
	cmd.Flags().StringVar(&jobId, jobIdFlag, "", "Import job id")
	err := cmd.MarkFlagRequired(organizationIdFlag)
	if err != nil {
		log.Fatal(err)
	}
	err = cmd.MarkFlagRequired(projectIdFlag)
	if err != nil {
		log.Fatal(err)
	}
	err = cmd.MarkFlagRequired(jobIdFlag)
	if err != nil {
		log.Fatal(err)
	}
	return cmd
}

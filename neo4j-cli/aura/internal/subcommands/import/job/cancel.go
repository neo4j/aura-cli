package job

import (
	"fmt"
	"log"
	"net/http"

	"github.com/neo4j/cli/common/clicfg"
	"github.com/neo4j/cli/neo4j-cli/aura/internal/api"
	"github.com/neo4j/cli/neo4j-cli/aura/internal/output"
	"github.com/spf13/cobra"
)

func NewCancelCommand(cfg *clicfg.Config) *cobra.Command {
	var (
		organizationId string
		projectId      string
		jobId          string
	)

	const (
		organizationIdFlag = "organization-id"
		projectIdFlag      = "project-id"
	)

	cmd := &cobra.Command{
		Use:   "cancel <id>",
		Short: "Cancel a job by id",
		Args:  cobra.ExactArgs(1),
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
			jobId = args[0]
			path := fmt.Sprintf("/organizations/%s/projects/%s/import/jobs/%s/cancellation", organizationId, projectId, jobId)
			resBody, statusCode, err := api.MakeRequest(cfg, path, &api.RequestConfig{
				Method:  http.MethodPost,
				Version: api.AuraApiVersion2,
			})
			if err != nil || statusCode != http.StatusOK {
				return err
			}
			output.PrintBody(cmd, cfg, resBody, []string{"id"})

			return nil
		},
	}
	cmd.Flags().StringVar(&organizationId, organizationIdFlag, "", "(required) Organization ID")
	cmd.Flags().StringVar(&projectId, projectIdFlag, "", "(required) Project/tenant ID")

	return cmd
}

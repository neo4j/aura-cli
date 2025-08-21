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

func NewCancelCommand(cfg *clicfg.Config) *cobra.Command {
	var (
		projectId string
		jobId     string
	)

	const (
		projectIdFlag = "project-id"
	)

	cmd := &cobra.Command{
		Use:   "cancel <id>",
		Short: "Cancel a job by id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			jobId = args[0]
			path := fmt.Sprintf("/projects/%s/import/jobs/%s/cancel", projectId, jobId)
			resBody, statusCode, err := api.MakeRequest(cfg, path, &api.RequestConfig{
				Method:  http.MethodPatch,
				Version: api.AuraApiVersion2,
			})
			if err != nil || statusCode != http.StatusOK {
				return err
			}
			output.PrintBody(cmd, cfg, resBody, []string{"id"})

			return nil
		},
	}
	cmd.Flags().StringVar(&projectId, projectIdFlag, "", "Project ID")
	err := cmd.MarkFlagRequired(projectIdFlag)
	if err != nil {
		log.Fatal(err)
	}

	return cmd
}

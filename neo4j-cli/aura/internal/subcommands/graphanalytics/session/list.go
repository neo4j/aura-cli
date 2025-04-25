package session

import (
	"net/http"

	"github.com/neo4j/cli/common/clicfg"
	"github.com/neo4j/cli/neo4j-cli/aura/internal/api"
	"github.com/neo4j/cli/neo4j-cli/aura/internal/output"
	"github.com/spf13/cobra"
)

func NewListCmd(cfg *clicfg.Config) *cobra.Command {
	var projectId string
	var instanceId string
	var organizationId string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "Returns a list of Graph Analytics Serverless sessions",
		Long: `This subcommand returns a list containing a summary of each of your Graph Analytics Serverless session
				By default, this subcommand lists all sessions a user has access to across all projects.
				You can filter sessions in a particular project using:
				--organization-id <organization-id>
				--project-id <project-id>
				--instance-id <instance-id>
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			path := "/graph-analytics/sessions"

			queryParams := make(map[string]string)
			if organizationId != "" {
				queryParams["organizationId"] = organizationId
			}
			if projectId != "" {
				queryParams["projectId"] = projectId
			}
			if instanceId != "" {
				queryParams["instanceId"] = instanceId
			}

			cmd.SilenceUsage = true
			resBody, statusCode, err := api.MakeRequest(cfg, path, &api.RequestConfig{
				Method:      http.MethodGet,
				QueryParams: queryParams,
			})
			if err != nil {
				return err
			}

			if statusCode == http.StatusOK {
				output.PrintBody(cmd, cfg, resBody, []string{"id", "name", "status", "project_id", "cloud_provider", "ttl"})
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&projectId, "project-id", "", "An optional Project ID to filter sessions in a project")
	cmd.Flags().StringVar(&organizationId, "organization-id", "", "An optional Organization ID to filter sessions in an organization")
	cmd.Flags().StringVar(&instanceId, "instance-id", "", "An optional Instance ID to filter for sessions attached to an instance")

	return cmd
}

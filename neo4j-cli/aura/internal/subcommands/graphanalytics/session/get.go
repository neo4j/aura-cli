package session

import (
	"fmt"
	"net/http"

	"github.com/neo4j/cli/common/clicfg"
	"github.com/neo4j/cli/neo4j-cli/aura/internal/api"
	"github.com/neo4j/cli/neo4j-cli/aura/internal/output"
	"github.com/spf13/cobra"
)

func NewGetCmd(cfg *clicfg.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Get a Graph Analytics Serverless session",
		Long:  `This subcommand returns the details of a Graph Analytics Serverless session.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			path := fmt.Sprintf("/graph-analytics/sessions/%s", args[0])

			cmd.SilenceUsage = true
			resBody, statusCode, err := api.MakeRequest(cfg, path, &api.RequestConfig{
				Method: http.MethodGet,
			})
			if err != nil {
				return err
			}

			if statusCode == http.StatusOK {
				output.PrintBody(cmd, cfg, resBody, []string{
					"id",
					"name",
					"memory",
					"status",
					"created_at",
					"user_id",
					"project_id",
					"cloud_provider",
					"region",
					"host",
					"expiry_date",
					"instance_id",
				})
			}
			return nil
		},
	}
	return cmd
}

package graphql

import (
	"fmt"
	"net/http"

	"github.com/neo4j/cli/common/clicfg"
	"github.com/neo4j/cli/neo4j/aura/internal/api"
	"github.com/neo4j/cli/neo4j/aura/internal/output"
	"github.com/spf13/cobra"
)

func NewResumeCmd(cfg *clicfg.Config) *cobra.Command {
	var (
		instanceId string
		await      bool
	)

	cmd := &cobra.Command{
		Use:   "resume <id>",
		Short: "Resume a GraphQL Data API",
		Long: `This command starts the resuming process of an existing GraphQL Data API.

Resuming a GraphQL Data API is an asynchronous operation. Use the --await flag to wait for the GraphQL Data API to be ready. Once the status transitions from "resuming" to "ready" you may begin to use your GraphQL Data API.	`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			path := fmt.Sprintf("/instances/%s/data-apis/graphql/%s/resume", instanceId, args[0])

			resBody, statusCode, err := api.MakeRequest(cfg, path, &api.RequestConfig{
				Method: http.MethodPost,
			})
			if err != nil {
				return err
			}

			// NOTE: resume should not return OK (200), it always returns 202, checking both just in case
			if statusCode == http.StatusAccepted || statusCode == http.StatusOK {
				output.PrintBody(cmd, cfg, resBody, []string{"id", "name", "status", "url"})

				if await {
					cmd.Println("Waiting for GraphQL Data API to be resumed...")
					pollResponse, err := api.PollGraphQLDataApi(cfg, instanceId, args[0], api.GraphQLDataApiStatusResuming)
					if err != nil {
						return err
					}

					cmd.Println("GraphQL Data API Status:", pollResponse.Data.Status)
				}
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&instanceId, "instance-id", "", "The ID of the instance to resume the Data API for")
	cmd.MarkFlagRequired("instance-id")

	cmd.Flags().BoolVar(&await, "await", false, "Waits until GraphQL Data API is resumed.")

	return cmd
}

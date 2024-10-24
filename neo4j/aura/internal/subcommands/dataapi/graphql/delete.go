package graphql

import (
	"fmt"
	"net/http"

	"github.com/neo4j/cli/common/clicfg"
	"github.com/neo4j/cli/neo4j/aura/internal/api"
	"github.com/neo4j/cli/neo4j/aura/internal/output"
	"github.com/spf13/cobra"
)

func NewDeleteCmd(cfg *clicfg.Config) *cobra.Command {
	var (
		instanceId string
		await      bool
	)

	cmd := &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a GraphQL Data API",
		Long:  "Deletes a GraphQL Data API. This action can _not_ be undone.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			path := fmt.Sprintf("/instances/%s/data-apis/graphql/%s", instanceId, args[0])

			resBody, statusCode, err := api.MakeRequest(cfg, path, &api.RequestConfig{
				Method: http.MethodDelete,
			})
			if err != nil {
				return err
			}

			// NOTE: delete should not return OK (200), it always returns 202, checking both just in case
			if statusCode == http.StatusAccepted || statusCode == http.StatusOK {
				err = output.PrintBody(cmd, cfg, resBody, []string{"id", "name", "status", "url"})
				if err != nil {
					return err
				}

				if await {
					cmd.Println("Waiting for GraphQL Data API to be deleted...")
					pollResponse, err := api.PollGraphQLDataApi(cfg, instanceId, args[0], api.GraphQLDataApiStatusDeleting)
					if err != nil {
						return err
					}

					cmd.Println("GraphQL Data API Status:", pollResponse.Data.Status)
				}
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&instanceId, "instance-id", "", "The ID of the instance to delete the Data API for")
	cmd.MarkFlagRequired("instance-id")

	cmd.Flags().BoolVar(&await, "await", false, "Waits until GraphQL Data API is deleted.")

	return cmd
}

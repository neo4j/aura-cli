package graphql

import (
	"fmt"
	"net/http"

	"github.com/neo4j/cli/common/clicfg"
	"github.com/neo4j/cli/neo4j/aura/internal/api"
	"github.com/neo4j/cli/neo4j/aura/internal/output"
	"github.com/spf13/cobra"
)

func NewUpdateCmd(cfg *clicfg.Config) *cobra.Command {
	const (
		instanceIdFlag             = "instance-id"
		nameFlag                   = "name"
		instanceUsernameFlag       = "instance-username"
		instancePasswordFlag       = "instance-password"
		typeDefsFlag               = "type-definitions"
		typeDefsFileFlag           = "type-definitions-file"
		featureSubgraphEnabledFlag = "feature-subgraph-enabled"
		awaitFlag                  = "await"
	)

	var (
		instanceId             string
		name                   string
		instanceUsername       string
		instancePassword       string
		typeDefs               string
		typeDefsFile           string
		featureSubgraphEnabled bool
		await                  bool
	)

	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Edit a GraphQL Data API",
		Long:  "This endpoint edits a specific GraphQL Data API.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			path := fmt.Sprintf("/instances/%s/data-apis/graphql/%s", instanceId, args[0])

			resBody, statusCode, err := api.MakeRequest(cfg, path, &api.RequestConfig{
				Method: http.MethodPatch,
			})
			if err != nil {
				return err
			}

			if statusCode == http.StatusAccepted {
				err = output.PrintBody(cmd, cfg, resBody, []string{"id", "name", "status", "url"})
				if err != nil {
					return err
				}
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&instanceId, instanceIdFlag, "", "The ID of the instance to get the Data API for")
	cmd.MarkFlagRequired("instance-id")

	return cmd
}

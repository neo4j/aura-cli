// Copyright (c) "Neo4j"
// Neo4j Sweden AB [http://neo4j.com]

package agent

import (
	"fmt"
	"net/http"

	"github.com/neo4j/cli/common/clicfg"
	"github.com/neo4j/cli/neo4j-cli/aura/internal/api"
	"github.com/neo4j/cli/neo4j-cli/aura/internal/output"
	"github.com/neo4j/cli/neo4j-cli/aura/internal/subcommands/utils"
	"github.com/spf13/cobra"
)

func NewListCmd(cfg *clicfg.Config) *cobra.Command {
	var (
		organizationId string
		projectId      string
	)

	const (
		organizationIdFlag = "organization-id"
		projectIdFlag      = "project-id"
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "Returns a list of agents",
		Long:  "Returns a list of agents for the specified project.",
		Args:  cobra.ExactArgs(0),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return utils.SetProjectFlagsAsRequired(cfg, cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			organizationId, projectId, err := utils.SetProjetDefaults(cfg, organizationId, projectId)
			if err != nil {
				return err
			}

			path := fmt.Sprintf("/organizations/%s/projects/%s/agents", organizationId, projectId)

			cmd.SilenceUsage = true
			resBody, statusCode, err := api.MakeRequest(cfg, path, &api.RequestConfig{
				Method:  http.MethodGet,
				Version: api.AuraApiVersion2,
			})
			if err != nil {
				return err
			}

			if api.IsSuccessful(statusCode) {
				output.PrintRawBody(cmd, cfg, resBody, []string{"id", "name", "description", "dbid", "enabled"})
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&organizationId, organizationIdFlag, "", "(required) Organization ID")
	cmd.Flags().StringVar(&projectId, projectIdFlag, "", "(required) Project/tenant ID")

	return cmd
}

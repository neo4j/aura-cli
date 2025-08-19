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

func NewGetCmd(cfg *clicfg.Config) *cobra.Command {
	var (
		projectId    string
		jobId        string
		showProgress bool
	)

	const (
		projectIdFlag    = "project-id"
		showProgressFlag = "progress"
	)

	cmd := &cobra.Command{
		Use:   "get <id>",
		Short: "Get a job by id",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			jobId = args[0]
			if jobId == "" {
				return fmt.Errorf("jobId is required")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			path := fmt.Sprintf("/projects/%s/import/jobs/%s", projectId, jobId)
			queryParams := make(map[string]string)
			queryParams["progress"] = fmt.Sprintf("%v", showProgress)

			resBody, statusCode, err := api.MakeRequest(cfg, path, &api.RequestConfig{
				Method:      http.MethodGet,
				Version:     api.AuraApiVersion2,
				QueryParams: queryParams,
			})
			if err != nil {
				return err
			}
			outputType := cfg.Aura.Output()

			if statusCode == http.StatusOK {
				output.PrintBody(cmd, cfg, resBody, []string{"id", "import_type", "info:state", "info:exit_status:state", "info:percentage_complete", "data_source:name", "aura_target:db_id"})
				if outputType != "json" {
					output.PrintBody(cmd, cfg, resBody, []string{"info:exit_status:message"})
				}
			}

			if showProgress && outputType != "json" {
				printJobProgressTable(cmd, cfg, resBody)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&projectId, projectIdFlag, "", "Project ID")
	cmd.Flags().BoolVar(&showProgress, showProgressFlag, false, "Show progress details")
	err := cmd.MarkFlagRequired(projectIdFlag)
	if err != nil {
		log.Fatal(err)
	}
	return cmd
}

func printJobProgressTable(cmd *cobra.Command, cfg *clicfg.Config, resBody []byte) {
	cmd.Println("# Progress details:")
	parsedBody := api.ParseBody(resBody)
	data, err := parsedBody.GetSingleOrError()
	if err != nil {
		panic(err)
	}
	info := data["info"].(map[string]interface{})
	progress := info["progress"].(map[string]interface{})
	nodes := progress["nodes"].([]interface{})
	wrappedNodes := make([]map[string]any, 0)
	for _, node := range nodes {
		wrappedNodes = append(wrappedNodes, node.(map[string]interface{}))
	}
	nodesResponseData := api.NewListResponseData(wrappedNodes)
	cmd.Println("# Nodes progress:")
	output.PrintBodyMap(cmd, cfg, nodesResponseData, []string{"id", "labels", "processed_rows", "total_rows", "created_nodes", "created_constraints", "created_indexes"})

	relationships := progress["relationships"].([]interface{})
	wrappedRelationships := make([]map[string]any, 0)
	for _, relationship := range relationships {
		wrappedRelationships = append(wrappedRelationships, relationship.(map[string]interface{}))
	}
	relationshipsResponseData := api.NewListResponseData(wrappedRelationships)
	cmd.Println("# Relationships progress:")
	output.PrintBodyMap(cmd, cfg, relationshipsResponseData, []string{"id", "type", "processed_rows", "total_rows", "created_relationships", "created_constraints", "created_indexes"})
}

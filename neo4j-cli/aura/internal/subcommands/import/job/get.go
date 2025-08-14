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
		jobIdFlag        = "job-id"
		showProgressFlag = "progress"
	)

	cmd := &cobra.Command{
		Use:   "get <id>",
		Short: "Get a job by id",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if projectId == "" {
				return fmt.Errorf("projectId is required")
			}
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
			if statusCode == http.StatusOK {
				output.PrintBody(cmd, cfg, resBody, []string{"id", "import_type", "info:state", "info:exit_status:state", "info:percentage_complete", "data_source:name", "aura_target:db_id"})
				output.PrintBody(cmd, cfg, resBody, []string{"info:exit_status:message"})
			}

			outputType := cfg.Aura.Output()
			if showProgress && outputType != "json" {
				cmd.Println("###############################")
				cmd.Println("# The progress details are shown as follows.")
				cmd.Println("###############################")
				parsedBody := api.ParseBody(resBody)
				data, err := parsedBody.GetSingleOrError()
				if err != nil {
					panic(err)
				}
				//log.Printf("nodes: %v", data)
				info := data["info"].(map[string]interface{})
				//log.Printf("info: %v", info)
				progress := info["progress"].(map[string]interface{})
				//log.Printf("progress: %v", progress)
				nodes := progress["nodes"].([]interface{})
				//log.Printf("nodes: %v", nodes)
				wrappedNodes := make([]map[string]any, 0)
				for _, node := range nodes {
					wrappedNodes = append(wrappedNodes, node.(map[string]interface{}))
				}
				nodesResponseData := api.NewListResponseData(wrappedNodes)
				//log.Printf("nodesResponseData: %v", nodesResponseData)
				cmd.Println("###############################")
				cmd.Println("# Nodes progress details are shown as follows.")
				cmd.Println("###############################")
				output.PrintBodyMap(cmd, cfg, nodesResponseData, []string{"id", "processed_rows", "total_rows", "created_nodes", "created_constraints", "created_indexes"})

				relationships := progress["relationships"].([]interface{})
				//log.Printf("nodes: %v", nodes)
				wrappedRelationships := make([]map[string]any, 0)
				for _, relationship := range relationships {
					wrappedRelationships = append(wrappedRelationships, relationship.(map[string]interface{}))
				}
				relationshipsResponseData := api.NewListResponseData(wrappedRelationships)
				//log.Printf("nodesResponseData: %v", nodesResponseData)
				cmd.Println("###############################")
				cmd.Println("# Relationships progress details are shown as follows.")
				cmd.Println("###############################")
				output.PrintBodyMap(cmd, cfg, relationshipsResponseData, []string{"id", "processed_rows", "total_rows", "created_relationships", "created_constraints", "created_indexes"})
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&projectId, projectIdFlag, "", "Project ID")
	cmd.Flags().StringVar(&jobId, jobIdFlag, "", "Import job ID")
	cmd.Flags().BoolVar(&showProgress, showProgressFlag, false, "Show progress details")
	err := cmd.MarkFlagRequired(projectIdFlag)
	if err != nil {
		log.Fatal(err)
	}
	err = cmd.MarkFlagRequired(jobIdFlag)
	if err != nil {
		log.Fatal(err)
	}
	return cmd
}

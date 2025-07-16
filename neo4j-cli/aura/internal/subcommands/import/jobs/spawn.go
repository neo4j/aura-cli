package jobs

import (
	"fmt"
	"github.com/neo4j/cli/common/clicfg"
	"github.com/neo4j/cli/neo4j-cli/aura/internal/api"
	"github.com/spf13/cobra"
	"log"
	"net/http"
)

func NewSpawnCmd(cfg *clicfg.Config) *cobra.Command {
	var (
		importModelId string
	)

	const (
		importModelIdFlag = "import-model-id"
	)
	cmd := &cobra.Command{
		Use:   "spawn",
		Short: "Allows you to spawn your import jobs",
		PostRunE: func(cmd *cobra.Command, args []string) error {
			err := cmd.MarkFlagRequired(importModelIdFlag)
			if err != nil {
				return err
			}
			if importModelId == "" {
				return fmt.Errorf("importModelId is required")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			path := "/import/jobs"

			resBody, statusCode, err := api.MakeRequest(cfg, path, &api.RequestConfig{
				Method:  http.MethodPost,
				Version: api.AuraApiVersion2,
			})
			log.Printf(fmt.Sprintf("Response body: %+v\n", resBody))
			log.Printf(fmt.Sprintf("Response status code: %d\n", statusCode))
			if err != nil {
				log.Fatal(err)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&importModelId, importModelIdFlag, "", "Import model id")
	return cmd
}

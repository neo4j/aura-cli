package allowedorigin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/neo4j/cli/common/clicfg"
	"github.com/neo4j/cli/common/clierr"
	"github.com/neo4j/cli/neo4j-cli/aura/internal/api"
	"github.com/spf13/cobra"
)

func NewCmd(cfg *clicfg.Config) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "allowed-origin",
		Short: "Adds an exact Cross-Origin Resource Sharing (CORS) allowed origin for a specific GraphQL Data API",
	}

	cmd.AddCommand(NewAddCmd(cfg))

	return cmd
}

type DetailedBody struct {
	Data DataBody `json:"data"`
}

type DataBody struct {
	Security SecurityBody `json:"security"`
}

type SecurityBody struct {
	CorsPolicy CorsPolicyBody `json:"cors_policy"`
}

type CorsPolicyBody struct {
	AllowedOrigins []string `json:"allowed_origins"`
}

func getGetExistingOrigins(cfg *clicfg.Config, dataApiId, instanceId string) ([]string, error) {
	getPath := fmt.Sprintf("/instances/%s/data-apis/graphql/%s", instanceId, dataApiId)
	getResBody, statusCode, err := api.MakeRequest(cfg, getPath, &api.RequestConfig{
		Method: http.MethodGet,
	})
	if err != nil {
		return nil, err
	}
	if statusCode != http.StatusOK {
		panic(clierr.NewFatalError("unexpected status code %d running CLI with args %s, please report an issue in https://github.com/neo4j/cli", statusCode, os.Args[1:]))
	}

	var parsedGetResBody DetailedBody
	err = json.Unmarshal(getResBody, &parsedGetResBody)
	if err != nil {
		panic(err)
	}

	return parsedGetResBody.Data.Security.CorsPolicy.AllowedOrigins, nil
}

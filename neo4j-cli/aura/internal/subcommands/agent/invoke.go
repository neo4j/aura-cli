// Copyright (c) "Neo4j"
// Neo4j Sweden AB [http://neo4j.com]

package agent

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/neo4j/cli/common/clicfg"
	"github.com/neo4j/cli/neo4j-cli/aura/internal/api"
	"github.com/neo4j/cli/neo4j-cli/aura/internal/subcommands/utils"
	"github.com/spf13/cobra"
)

func NewInvokeCmd(cfg *clicfg.Config) *cobra.Command {
	var (
		organizationId string
		projectId      string
		input          string
	)

	const (
		organizationIdFlag = "organization-id"
		projectIdFlag      = "project-id"
		inputFlag          = "input"
	)

	cmd := &cobra.Command{
		Use:   "invoke <id>",
		Short: "Invokes an agent with the given input",
		Long:  "Invokes an agent with the provided input string. Use --output json for the full response including content blocks.",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return utils.SetProjectFlagsAsRequired(cfg, cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			organizationId, projectId, err := utils.SetProjetDefaults(cfg, organizationId, projectId)
			if err != nil {
				return err
			}

			agentId := args[0]
			path := fmt.Sprintf("/organizations/%s/projects/%s/agents/%s/invoke", organizationId, projectId, agentId)

			body := map[string]any{
				"input": input,
			}

			cmd.SilenceUsage = true
			resBody, statusCode, err := api.MakeRequest(cfg, path, &api.RequestConfig{
				Method:   http.MethodPost,
				PostBody: body,
				Version:  api.AuraApiVersion2,
			})
			if err != nil {
				if statusCode == http.StatusForbidden {
					return fmt.Errorf("agent invocation forbidden: agent may be disabled or private")
				}
				return err
			}

			if api.IsSuccessful(statusCode) {
				// Check for application-level errors in the response (type: "error")
				var result map[string]any
				if jsonErr := json.Unmarshal(resBody, &result); jsonErr == nil {
					if resultType, ok := result["type"].(string); ok && resultType == "error" {
						if errObj, ok := result["error"].(map[string]any); ok {
							if msg, ok := errObj["message"].(string); ok {
								return fmt.Errorf("agent invocation failed: %s", msg)
							}
						}
						return fmt.Errorf("agent invocation failed")
					}
				}
				printInvokeResult(cmd, cfg, resBody)
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&organizationId, organizationIdFlag, "", "(required) Organization ID")
	cmd.Flags().StringVar(&projectId, projectIdFlag, "", "(required) Project/tenant ID")
	cmd.Flags().StringVar(&input, inputFlag, "", "(required) Input message to send to the agent")

	if err := cmd.MarkFlagRequired(inputFlag); err != nil {
		log.Fatal(err)
	}

	return cmd
}

// printInvokeResult prints the invoke response. In JSON mode the raw body is printed.
// In table/default mode the text content is printed followed by a stats line.
func printInvokeResult(cmd *cobra.Command, cfg *clicfg.Config, resBody []byte) {
	if cfg.Aura.Output() == "json" {
		printAgentItem(cmd, cfg, resBody, nil)
		return
	}
	var result map[string]any
	if err := json.Unmarshal(resBody, &result); err != nil {
		cmd.Println(string(resBody))
		return
	}

	contentBlocks, _ := result["content"].([]any)
	var texts []string
	toolCalls := 0
	for _, block := range contentBlocks {
		if m, ok := block.(map[string]any); ok {
			blockType, _ := m["type"].(string)
			switch {
			case blockType == "text":
				if text, ok := m["text"].(string); ok {
					texts = append(texts, text)
				}
			case strings.HasSuffix(blockType, "tool_use"):
				toolCalls++
			}
		}
	}

	if len(texts) > 0 {
		cmd.Println(strings.Join(texts, "\n"))
	}

	// Stats line
	status, _ := result["status"].(string)
	endReason, _ := result["end_reason"].(string)
	var reqTokens, resTokens, totalTokens int
	if usage, ok := result["usage"].(map[string]any); ok {
		reqTokens = int(toFloat64(usage["request_tokens"]))
		resTokens = int(toFloat64(usage["response_tokens"]))
		totalTokens = int(toFloat64(usage["total_tokens"]))
	}
	cmd.Printf("\nStatus: %s | End reason: %s | Tool calls: %d | Tokens: %d req / %d res / %d total\n",
		status, endReason, toolCalls, reqTokens, resTokens, totalTokens)
}

func toFloat64(v any) float64 {
	if f, ok := v.(float64); ok {
		return f
	}
	return 0
}

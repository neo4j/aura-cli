// Copyright (c) "Neo4j"
// Neo4j Sweden AB [http://neo4j.com]

package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/neo4j/cli/common/clicfg"
	"github.com/neo4j/cli/common/clierr"
	"github.com/neo4j/cli/neo4j-cli/aura/internal/api"
	"github.com/neo4j/cli/neo4j-cli/aura/internal/output"
	"github.com/spf13/cobra"
)

func NewCmd(cfg *clicfg.Config) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "agent",
		Short: "Relates to Aura Agents",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			cfg.Aura.BindBaseUrl(cmd.Flags().Lookup("base-url"))
			cfg.Aura.BindAuthUrl(cmd.Flags().Lookup("auth-url"))

			outputValue := cmd.Flags().Lookup("output").Value.String()
			if outputValue != "" {
				for _, v := range clicfg.ValidOutputValues {
					if v == outputValue {
						cfg.Aura.BindOutput(cmd.Flags().Lookup("output"))
						return nil
					}
				}
				return clierr.NewUsageError("invalid output value specified: %s", outputValue)
			}
			cfg.Aura.BindOutput(cmd.Flags().Lookup("output"))

			return nil
		},
	}

	cmd.AddCommand(NewGetCmd(cfg))
	cmd.AddCommand(NewListCmd(cfg))
	cmd.AddCommand(NewCreateCmd(cfg))
	cmd.AddCommand(NewUpdateCmd(cfg))
	cmd.AddCommand(NewPatchCmd(cfg))
	cmd.AddCommand(NewDeleteCmd(cfg))
	cmd.AddCommand(NewInvokeCmd(cfg))

	cmd.PersistentFlags().String("auth-url", "", "")
	cmd.PersistentFlags().String("base-url", "", "")
	cmd.PersistentFlags().String("output", "", fmt.Sprintf("Format to print console output in, from a choice of [%s]", strings.Join(clicfg.ValidOutputValues[:], ", ")))

	return cmd
}

// printAgentList prints a raw JSON array response from the agents API.
// In JSON mode the raw API body is printed; in table mode selected fields are shown.
func printAgentList(cmd *cobra.Command, cfg *clicfg.Config, resBody []byte, fields []string) {
	if cfg.Aura.Output() == "json" {
		var buf bytes.Buffer
		json.Indent(&buf, resBody, "", "\t")
		cmd.Println(buf.String())
		return
	}
	var items []map[string]any
	if err := json.Unmarshal(resBody, &items); err != nil {
		panic(err)
	}
	output.PrintBodyMap(cmd, cfg, api.NewListResponseData(items), fields)
}

// printAgentItem prints a raw JSON object response from the agents API.
// In JSON mode the raw API body is printed; in table mode selected fields are shown.
func printAgentItem(cmd *cobra.Command, cfg *clicfg.Config, resBody []byte, fields []string) {
	if cfg.Aura.Output() == "json" {
		var buf bytes.Buffer
		json.Indent(&buf, resBody, "", "\t")
		cmd.Println(buf.String())
		return
	}
	var item map[string]any
	if err := json.Unmarshal(resBody, &item); err != nil {
		panic(err)
	}
	output.PrintBodyMap(cmd, cfg, api.NewSingleValueResponseData(item), fields)
}

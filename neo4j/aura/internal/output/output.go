package output

import (
	"encoding/json"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"

	"github.com/neo4j/cli/common/clicfg"
	"github.com/neo4j/cli/neo4j/aura/internal/api"
)

func PrintBodyMap(cmd *cobra.Command, cfg *clicfg.Config, values api.ResponseData, fields []string) error {
	outputType := cfg.Aura.Output()

	switch output := outputType; output {
	case "json":
		bytes, err := json.MarshalIndent(values, "", "\t")
		if err != nil {
			return err
		}
		cmd.Println(string(bytes))
	case "table", "default":
		err := printTable(cmd, values, fields)
		if err != nil {
			return err
		}

	default:
		// This is in case the value is unknown
		cmd.Println(values)
	}

	return nil
}

func PrintBody(cmd *cobra.Command, cfg *clicfg.Config, body []byte, fields []string) error {
	if len(body) == 0 {
		return nil
	}
	values, err := api.ParseBody(body)
	if err != nil {
		return err
	}
	err = PrintBodyMap(cmd, cfg, values, fields)
	if err != nil {
		return err
	}

	return nil
}

func printTable(cmd *cobra.Command, responseData api.ResponseData, fields []string) error {
	t := table.NewWriter()

	header := table.Row{}
	for _, f := range fields {
		header = append(header, f)
	}

	t.AppendHeader(header)
	for _, v := range responseData.Data {
		row := table.Row{}
		for _, f := range fields {
			formattedValue := v[f]

			if v[f] == nil {
				formattedValue = ""
			}

			row = append(row, formattedValue)
		}
		t.AppendRow(row)
	}

	t.SetStyle(table.StyleLight)
	cmd.Println(t.Render())
	return nil
}

package output

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"

	"github.com/neo4j/cli/common/clicfg"
	"github.com/neo4j/cli/neo4j-cli/aura/internal/api"
)

func PrintBodyMap(cmd *cobra.Command, cfg *clicfg.Config, values api.ResponseData, fields []string) {
	outputType := cfg.Aura.Output()

	switch output := outputType; output {
	case "json":
		bytes, err := json.MarshalIndent(values, "", "\t")
		if err != nil {
			panic(err)
		}
		cmd.Println(string(bytes))
	case "table", "default":
		printTable(cmd, values, fields)
	default:
		// This is in case the value is unknown
		cmd.Println(values)
	}
}

func PrintBody(cmd *cobra.Command, cfg *clicfg.Config, body []byte, fields []string) {
	if len(body) == 0 {
		return
	}
	log.Println(string(body))
	values := api.ParseBody(body)
	log.Printf("values: %+v", values)

	PrintBodyMap(cmd, cfg, values, fields)
}

func printTable(cmd *cobra.Command, responseData api.ResponseData, fields []string) {
	t := table.NewWriter()

	header := table.Row{}
	for _, f := range fields {
		header = append(header, f)
	}

	t.AppendHeader(header)
	for _, v := range responseData.AsArray() {
		row := table.Row{}
		for _, f := range fields {
			subfields := strings.Split(f, ":")
			formattedValue := ""
			var nestedValues []map[string]any
			nestedValues = append(nestedValues, v)
			//log.Printf("nestedValues: %+v", nestedValues)
			for _, subfield := range subfields {
				log.Printf("subfield: %s", subfield)
				lastItem := nestedValues[len(nestedValues)-1]
				//log.Printf("value: %+v", lastItem[subfield])
				switch val := lastItem[subfield].(type) {
				case map[string]any:
					nestedValues = append(nestedValues, val)
				default:
					formattedValue = fmt.Sprintf("%v", val)
				}
			}

			if reflect.TypeOf(formattedValue).Kind() == reflect.Slice {
				marshaledSlice, _ := json.MarshalIndent(formattedValue, "", "  ")
				formattedValue = string(marshaledSlice)
			}

			row = append(row, formattedValue)
		}
		t.AppendRow(row)
	}

	t.SetStyle(table.StyleLight)
	cmd.Println(t.Render())
}

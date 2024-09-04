package output

import (
	"bytes"
	"encoding/json"
	"errors"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/neo4j/cli/common/clictx"
	"github.com/spf13/cobra"
)

// Prints a response body
func PrintBody(cmd *cobra.Command, body []byte, fields []string) error {
	config, ok := clictx.Config(cmd.Context())
	if !ok {
		return errors.New("error fetching cli configuration values")
	}

	outputType, err := config.GetString("aura.output")
	if err != nil {
		return err
	}

	if len(body) > 0 {
		switch output := outputType; output {
		case "json":
			var pretty bytes.Buffer
			err := json.Indent(&pretty, body, "", "\t")
			if err != nil {
				return err
			}
			cmd.Println(pretty.String())
		case "table", "default":
			err := PrintTable(cmd, body, fields)
			if err != nil {
				return err
			}

		default:
			// This is in case the value is unknown
			cmd.Println(string(body))
		}
	}

	return nil
}

func PrintTable(cmd *cobra.Command, body []byte, fields []string) error {
	values, err := parseBody(body)
	if err != nil {
		return err
	}

	t := table.NewWriter()

	header := table.Row{}
	for _, f := range fields {
		header = append(header, f)
	}

	t.AppendHeader(header)
	for _, v := range values {
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

func parseBody(body []byte) ([]map[string]any, error) {
	var values []map[string]any
	var jsonWithArray struct{ Data []map[string]any }

	err := json.Unmarshal(body, &jsonWithArray)

	// Try unmarshalling array first, if not it creates an array from the single item
	if err == nil {
		values = jsonWithArray.Data
	} else {
		var jsonWithSingleItem struct{ Data map[string]any }
		err := json.Unmarshal(body, &jsonWithSingleItem)
		if err != nil {
			return nil, err
		}
		values = []map[string]any{jsonWithSingleItem.Data}
	}

	return values, nil
}

package output

import (
	"bytes"
	"encoding/json"
	"errors"

	"github.com/neo4j/cli/common/clictx"
	"github.com/spf13/cobra"
)

// Prints a response body
func PrintBody(cmd *cobra.Command, body []byte) error {
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
		default:
			cmd.Println(string(body))
		}
	}

	return nil
}

package output

import (
	"bytes"
	"encoding/json"

	"github.com/neo4j/cli/common/clicfg"
	"github.com/spf13/cobra"
)

// Prints a response body
func PrintBody(cmd *cobra.Command, cfg *clicfg.Config, body []byte) error {
	outputType := cfg.Aura.Output()

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

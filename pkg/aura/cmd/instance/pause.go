package instance

import (
	"fmt"

	"github.com/neo4j/cli/pkg/aura/api"
	"github.com/spf13/cobra"
)

func NewPauseCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "pause",
		Short: "Pauses an instance",
		Long: `Starts the pause process of an Aura instance.

Pausing an instance is an asynchronous operation. You can poll the current status of this operation by periodically getting the instance details for the instance ID using the get subcommand.

The pause time depends on the amount of data stored in the instance; larger quantities of data will take longer. The exact time this will take is dependent on the size of your data store.

If another operation is being performed on the instance you are trying to pause, an error will be returned that indicates that the pause operation cannot be performed.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return api.MakeRequest(cmd, "POST", fmt.Sprintf("/instances/%s/pause", args[0]), nil)
		},
	}
}

package instance

import (
	"fmt"

	"github.com/neo4j/cli/pkg/aura/api"
	"github.com/spf13/cobra"
)

func NewResumeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "resume",
		Short: "Resumes an instance",
		Long: `Starts the resume process of an Aura instance.

Resuming an instance is an asynchronous operation. You can poll the current status of this operation by periodically getting the instance details for the instance ID using the get subcommand.

If another operation is being performed on the instance you are trying to resume, an error will be returned that indicates that resume cannot be performed.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return api.MakeRequest(cmd, "POST", fmt.Sprintf("/instances/%s/resume", args[0]), nil)
		},
	}
}

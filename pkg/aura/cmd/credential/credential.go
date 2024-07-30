package credential

import (
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "credential",
		Short: "Manage and view credential values",
	}

	cmd.AddCommand(NewAddCmd())
	cmd.AddCommand(NewRemoveCmd())
	cmd.AddCommand(NewUseCmd())

	return cmd
}

package credential

import (
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "credential",
	Short: "Manage and view credential values",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
}

func init() {
	Cmd.AddCommand(AddCmd)
	Cmd.AddCommand(RemoveCmd)
	Cmd.AddCommand(UseCmd)
}

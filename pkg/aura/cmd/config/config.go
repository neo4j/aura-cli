package config

import (
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage and view configuration values",
		Long: `A longer description that spans multiple lines and likely contains examples
	and usage of using your command. For example:
	
	Cobra is a CLI library for Go that empowers applications.
	This application is a tool to generate the needed files
	to quickly create a Cobra application.`,
	}

	cmd.AddCommand(NewGetCmd())
	cmd.AddCommand(NewListCmd())
	cmd.AddCommand(NewSetCmd())

	return cmd
}

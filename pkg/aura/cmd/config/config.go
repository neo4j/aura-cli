package config

import (
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage and view configuration values",
	}

	cmd.AddCommand(NewGetCmd())
	cmd.AddCommand(NewListCmd())
	cmd.AddCommand(NewSetCmd())

	return cmd
}

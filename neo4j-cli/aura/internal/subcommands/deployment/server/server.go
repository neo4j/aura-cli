package server

import (
	"github.com/neo4j/cli/common/clicfg"
	serverdatabase "github.com/neo4j/cli/neo4j-cli/aura/internal/subcommands/deployment/server/database"
	"github.com/spf13/cobra"
)

func NewCmd(cfg *clicfg.Config) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "server",
		Short: "Relates to deployment servers",
	}

	cmd.AddCommand(NewListCmd(cfg))
	cmd.AddCommand(serverdatabase.NewCmd(cfg))

	return cmd
}

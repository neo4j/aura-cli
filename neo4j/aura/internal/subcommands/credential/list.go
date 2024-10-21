package credential

import (
	"fmt"

	"github.com/neo4j/cli/common/clicfg"
	"github.com/spf13/cobra"
)

func NewListCmd(cfg *clicfg.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "list credentials",
		RunE: func(cmd *cobra.Command, args []string) error {
			creds := cfg.Credentials.Aura.List()

			fmt.Println(creds)
			return nil
		},
	}
}

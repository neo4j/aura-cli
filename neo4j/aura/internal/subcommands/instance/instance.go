package instance

import (
	"github.com/neo4j/cli/common/clicfg"
	"github.com/neo4j/cli/neo4j/aura/internal/subcommands/instance/snapshot"
	"github.com/spf13/cobra"
)

func NewCmd(cfg *clicfg.Config) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "instance",
		Short: "Relates to AuraDB or AuraDS instances",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := cfg.Aura.BindBaseUrl(cmd.Flags().Lookup("base-url")); err != nil {
				return err
			}
			if err := cfg.Aura.BindAuthUrl(cmd.Flags().Lookup("auth-url")); err != nil {
				return err
			}
			if err := cfg.Aura.BindOutput(cmd.Flags().Lookup("output")); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.AddCommand(NewCreateCmd(cfg))
	cmd.AddCommand(NewDeleteCmd(cfg))
	cmd.AddCommand(NewGetCmd(cfg))
	cmd.AddCommand(NewListCmd(cfg))
	cmd.AddCommand(NewPauseCmd(cfg))
	cmd.AddCommand(NewResumeCmd(cfg))
	cmd.AddCommand(NewUpdateCmd(cfg))
	cmd.AddCommand(NewOverwriteCmd(cfg))
	cmd.AddCommand(snapshot.NewCmd(cfg))

	cmd.PersistentFlags().String("auth-url", "", "")
	cmd.PersistentFlags().String("base-url", "", "")
	cmd.PersistentFlags().String("output", "", "")

	return cmd
}

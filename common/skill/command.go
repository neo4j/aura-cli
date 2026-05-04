// Copyright (c) "Neo4j"
// Neo4j Sweden AB [http://neo4j.com]

package skill

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"

	"github.com/neo4j/cli/common/clicfg"
	"github.com/neo4j/cli/common/clierr"
)

// NewCmd builds the per-binary `skill` cobra command tree. `bundle` is the
// embedded SKILL.md + references/ tree for the calling binary; `skillName`
// is the lowercase id used as the on-disk install dir
// (`<agentSkillsDir>/<skillName>/`) and the SKILL.md frontmatter `name:`.
//
// Each leaf renders results as a table by default and emits a JSON envelope
// when `--output json` is passed (matching neo4j-cli/aura/internal/output
// conventions).
func NewCmd(cfg *clicfg.Config, bundle fs.FS, skillName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "skill",
		Short: "Install agent skills for this CLI into supported AI agents",
		Long: "Install, remove, list, and check the per-binary agent-skill " +
			"bundle. The bundle teaches AI agents (Claude Code, Cursor, " +
			"Windsurf, etc.) how to drive this CLI.",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			outputValue := cmd.Flags().Lookup("output").Value.String()
			if outputValue != "" {
				valid := false
				for _, v := range clicfg.ValidOutputValues {
					if v == outputValue {
						valid = true
						break
					}
				}
				if !valid {
					return clierr.NewUsageError("invalid output value specified: %s", outputValue)
				}
			}
			cfg.Aura.BindOutput(cmd.Flags().Lookup("output"))
			return nil
		},
	}

	cmd.PersistentFlags().String("output", "", fmt.Sprintf("Format to print console output in, from a choice of [%s]", strings.Join(clicfg.ValidOutputValues[:], ", ")))

	cmd.AddCommand(newInstallCmd(cfg, bundle, skillName))
	cmd.AddCommand(newRemoveCmd(cfg, skillName))
	cmd.AddCommand(newListCmd(cfg, skillName))
	cmd.AddCommand(newCheckCmd(cfg, skillName))

	return cmd
}

func newInstallCmd(cfg *clicfg.Config, bundle fs.FS, skillName string) *cobra.Command {
	return &cobra.Command{
		Use:   "install [agent]",
		Short: "Install the skill bundle into supported AI agents",
		Long: "Without an argument, installs into every detected agent. " +
			"With an [agent] argument (case-insensitive), installs into that " +
			"single agent. Unknown agent names exit non-zero with the list " +
			"of valid names.",
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			filter := ""
			if len(args) == 1 {
				filter = args[0]
			}
			targets, err := Install(cfg.Aura.Fs(), bundle, skillName, cfg.Version, filter)
			if err != nil {
				return formatAgentErr(err)
			}
			renderInstallResult(cmd, cfg, skillName, "installed", targets)
			return nil
		},
	}
}

func newRemoveCmd(cfg *clicfg.Config, skillName string) *cobra.Command {
	return &cobra.Command{
		Use:   "remove [agent]",
		Short: "Remove the installed skill bundle",
		Long: "Without an argument, removes from every detected agent. " +
			"With an [agent] argument (case-insensitive), removes from that " +
			"single agent. Idempotent: a second run on a clean target is a " +
			"no-op.",
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			filter := ""
			if len(args) == 1 {
				filter = args[0]
			}
			targets, err := Remove(cfg.Aura.Fs(), skillName, filter)
			if err != nil {
				return formatAgentErr(err)
			}
			renderInstallResult(cmd, cfg, skillName, "removed", targets)
			return nil
		},
	}
}

func newListCmd(cfg *clicfg.Config, skillName string) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List supported agents and per-agent install state",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			rows, err := List(cfg.Aura.Fs(), skillName)
			if err != nil {
				return err
			}
			renderListResult(cmd, cfg, rows)
			return nil
		},
	}
}

func newCheckCmd(cfg *clicfg.Config, skillName string) *cobra.Command {
	return &cobra.Command{
		Use:   "check",
		Short: "Check installed skills for version drift against this binary",
		Long: "Reads each installed SKILL.md frontmatter `version:` and " +
			"compares to the running binary version. Exits non-zero on any " +
			"drift; prints a per-agent table.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			rows, drift, err := Check(cfg.Aura.Fs(), skillName, cfg.Version)
			if err != nil {
				return err
			}
			renderCheckResult(cmd, cfg, rows)
			if drift {
				return fmt.Errorf("skill: drift detected — run `skill install` to refresh")
			}
			return nil
		},
	}
}

// formatAgentErr converts skill-package sentinel errors into user-facing
// usage errors that include the valid agent names. Other errors pass
// through unchanged.
func formatAgentErr(err error) error {
	switch {
	case errors.Is(err, ErrUnknownAgent):
		return clierr.NewUsageError("%v\nvalid agents: %s", err, strings.Join(agentNames(), ", "))
	case errors.Is(err, ErrAgentNotDetected):
		return clierr.NewUsageError("%v", err)
	case errors.Is(err, ErrNoAgentsDetected):
		return clierr.NewUsageError("%v\nvalid agents: %s", err, strings.Join(agentNames(), ", "))
	default:
		return err
	}
}

func agentNames() []string {
	names := make([]string, 0, len(AGENTS))
	for i := range AGENTS {
		names = append(names, AGENTS[i].Name)
	}
	return names
}

// installResultRow is the JSON shape emitted by install/remove.
type installResultRow struct {
	Agent       string `json:"agent"`
	DisplayName string `json:"display_name"`
	SkillsPath  string `json:"skills_path"`
	Action      string `json:"action"`
}

// renderInstallResult prints the install/remove outcome as a table or
// JSON. `action` is "installed" or "removed" — printed in the Action
// column / JSON field. Empty target list emits a friendly note in table
// mode and an empty array in JSON mode.
func renderInstallResult(cmd *cobra.Command, cfg *clicfg.Config, skillName, action string, targets []*Agent) {
	rows := make([]installResultRow, 0, len(targets))
	for _, a := range targets {
		sp, _ := a.SkillsPath()
		var path string
		if sp != "" {
			path = sp + "/" + skillName
		}
		rows = append(rows, installResultRow{
			Agent:       a.Name,
			DisplayName: a.DisplayName,
			SkillsPath:  path,
			Action:      action,
		})
	}

	if cfg.Aura.Output() == "json" {
		printJSON(cmd, rows)
		return
	}

	if len(rows) == 0 {
		cmd.Printf("No agents to %s.\n", strings.TrimSuffix(action, "ed"))
		return
	}

	t := table.NewWriter()
	t.AppendHeader(table.Row{"agent", "display", "path", "action"})
	for _, r := range rows {
		t.AppendRow(table.Row{r.Agent, r.DisplayName, r.SkillsPath, r.Action})
	}
	t.SetStyle(table.StyleLight)
	cmd.Println(t.Render())
}

// listResultRow is the JSON shape emitted by list.
type listResultRow struct {
	Agent            string `json:"agent"`
	DisplayName      string `json:"display_name"`
	Detected         bool   `json:"detected"`
	Installed        bool   `json:"installed"`
	InstalledVersion string `json:"installed_version"`
}

func renderListResult(cmd *cobra.Command, cfg *clicfg.Config, rows []AgentInstall) {
	out := make([]listResultRow, 0, len(rows))
	for _, r := range rows {
		out = append(out, listResultRow{
			Agent:            r.Agent.Name,
			DisplayName:      r.Agent.DisplayName,
			Detected:         r.Detected,
			Installed:        r.Installed,
			InstalledVersion: r.InstalledVersion,
		})
	}

	if cfg.Aura.Output() == "json" {
		printJSON(cmd, out)
		return
	}

	t := table.NewWriter()
	t.AppendHeader(table.Row{"agent", "display", "detected", "installed", "installed-version"})
	for _, r := range out {
		t.AppendRow(table.Row{r.Agent, r.DisplayName, boolStr(r.Detected), boolStr(r.Installed), r.InstalledVersion})
	}
	t.SetStyle(table.StyleLight)
	cmd.Println(t.Render())
}

// checkResultRow is the JSON shape emitted by check.
type checkResultRow struct {
	Agent            string `json:"agent"`
	InstalledVersion string `json:"installed_version"`
	CurrentVersion   string `json:"current_version"`
	Status           string `json:"status"`
}

func renderCheckResult(cmd *cobra.Command, cfg *clicfg.Config, rows []CheckRow) {
	out := make([]checkResultRow, 0, len(rows))
	for _, r := range rows {
		out = append(out, checkResultRow{
			Agent:            r.Agent.Name,
			InstalledVersion: r.InstalledVersion,
			CurrentVersion:   r.CurrentVersion,
			Status:           r.Status,
		})
	}

	if cfg.Aura.Output() == "json" {
		printJSON(cmd, out)
		return
	}

	if len(out) == 0 {
		cmd.Println("No installed skills found.")
		return
	}

	t := table.NewWriter()
	t.AppendHeader(table.Row{"agent", "installed-version", "current-version", "status"})
	for _, r := range out {
		t.AppendRow(table.Row{r.Agent, r.InstalledVersion, r.CurrentVersion, r.Status})
	}
	t.SetStyle(table.StyleLight)
	cmd.Println(t.Render())
}

func printJSON(cmd *cobra.Command, v any) {
	bytes, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		// Marshalling our own structs cannot fail in practice; mirror the
		// existing output package's posture (panic on impossible-state).
		panic(err)
	}
	cmd.Println(string(bytes))
}

func boolStr(b bool) string {
	if b {
		return "yes"
	}
	return "no"
}

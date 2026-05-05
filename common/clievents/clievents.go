package clievents

import (
	"fmt"
	"strings"

	"github.com/neo4j/cli/common/analytics"
	"github.com/spf13/pflag"
)

func Emit(events *analytics.Analytics, args []string, state bool) {
	cliCommand := strings.Trim(fmt.Sprint(args), "[]")

	flags := pflag.NewFlagSet("cliEvents", pflag.ContinueOnError)
	var output string
	var help bool

	flags.StringVar(&output, "output", "", "")
	flags.BoolVarP(&help, "help", "h", false, "")

	_ = flags.Parse(args)

	if help {
		events.EmitHelpEvent(cliCommand)
	} else {
		events.EmitCommandEvent(cliCommand, state)
	}
}

package step

import (
	"fmt"
	"strings"
	"text/tabwriter"
)

// Kind represents the type of operation in a profile execution plan.
type Kind string

const (
	// DotEnv indicates a dotenv step.
	DotEnv Kind = "dotenv"
	// Profile indicates a profile step.
	Profile Kind = "env"
	// Overlay indicates an overlay step.
	Overlay Kind = "overlay"
)

// Step represents a single operation in a profile execution plan.
type Step struct {
	// Kind specifies the type of step operation.
	Kind Kind
	// Owner identifies the profile that owns this step.
	Owner string
	// Name contains the profile name or dotenv path.
	Name string
}

// Steps represents a sequence of profile execution steps.
type Steps []Step

// Table formats the steps as a human-readable table.
func (s Steps) Table() string {
	var builder strings.Builder

	//nolint:mnd 	// minwidth, tabwidth, padding, padchar, flags
	writer := tabwriter.NewWriter(&builder, 0, 4, 2, ' ', 0)

	_, _ = fmt.Fprintln(writer, "STEP\tPROFILE\tKIND\tNAME")

	for number, step := range s {
		if step.Owner == "" {
			step.Owner = step.Name // for env steps, owner == profile
			step.Name = ""         // no name for env steps
		}

		fmt.Fprintf(writer, "%02d\t%s\t%s\t%s\n", number+1, step.Owner, step.Kind, step.Name)
	}

	_ = writer.Flush()

	return builder.String()
}

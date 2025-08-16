package cli

import (
	"fmt"
	"strings"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"
)

// List returns the cobra command for listing profiles and their variables.
//
//nolint:forbidigo	// Command prints out to the console.
func List(options *Options) *cobra.Command {
	var oneline bool

	cmd := &cobra.Command{
		Use:   "list [key]",
		Short: "List profiles and their variables",
		Long: heredoc.Doc(`
			List all variables for a specific profile,
			or the value of a variable for a specific profile.
		`),
		Example: heredoc.Doc(`
			# List all variables for 'dev' with sources
			$ envprof --profile dev -v list

			# Show the value of HOST in 'dev'
			$ envprof --profile dev list HOST
		`),
		Aliases: []string{"ls"},
		Args:    cobra.MaximumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			nArgs := len(args)

			vars, err := loadProfile(options.EnvProf, options.Profile)
			if err != nil {
				return err
			}

			if oneline {
				options.Verbose = false
			}

			var output string

			if nArgs > 0 {
				variable := args[0]
				if !vars.Env.Exists(variable) {
					return fmt.Errorf("key %q not found in profile %q", variable, vars.Name)
				}

				output = vars.Format(variable, options.Verbose, false)
			} else {
				output = vars.FormatAll("", options.Verbose)
			}

			if oneline {
				output = strings.Join(strings.Fields(output), " ")
			}

			fmt.Println(output)

			return nil
		},
	}

	cmd.Flags().SortFlags = false

	cmd.Flags().BoolVarP(&oneline, "oneline", "o", false, "Emit variables on a single line")

	return cmd
}

package cli

import (
	"fmt"
	"strings"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"github.com/idelchi/envprof/internal/environment"
)

// List returns the cobra command for listing profiles and their variables.
//
//nolint:forbidigo	// Command prints out to the console.
func List(options *Options) *cobra.Command {
	var oneline bool
	var planOnly bool

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
			if planOnly {
				steps, err := loadPlan(options)
				if err != nil {
					return err
				}

				fmt.Println(steps.Chain())

				return nil
			}

			nArgs := len(args)

			env, err := loadProfile(options)
			if err != nil {
				return err
			}

			var output string

			const padding = 60

			formatter := environment.Formatter{
				// oneline implies not verbose
				WithOrigin: options.Verbose && !oneline,
				WithKey:    true,
				Padding:    padding,
			}

			if nArgs == 1 {
				variable := args[0]
				if !env.Env.Exists(variable) {
					return fmt.Errorf("key %q not found in profile %q", variable, env.Name)
				}

				formatter.WithKey = false

				output = formatter.Key(variable, env)
			} else {
				output = formatter.All(env)
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
	cmd.Flags().BoolVarP(&planOnly, "plan-only", "p", false, "Show only the plan for the profile")

	return cmd
}

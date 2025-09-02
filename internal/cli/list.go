package cli

import (
	"fmt"
	"strings"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"github.com/idelchi/envprof/internal/environment"
)

// List returns the cobra command for listing profiles and their variables.
func List(options *Options) *cobra.Command {
	var (
		oneline bool
		dry     bool
	)

	cmd := &cobra.Command{
		Use:   "list [key]",
		Short: "List profiles and their variables",
		Long: heredoc.Doc(`
			List all variables for a specific profile,
			or the value of a variable for a specific profile.
		`),
		Example: heredoc.Doc(`
			# List all variables for 'dev' with sources
			envprof --profile dev -v list

			# Show the value of HOST in 'dev'
			envprof --profile dev list HOST

			# List the layering order only
			envprof --profile dev list --dry
		`),
		Aliases: []string{"ls"},
		Args: func(cmd *cobra.Command, args []string) error {
			if err := cobra.MaximumNArgs(1)(cmd, args); err != nil {
				return fmt.Errorf(
					"%q only accepts [key] as an optional argument, received %d arguments: %v",
					cmd.Name(),
					len(args),
					args,
				)
			}

			return nil
		},
		RunE: func(_ *cobra.Command, args []string) error {
			profiles, profile, steps, err := LoadPlan(options)
			if err != nil {
				return err
			}

			if dry {
				//nolint:forbidigo	// Command prints out to the console.
				fmt.Println(steps.Table())

				return nil
			}

			env, err := profiles.Environment(profile, steps)
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

			if len(args) == 1 {
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

			//nolint:forbidigo	// Command prints out to the console.
			fmt.Println(output)

			return nil
		},
	}

	cmd.Flags().SortFlags = false

	cmd.Flags().BoolVarP(&oneline, "oneline", "o", false, "Emit variables on a single line")
	cmd.Flags().BoolVarP(&dry, "dry", "d", false, "Show only the plan for the profile")

	return cmd
}

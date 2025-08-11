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
func List(envprof *[]string) *cobra.Command {
	var (
		verbose bool
		oneline bool
	)

	cmd := &cobra.Command{
		Use:   "list [profile] [key]",
		Short: "List profiles and their variables",
		Long: heredoc.Doc(`
			Lists all profiles (sorted),
			all variables for a specific profile,
			or the value of a variable for a specific profile.
		`),
		Example: heredoc.Doc(`
			# List all profiles
			$ envprof list

			# List all variables for 'dev' with sources
			$ envprof list dev -v

			# Show the value of HOST in 'dev'
			$ envprof list dev HOST
		`),
		Aliases: []string{"ls"},
		Args:    cobra.MaximumNArgs(2), //nolint:mnd	// The command takes up to 2 arguments as documented.
		RunE: func(_ *cobra.Command, args []string) error {
			if len(args) == 0 {
				profiles, err := load(*envprof)
				if err != nil {
					return err
				}
				for _, profile := range profiles.Names() {
					fmt.Println(profile)
				}

				return nil
			}

			prof := args[0]

			vars, err := loadProfileVars(*envprof, prof)
			if err != nil {
				return err
			}

			if oneline {
				verbose = false
			}

			var output string

			if len(args) > 1 {
				if !vars.Env.Exists(args[1]) {
					//nolint:err113	// Occasional dynamic errors are fine.
					return fmt.Errorf("key %q not found in profile %q", args[1], prof)
				}

				output = vars.Format(args[1], verbose, false)
			} else {
				output = vars.FormatAll("", verbose)
			}

			if oneline {
				output = strings.Join(strings.Fields(output), " ")
			}

			fmt.Println(output)

			return nil
		},
	}

	cmd.Flags().SortFlags = false

	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Show from which source each variable is inherited")
	cmd.Flags().BoolVarP(&oneline, "oneline", "o", false, "Emit variables on a single line")

	return cmd
}

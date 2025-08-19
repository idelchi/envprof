package cli

import (
	"fmt"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"github.com/idelchi/envprof/internal/environment"
)

// Diff returns the cobra command for diffing profiles.
func Diff(options *Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "diff <profile>",
		Short: "Diff profiles",
		Long: heredoc.Doc(`
			Compare the specified profile with the currently loaded profile.

			Outputs changes in a diff-like format:
			 - KEY="VALUE"   means the key was removed
			 + KEY="VALUE"   means the key was added
			 ~ KEY: "OLD" -> "NEW"   means the key changed
		`),
		Args: func(cmd *cobra.Command, args []string) error {
			if err := cobra.ExactArgs(1)(cmd, args); err != nil {
				return fmt.Errorf(
					"%q requires a <profile> as it's only positional argument, received %d arguments: %v",
					cmd.Name(),
					len(args),
					args,
				)
			}

			return nil
		},
		RunE: func(_ *cobra.Command, args []string) error {
			env1, err := LoadProfile(options)
			if err != nil {
				return err
			}

			options.Profile = args[0]

			env2, err := LoadProfile(options)
			if err != nil {
				return err
			}

			diff := environment.Diffs(env1.Env, env2.Env)

			if diff.Equal() {
				//nolint:forbidigo	// Command prints out to the console.
				fmt.Println("No differences found.")

				return nil
			}

			if err := diff.Render(env1.Name, env2.Name); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().SortFlags = false

	return cmd
}

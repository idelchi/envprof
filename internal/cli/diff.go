package cli

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/idelchi/envprof/internal/environment"
)

// Diff returns the cobra command for diffing profiles.
//
//nolint:forbidigo	// Command prints out to the console.
func Diff(options *Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "diff <profile>",
		Short: "Diff profiles",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			env1, err := loadProfile(options)
			if err != nil {
				return err
			}

			options.Profile = args[0]

			env2, err := loadProfile(options)
			if err != nil {
				return err
			}

			environment.Diffs(env1.Env, env2.Env).
				RenderUnified(os.Stdout, env1.Name, env2.Name, environment.RenderOptions{Color: options.Verbose})

			return nil
		},
	}

	cmd.Flags().SortFlags = false

	return cmd
}

package cli

import (
	"fmt"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"github.com/idelchi/envprof/pkg/terminal"
	"github.com/idelchi/godyl/pkg/env"
	"github.com/idelchi/godyl/pkg/path/file"
)

// Shell returns the cobra command for entering a scoped shell with the active environment.
func Shell(options *Options) *cobra.Command {
	environment := env.FromEnv()

	var (
		shell   = environment.GetAny("SHELL", "STARSHIP_SHELL")
		isolate bool
		path    bool
		envs    []string
	)

	cmd := &cobra.Command{
		Use:   "shell",
		Short: "Spawn a subshell with a profile",
		Long: heredoc.Doc(`
			Launch a new shell with the selected profile's environment.

			Customize shell and level of isolation with --shell, --isolate, and --path.
		`),
		Example: heredoc.Doc(`
			# Subshell with profile
			envprof --profile dev shell

			# Isolated subshell
			envprof --profile dev shell --isolate

			# Isolated but keep PATH
			envprof --profile dev shell --isolate --path

			# Use a specific shell
			envprof --profile dev shell --shell zsh
		`),
		Aliases: []string{"sh"},
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if active := environment.Get("ENVPROF_ACTIVE_PROFILE"); environment.Exists(
				"ENVPROF_ACTIVE_PROFILE",
			) {
				return fmt.Errorf(
					"already inside profile %q, nested profiles are not allowed, please exit first",
					active,
				)
			}

			profile, err := LoadProfile(options)
			if err != nil {
				return err
			}

			if err = profile.Env.AddPair("ENVPROF_ACTIVE_PROFILE", profile.Name); err != nil {
				return err
			}

			profile.Env = Merge(profile.Env, environment, isolate, path, envs)

			if shell == "" {
				shell = string(terminal.Current())
			}

			//nolint:forbidigo	// Command prints out to the console.
			fmt.Printf(
				"Entering shell %q with profile %q...\n",
				file.New(shell).WithoutExtension().Base(),
				profile.Name,
			)

			if err := terminal.Spawn(shell, profile.Env.AsSlice()); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().
		StringVarP(&shell, "shell", "s", shell, "Shell to launch (leave empty to auto-detect)")
	cmd.Flags().BoolVarP(&isolate, "isolate", "i", false, "Isolate from parent environment")
	cmd.Flags().BoolVarP(&path, "path", "p", false, "Include the current PATH in the environment")
	cmd.Flags().
		StringSliceVarP(&envs, "env", "e", nil, "Passthrough environment variables (combined with --isolate)")

	return cmd
}

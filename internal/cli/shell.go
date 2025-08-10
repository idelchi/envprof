package cli

import (
	"fmt"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"github.com/idelchi/envprof/internal/terminal"
	"github.com/idelchi/godyl/pkg/env"
	"github.com/idelchi/godyl/pkg/path/file"
)

// Shell returns the cobra command for entering a scoped shell with the environment active.
//
//nolint:forbidigo	// Command prints out to the console.
func Shell(envprof *[]string) *cobra.Command {
	env := env.FromEnv()

	var (
		shell   = env.GetAny("SHELL", "STARSHIP_SHELL")
		isolate bool
	)

	cmd := &cobra.Command{
		Use:   "shell <profile>",
		Short: "Spawn a new shell with the profile's environment",
		Long: heredoc.Doc(`
      Spawn a new shell with the profile's environment.
      The shell will be spawned with the environment variables set to the profile's values.

      Use the --shell flag to specify the shell to spawn, otherwise it will try to identify the current shell.

      Use the --inherit flag to inherit the environment variables from the parent shell.
    `),
		Example: heredoc.Doc(`
      # Spawn a new shell with the profile's environment
      $ godyl shell dev

      # Spawn a new shell with the profile's environment and inherit the environment variables from the parent shell
      $ godyl shell dev --inherit

      # Spawn a new shell with the profile's environment and use zsh as the shell
      $ godyl shell dev --shell zsh
      `),
		Aliases: []string{"sh"},
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if active := env.Get("ENVPROF_ACTIVE_PROFILE"); env.Exists("ENVPROF_ACTIVE_PROFILE") {
				//nolint:err113	// Occasional dynamic errors are fine.
				return fmt.Errorf(
					"already inside profile %q, nested profiles are not allowed, please exit first",
					active,
				)
			}

			profiles, err := load(*envprof)
			if err != nil {
				return err
			}

			if shell == "" {
				shell = terminal.Current()
			}

			prof := args[0]

			vars, err := profiles.Environment(prof)
			if err != nil {
				return err //nolint:wrapcheck	// Error does not need additional wrapping.
			}

			if err = vars.Env.AddPair("ENVPROF_ACTIVE_PROFILE", prof); err != nil {
				return err //nolint:wrapcheck	// Error does not need additional wrapping.
			}

			if !isolate {
				vars.Env.Merge(env)
			}

			fmt.Printf("Entering shell %q with profile %q...\n", file.New(shell).WithoutExtension().Base(), prof)

			if err := terminal.Spawn(shell, vars.Env.AsSlice()); err != nil {
				return err //nolint:wrapcheck	// Error does not need additional wrapping.
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&shell, "shell", "s", shell, "Shell to launch (leave empty to auto-detect).")
	cmd.Flags().BoolVarP(&isolate, "isolate", "i", false, "Isolate from parent environment.")

	return cmd
}

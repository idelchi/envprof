package cli

import (
	"fmt"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"github.com/idelchi/envprof/internal/terminal"
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

			prof, err := LoadProfile(options)
			if err != nil {
				return err
			}

			if err = prof.Env.AddPair("ENVPROF_ACTIVE_PROFILE", prof.Name); err != nil {
				return err
			}

			if !isolate {
				prof.Env.Merge(environment)
			} else if path {
				prof.Env.Merge(env.Env{"PATH": environment.Get("PATH")})
			}

			if shell == "" {
				shell = terminal.Current()
			}

			//nolint:forbidigo	// Command prints out to the console.
			fmt.Printf(
				"Entering shell %q with profile %q...\n",
				file.New(shell).WithoutExtension().Base(),
				prof.Name,
			)

			if err := terminal.Spawn(shell, prof.Env.AsSlice()); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().
		StringVarP(&shell, "shell", "s", shell, "Shell to launch (leave empty to auto-detect).")
	cmd.Flags().BoolVarP(&isolate, "isolate", "i", false, "Isolate from parent environment.")
	cmd.Flags().BoolVarP(&path, "path", "p", false, "Include the current PATH in the environment.")

	return cmd
}

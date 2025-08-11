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
//
//nolint:forbidigo	// Command prints out to the console.
func Shell(envprof *[]string) *cobra.Command {
	environment := env.FromEnv()

	var (
		shell   = environment.GetAny("SHELL", "STARSHIP_SHELL")
		isolate bool
		path    bool
	)

	cmd := &cobra.Command{
		Use:   "shell <profile>",
		Short: "Spawn a subshell with a profile",
		Long: heredoc.Doc(`
			Launch a new shell with the selected profile's environment.

			Customize shell and level of isolation with --shell, --isolate, and --path.
		`),
		Example: heredoc.Doc(`
			# Subshell with profile
			envprof shell dev

			# Isolated subshell
			envprof shell dev --isolate

			# Isolated but keep PATH
			envprof shell dev --isolate --path

			# Use a specific shell
			envprof shell dev --shell zsh
		`),
		Aliases: []string{"sh"},
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if active := environment.Get("ENVPROF_ACTIVE_PROFILE"); environment.Exists("ENVPROF_ACTIVE_PROFILE") {
				//nolint:err113	// Occasional dynamic errors are fine.
				return fmt.Errorf(
					"already inside profile %q, nested profiles are not allowed, please exit first",
					active,
				)
			}

			prof := args[0]

			profEnv, err := loadProfileEnv(*envprof, prof)
			if err != nil {
				return err
			}

			if err = profEnv.AddPair("ENVPROF_ACTIVE_PROFILE", prof); err != nil {
				return err //nolint:wrapcheck	// Error does not need additional wrapping.
			}

			if !isolate {
				profEnv.Merge(environment)
			} else if path {
				profEnv.Merge(env.Env{"PATH": environment.Get("PATH")})
			}

			if shell == "" {
				shell = terminal.Current()
			}

			fmt.Printf("Entering shell %q with profile %q...\n", file.New(shell).WithoutExtension().Base(), prof)

			if err := terminal.Spawn(shell, profEnv.AsSlice()); err != nil {
				return err //nolint:wrapcheck	// Error does not need additional wrapping.
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&shell, "shell", "s", shell, "Shell to launch (leave empty to auto-detect).")
	cmd.Flags().BoolVarP(&isolate, "isolate", "i", false, "Isolate from parent environment.")
	cmd.Flags().BoolVarP(&path, "path", "p", false, "Include the current PATH in the environment.")

	return cmd
}

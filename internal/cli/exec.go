package cli

import (
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	execx "github.com/idelchi/envprof/internal/exec"
	"github.com/idelchi/godyl/pkg/env"
)

// Exec returns the cobra command for executing a command with the active environment.
func Exec(envprof *[]string) *cobra.Command {
	environment := env.FromEnv()

	var (
		isolate bool
		path    bool
	)

	cmd := &cobra.Command{
		Use:   "exec <profile> <command> [args...]",
		Short: "Execute a command with a profile",
		Long: heredoc.Doc(`
			Run a command with the selected profile's environment.

			On Unix, replaces the current process.
			On Windows, runs and exits with the same code.
    	`),
		Example: heredoc.Doc(`
			# Run a command with 'dev'
			envprof exec dev -- make build

			# Isolated exec
			envprof exec dev --isolate -- npm run test

			# Isolated exec but keep PATH
			envprof exec dev --isolate --path -- python -V
      	`),
		Aliases: []string{"ex"},
		Args:    cobra.MinimumNArgs(2), //nolint:mnd	// The command a minimum of 2 arguments as documented.
		RunE: func(_ *cobra.Command, args []string) error {
			prof := args[0]

			profEnv, err := loadProfileEnv(*envprof, prof)
			if err != nil {
				return err
			}

			if !isolate {
				profEnv.Merge(environment)
			} else if path {
				profEnv.Merge(env.Env{"PATH": environment.Get("PATH")})
			}

			cmd := args[1]

			if len(args) > 2 { //nolint:mnd	// The number is obvious from the context.
				args = args[2:]
			} else {
				args = nil
			}

			if err := execx.Replace(cmd, args, profEnv.AsSlice()); err != nil {
				return err //nolint:wrapcheck	// Error does not need additional wrapping.
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&isolate, "isolate", "i", false, "Isolate from parent environment.")
	cmd.Flags().BoolVarP(&path, "path", "p", false, "Include the current PATH in the environment.")

	return cmd
}

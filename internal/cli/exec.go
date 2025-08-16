package cli

import (
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	execx "github.com/idelchi/envprof/internal/exec"
	"github.com/idelchi/godyl/pkg/env"
)

// Exec returns the cobra command for executing a command with the active environment.
func Exec(options *Options) *cobra.Command {
	environment := env.FromEnv()

	var (
		isolate bool
		path    bool
	)

	cmd := &cobra.Command{
		Use:   "exec <command> [args...]",
		Short: "Execute a command with a profile",
		Long: heredoc.Doc(`
			Run a command with the selected profile's environment.

			On Unix, replaces the current process.
			On Windows, runs and exits with the same code.
    	`),
		Example: heredoc.Doc(`
			# Run a command with 'dev'
			envprof --profile dev exec -- make build

			# Isolated exec
			envprof --profile dev exec --isolate -- npm run test

			# Isolated exec but keep PATH
			envprof --profile dev exec --isolate --path -- python --version
      	`),
		Aliases: []string{"ex"},
		Args:    cobra.MinimumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			prof, err := loadProfile(options.EnvProf, options.Profile)
			if err != nil {
				return err
			}

			if !isolate {
				prof.Env.Merge(environment)
			} else if path {
				prof.Env.Merge(env.Env{"PATH": environment.Get("PATH")})
			}

			cmd := args[0]

			if len(args) > 1 {
				args = args[1:]
			} else {
				args = nil
			}

			if err := execx.Replace(cmd, args, prof.Env.AsSlice()); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&isolate, "isolate", "i", false, "Isolate from parent environment.")
	cmd.Flags().BoolVarP(&path, "path", "p", false, "Include the current PATH in the environment.")

	return cmd
}

package cli

import (
	"errors"
	"fmt"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	execx "github.com/idelchi/envprof/pkg/exec"
	"github.com/idelchi/envprof/pkg/terminal"
	"github.com/idelchi/godyl/pkg/env"
)

// Exec returns the cobra command for executing a command with the active environment.
//
//nolint:gocognit,funlen  // stdin addition makes this function slightly complex.
func Exec(options *Options) *cobra.Command {
	environment := env.FromEnv()

	var (
		isolate     bool
		path        bool
		interactive bool
		envs        []string
	)

	cmd := &cobra.Command{
		Use:   "exec <command> [args...]",
		Short: "Execute a command with a profile",
		Long: heredoc.Doc(`
			Run a command with the selected profile's environment.

			On Unix, replaces the current process.
			On Windows, runs and exits with the same code.

			Optionally allows to pass <command> and [args...] via stdin when <command> is "-".
    	`),
		Example: heredoc.Doc(`
			# Run a command with 'dev'
			envprof --profile dev exec -- make build

			# Isolated exec
			envprof --profile dev exec --isolate -- npm run test

			# Isolated exec but keep PATH
			envprof --profile dev exec --isolate --path -- python --version

			# Run in interactive mode (e.g. zsh -i -c "<command> <args...>")
			envprof --profile dev exec --interactive -- <some alias from .zshrc>

			# Run command and arguments passed with stdin
			echo "node --version" | envprof --profile dev exec --interactive -
      	`),
		Aliases: []string{"ex"},
		Args:    cobra.MinimumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			profile, err := LoadProfile(options)
			if err != nil {
				return err
			}

			profile.Env = Merge(profile.Env, environment, isolate, path, envs)

			cmd := args[0]

			if cmd == "-" {
				if ok, err := MaybePiped(); err != nil {
					return err
				} else if !ok {
					return errors.New("no input from stdin")
				}

				args, err = Read()
				if err != nil {
					return err
				}

				cmd = args[0]

				if cmd == "" {
					return errors.New("no input from stdin")
				}
			}

			if len(args) > 1 {
				args = args[1:]
			} else {
				args = nil
			}

			var shell terminal.Shell

			if interactive {
				shell = terminal.Current()

				//nolint:forbidigo	// Command prints out to the console.
				if options.Verbose {
					fmt.Printf("Using login shell: %q\n", shell)
				}
			}

			//nolint:forbidigo	// Command prints out to the console.
			if options.Verbose {
				fmt.Printf("Executing command %q with args %q\n", cmd, args)
			}

			if err := execx.Replace(cmd, args, profile.Env.AsSlice(), shell); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&isolate, "isolate", "i", false, "Isolate from parent environment")
	cmd.Flags().BoolVarP(&path, "path", "p", false, "Include the current PATH in the environment")
	cmd.Flags().BoolVarP(&interactive, "interactive", "I", false, "Run in interactive mode")
	cmd.Flags().
		StringSliceVarP(&envs, "env", "e", nil, "Passthrough environment variables (combined with --isolate)")

	return cmd
}

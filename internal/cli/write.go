package cli

import (
	"errors"
	"fmt"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"github.com/idelchi/envprof/internal/environment"
	"github.com/idelchi/godyl/pkg/path/file"
)

// Write defines the command for writing profile variables to one or multiple dotenv files.
func Write(options *Options) *cobra.Command {
	var all bool

	cmd := &cobra.Command{
		Use:   "write [file]",
		Short: "Write profile variables to dotenv files",
		Long: heredoc.Doc(`
			Write variables to dotenv files.

			Files are written out as <profile>.env, unless set in the profile configuration file, or
			overridden with the [file] argument.
		`),
		Example: heredoc.Doc(`
			# Write 'dev' to dev.env
			envprof --profile dev write

			# Write 'dev' to a specific file
			envprof --profile dev write dev

			# Write all profiles to <profile>.env files
			envprof write --all
		`),
		Aliases: []string{"w"},
		Args:    cobra.RangeArgs(0, 1),
		RunE: func(_ *cobra.Command, args []string) error {
			nArgs := len(args)
			if nArgs > 0 && all {
				return errors.New(
					"'--all' and [file] are incompatible and may not be used together",
				)
			}

			environments, err := environments(all, options, args)
			if err != nil {
				return err
			}

			for _, environment := range environments {
				if err := environment.Write(); err != nil {
					return err
				}

				//nolint:forbidigo		// Command prints out to the console.
				fmt.Printf("Wrote profile %q to %q\n", environment.Name, environment.Output)
			}

			return nil
		},
	}

	cmd.Flags().SortFlags = false

	cmd.Flags().BoolVarP(&all, "all", "a", false, "Write all profiles, ignoring the active profile")

	return cmd
}

// environments returns the selected subset of environments based on the CLI settings.
func environments(all bool, options *Options, args []string) (environments []environment.Environment, err error) {
	switch {
	case all:
		profiles, err := LoadProfiles(options)
		if err != nil {
			return nil, err
		}

		environments, err = profiles.Environments()
		if err != nil {
			return nil, err
		}

	default:
		env, err := LoadProfile(options)
		if err != nil {
			return nil, err
		}

		if len(args) == 1 {
			env.Output = file.New(args[0])
		}

		environments = []environment.Environment{env}
	}

	return environments, nil
}

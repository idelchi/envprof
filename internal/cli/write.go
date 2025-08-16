package cli

import (
	"errors"
	"fmt"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"github.com/idelchi/envprof/internal/profile"
	"github.com/idelchi/godyl/pkg/path/file"
)

// Write defines the command for writing profile variables to one or multiple dotenv files.
//
// TODO(Idelchi): Refactor this function to reduce cognitive complexity.
//
//nolint:forbidigo,gocognit	// Command prints out to the console.
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

			profiles, err := loadProfiles(options.EnvProf)
			if err != nil {
				return err
			}

			var dotenvs []profile.DotEnv
			switch {
			case all:
				profs := profiles.Names()
				dotenvs = make([]profile.DotEnv, 0, len(profs))
				for _, prof := range profs {
					output := file.New(profiles.Output(prof))
					dotenvs = append(dotenvs, profile.DotEnv{Profile: prof, Path: output})
				}
			default:
				prof, err := profileOrDefault(profiles, options.Profile)
				if err != nil {
					return err
				}

				output := file.New(profiles.Output(prof))
				if nArgs == 1 {
					output = file.New(args[0])
				}

				dotenvs = []profile.DotEnv{{Profile: prof, Path: output}}
			}

			for _, dotenv := range dotenvs {
				profile, err := profiles.Environment(dotenv.Profile)
				if err != nil {
					return err
				}

				if err := profile.ToDotEnv(dotenv.Path); err != nil {
					return err
				}

				fmt.Printf("Wrote profile %q to %q\n", dotenv.Profile, dotenv.Path)
			}

			return nil
		},
	}

	cmd.Flags().SortFlags = false

	cmd.Flags().BoolVarP(&all, "all", "a", false, "Write all profiles, ignoring the active profile")

	return cmd
}

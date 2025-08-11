package cli

import (
	"fmt"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"github.com/idelchi/envprof/internal/profile"
	"github.com/idelchi/godyl/pkg/path/file"
)

// Write defines the command for writing profile variables to one or multiple dotenv files.
//
//nolint:forbidigo	// Command prints out to the console.
func Write(envprof *[]string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "write [profile] [file]",
		Short: "Write profile variables to dotenv files",
		Long: heredoc.Doc(`
			Write variables to dotenv files.

			No arguments: write all profiles to <profile>.env files.
			<profile>: write <profile> to <profile>.env.
			<profile> <file>: write <profile> to <file>.
		`),
		Example: heredoc.Doc(`
			# Write 'dev' to dev.env
			envprof write dev

			# Write 'dev' to a specific file
			envprof write dev .env

			# Write all profiles to <profile>.env files
			envprof write
		`),
		Aliases: []string{"w"},
		Args:    cobra.RangeArgs(0, 2), //nolint:mnd	// The command takes between 0 and 2 arguments as documented.
		RunE: func(_ *cobra.Command, args []string) error {
			profiles, err := load(*envprof)
			if err != nil {
				return err
			}

			var dotenvs []profile.DotEnv
			switch len(args) {
			case 0:
				profiles := profiles.Names()
				dotenvs = make([]profile.DotEnv, 0, len(profiles))
				for _, prof := range profiles {
					dotenvs = append(dotenvs, profile.DotEnv{Profile: prof, Path: file.New(prof + ".env")})
				}
			case 1:
				prof := args[0]
				dotenvs = []profile.DotEnv{{Profile: prof, Path: file.New(prof + ".env")}}
			case 2: //nolint:mnd	// The number is obvious from the context.
				dotenvs = []profile.DotEnv{{Profile: args[0], Path: file.New(args[1])}}
			}

			for _, dotenv := range dotenvs {
				profile, err := profiles.Environment(dotenv.Profile)
				if err != nil {
					return err //nolint:wrapcheck // Error does not need additional wrapping.
				}

				if err := profile.ToDotEnv(dotenv.Path); err != nil {
					return err //nolint:wrapcheck // Error does not need additional wrapping.
				}

				fmt.Printf("Wrote profile %q to %q\n", dotenv.Profile, dotenv.Path)
			}

			return nil
		},
	}

	return cmd
}

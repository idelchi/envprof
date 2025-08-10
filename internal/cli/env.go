package cli

import (
	"fmt"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"github.com/idelchi/envprof/internal/profile"
	"github.com/idelchi/godyl/pkg/path/file"
)

// Env defines the command for exporting profile variables to one or multiple dotenv files.
//
//nolint:forbidigo	// Command prints out to the console.
func Env(envprof *[]string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "env <profile> [file]",
		Short: "Write profile variables to dotenv files",
		Long: heredoc.Doc(`
			Write profile variables to one or multiple dotenv files.

			If no profile is provided, all profiles will be written out to files named as <profile>.env.

			When specifying a specific profile, the default output file will be named as <profile>.env.

			This can be overridden by providing a file argument.
		`),
		Example: heredoc.Doc(`
			# Write out to 'dev.env'
			$ envprof env dev

			# Write out to a specific dotenv file
			$ envprof env dev .env

			# Write out all profiles to separate files named after <profile>
			$ envprof env
		`),

		Aliases: []string{"e", "write"},
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
				if err := profiles.ToDotEnv(dotenv.Profile, dotenv.Path); err != nil {
					return err //nolint:wrapcheck // Error does not need additional wrapping.
				}

				fmt.Printf("Wrote profile %q to %q\n", dotenv.Profile, dotenv.Path)
			}

			return nil
		},
	}

	return cmd
}

package cli

import (
	"fmt"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"github.com/idelchi/envprof/internal/dotenv"
	"github.com/idelchi/godyl/pkg/path/file"
)

// Env defines the command for exporting profile variables to one or multiple dotenv files.
//
//nolint:forbidigo	// Command print out to the console.
func Env(envprof *[]string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "env <profile> [file]",
		Short: "Write profile variables to one or multiple dotenv files",
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

		Aliases: []string{"e"},
		Args:    cobra.RangeArgs(0, 2), //nolint:mnd	// The command takes between 0 and 2 arguments as documented.
		RunE: func(_ *cobra.Command, args []string) error {
			profiles, err := load(*envprof)
			if err != nil {
				return err
			}

			var targets []dotenv.DotEnv
			switch len(args) {
			case 0:
				names := profiles.Names()
				targets = make([]dotenv.DotEnv, 0, len(names))
				for _, p := range names {
					targets = append(targets, dotenv.DotEnv{Profile: p, Path: file.New(p + ".env")})
				}
			case 1:
				p := args[0]
				targets = []dotenv.DotEnv{{Profile: p, Path: file.New(p + ".env")}}
			case 2: //nolint:mnd	// The number is obvious from the context.
				targets = []dotenv.DotEnv{{Profile: args[0], Path: file.New(args[1])}}
			}

			for _, d := range targets {
				if err := d.WriteFrom(&profiles); err != nil {
					return err //nolint:wrapcheck // Error does not need additional wrapping.
				}

				fmt.Printf("Wrote profile %q to %q\n", d.Profile, d.Path)
			}

			return nil
		},
	}

	return cmd
}

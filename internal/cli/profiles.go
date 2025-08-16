package cli

import (
	"fmt"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"
)

// Profiles returns the cobra command for listing profiles.
//
//nolint:forbidigo	// Command prints out to the console.
func Profiles(options *Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "profiles",
		Short: "List all available profiles",
		Long:  "List all profiles sorted",
		Example: heredoc.Doc(`
			# List all profiles
			$ envprof profiles
		`),
		Aliases: []string{"profs"},
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			profiles, err := loadProfiles(options.EnvProf)
			if err != nil {
				return err
			}

			name, _ := profileOrDefault(profiles, options.Profile)

			prefix := ""
			for _, profile := range profiles.Names() {
				if options.Verbose {
					prefix = "- "
					if profile == name {
						prefix = "* "
					}
				}
				fmt.Printf("%s%s\n", prefix, profile)
			}

			return nil
		},
	}

	cmd.Flags().SortFlags = false

	return cmd
}

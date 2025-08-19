package cli

import (
	"fmt"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"
)

// Profiles returns the cobra command for listing profiles.
func Profiles(options *Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "profiles",
		Short: "List all available profiles",
		Long:  "List all available profiles sorted",
		Example: heredoc.Doc(`
			# List all profiles
			$ envprof profiles

			# Highlight active profile
			$ envprof profiles -v
		`),
		Aliases: []string{"profs"},
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			profiles, _, err := loadProfiles(options.EnvProf)
			if err != nil {
				return err
			}

			// It's not important if the active profile is existing or not.
			name, _ := profiles.GetOrDefault(options.Profile)

			for _, profile := range profiles.Names() {
				//nolint:forbidigo	// Command prints out to the console.
				fmt.Println(formatProfile(profile, options.Verbose, profile == name))
			}

			return nil
		},
	}

	cmd.Flags().SortFlags = false

	return cmd
}

// formatProfile optionally marks out the active profile.
func formatProfile(profile string, decorate, isActive bool) string {
	if !decorate {
		return profile
	}

	if isActive {
		return "* " + profile
	}

	return "- " + profile
}

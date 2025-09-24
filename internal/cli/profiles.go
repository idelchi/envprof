package cli

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"github.com/idelchi/envprof/internal/profiles"
	"github.com/idelchi/godyl/pkg/pretty"
)

// Profiles returns the cobra command for listing profiles.
func Profiles(options *Options) *cobra.Command {
	var rendered bool

	cmd := &cobra.Command{
		Use:   "profiles",
		Short: "List all available profiles",
		Long:  "List all available profiles sorted",
		Example: heredoc.Doc(`
			# List all profiles
			envprof profiles

			# Highlight active profile
			envprof profiles -v
		`),
		Aliases: []string{"profs"},
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if rendered {
				envprof, err := LoadEnvProf(options)
				if err != nil && !errors.Is(err, profiles.ErrValidation) {
					return err
				}

				pretty.PrintYAML(envprof.Profiles())

				return nil
			}

			envprof, err := LoadEnvProf(options)
			if err != nil {
				return err
			}

			// It's not important if the active profile is existing or not.
			name, _ := envprof.GetOrDefault(options.Profile)

			profiles := envprof.Profiles().Names()

			if options.Verbose {
				slices.SortFunc(profiles, func(first, second string) int {
					if first == name {
						return -1
					}

					if second == name {
						return 1
					}

					return strings.Compare(first, second)
				})
			}

			for _, profile := range profiles {
				//nolint:forbidigo	// Command prints out to the console.
				fmt.Println(formatProfile(profile, options.Verbose, profile == name))
			}

			return nil
		},
	}

	cmd.Flags().SortFlags = false

	cmd.Flags().BoolVarP(&rendered, "rendered", "r", false, "Write the rendered profile to stdout")

	return cmd
}

// formatProfile formats a profile name with optional decoration to mark the active profile.
func formatProfile(profile string, decorate, isActive bool) string {
	if !decorate {
		return profile
	}

	if isActive {
		return "* " + profile
	}

	return "- " + profile
}

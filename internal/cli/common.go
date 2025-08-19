package cli

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/idelchi/envprof/internal/environment"
	"github.com/idelchi/envprof/internal/envprof"
	"github.com/idelchi/envprof/internal/profiles"
	"github.com/idelchi/godyl/pkg/path/file"
	"github.com/idelchi/godyl/pkg/path/files"
)

// loadProfiles loads the profile store from the specified file and fallbacks.
func loadProfiles(paths []string) (profiles.Profiles, file.File, error) {
	profiles, file, err := envprof.Load(files.New("", paths...))
	if err != nil {
		return nil, file, err
	}

	return profiles, file, nil
}

// loadProfile fully loads and resolves the profile variables from the specified file and fallbacks.
func loadProfile(options *Options) (environment.Environment, error) {
	profiles, _, err := loadProfiles(options.EnvProf)
	if err != nil {
		return environment.Environment{}, err
	}

	profile, err := profiles.GetOrDefault(options.Profile)
	if err != nil {
		return environment.Environment{}, err
	}

	return profiles.Environment(profile, options.Overlay...)
}

// loadPlan returns the plan for the specified profile.
func loadPlan(options *Options) (profiles.Steps, error) {
	profiles, _, err := loadProfiles(options.EnvProf)
	if err != nil {
		return nil, err
	}

	profile, err := profiles.GetOrDefault(options.Profile)
	if err != nil {
		return nil, err
	}

	return profiles.Plan(profile)
}

// UnknownSubcommandAction handles unknown cobra subcommands.
// Implements cobra.Command.RunE to provide helpful error messages
// and suggestions when an unknown subcommand is used. Required
// when TraverseChildren is true, as this disables cobra's built-in
// suggestion system. See:
// - https://github.com/spf13/cobra/issues/981
// - https://github.com/containerd/nerdctl/blob/242e6fc6e861b61b878bd7df8bf25e95674c036d/cmd/nerdctl/main.go#L401-L418
func UnknownSubcommandAction(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return cmd.Help() //nolint: wrapcheck	// The help message (error) should be returned as is
	}

	err := fmt.Sprintf("unknown subcommand %q for %q", args[0], cmd.Name())

	if suggestions := cmd.SuggestionsFor(args[0]); len(suggestions) > 0 {
		err += "\n\nDid you mean this?\n"

		for _, s := range suggestions {
			err += fmt.Sprintf("\t%v\n", s)
		}
	}

	return errors.New(err) //nolint: err113 	 // The error should be returned as is
}

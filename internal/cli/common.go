package cli

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/idelchi/envprof/internal/environment"
	"github.com/idelchi/envprof/internal/envprof"
	"github.com/idelchi/envprof/internal/profiles"
	"github.com/idelchi/envprof/internal/step"
	"github.com/idelchi/godyl/pkg/path/files"
)

// EnvProf returns the envprof instance, resolving only the file path.
func EnvProf(options *Options) (*envprof.EnvProf, error) {
	envprof, err := envprof.NewFrom(files.New("", options.EnvProf...))
	if err != nil {
		return nil, err
	}

	return envprof, nil
}

// LoadEnvProf returns the loaded envprof instance.
func LoadEnvProf(options *Options) (*envprof.EnvProf, error) {
	envprof, err := EnvProf(options)
	if err != nil {
		return nil, err
	}

	if err := envprof.Load(); err != nil {
		return nil, err
	}

	return envprof, nil
}

// LoadProfiles returns the existing profiles.
func LoadProfiles(options *Options) (profiles.Profiles, error) {
	envprof, err := LoadEnvProf(options)
	if err != nil {
		return nil, err
	}

	return envprof.Profiles(), nil
}

// LoadPlan returns the plan for the specified profile.
func LoadPlan(options *Options) (profiles.Profiles, string, step.Steps, error) {
	envprof, err := LoadEnvProf(options)
	if err != nil {
		return nil, "", nil, err
	}

	profile, err := envprof.GetOrDefault(options.Profile)
	if err != nil {
		return nil, "", nil, err
	}

	profiles := envprof.Profiles()

	steps, err := profiles.Plan(profile, options.Overlay...)
	if err != nil {
		return nil, "", nil, err
	}

	return profiles, profile, steps, nil
}

// LoadProfile returns the loaded and resolved profile.
func LoadProfile(options *Options) (environment.Environment, error) {
	profiles, profile, steps, err := LoadPlan(options)
	if err != nil {
		return environment.Environment{}, err
	}

	return profiles.Environment(profile, steps)
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

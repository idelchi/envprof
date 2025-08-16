package cli

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/idelchi/envprof/internal/profile"
	"github.com/idelchi/godyl/pkg/path/files"
)

// loadProfiles loads the profile store from the specified file and fallbacks.
func loadProfiles(paths []string) (profile.Profiles, error) {
	file, ok := files.New("", paths...).Exists()
	if !ok {
		return nil, fmt.Errorf("profile file not found: searched for %v", paths)
	}

	profiles, err := profile.New(file)
	if err != nil {
		return nil, err
	}

	store, err := profiles.Load()
	if err != nil {
		return nil, fmt.Errorf("loading profile from %s: %w", file.String(), err)
	}

	return store.Profiles, nil
}

// loadProfile loads the profile variables from the specified file and fallbacks.
func loadProfile(paths []string, name string) (*profile.InheritanceTracker, error) {
	profiles, err := loadProfiles(paths)
	if err != nil {
		return nil, err
	}

	name, err = profileOrDefault(profiles, name)
	if err != nil {
		return nil, err
	}

	vars, err := profiles.Environment(name)
	if err != nil {
		return nil, err
	}

	return vars, nil
}

// profileOrDefault returns the profile or the default profile if not found.
func profileOrDefault(profiles profile.Profiles, name string) (string, error) {
	if name == "" {
		name = profiles.Default()
		if name == "" {
			return "", errors.New("no default profile found and none specified with --profile")
		}
	}

	return name, nil
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

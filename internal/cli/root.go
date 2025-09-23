package cli

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"
)

// Options represents the root level configuration for the CLI application.
type Options struct {
	// EnvProf is the list of candidate profile files to load.
	EnvProf []string
	// Profile is the selected profile.
	Profile string
	// Verbose enables verbose output.
	Verbose bool
	// Overlay contains the profiles to overlay on top of the current profile.
	Overlay []string
}

// Execute runs the root command for the envprof CLI application.
func Execute(version string) error {
	options := &Options{
		EnvProf: []string{
			"envprof.yaml",
			"envprof.yml",
			"envprof.toml",
		},
	}

	home, err := os.UserHomeDir()
	if err == nil {
		base := filepath.ToSlash(filepath.Join(home, ".config", "envprof", "envprof"))

		options.EnvProf = append(options.EnvProf, base+".yaml", base+".yml", base+".toml")
	}

	root := &cobra.Command{
		Use:   "envprof",
		Short: "Manage env profiles in YAML/TOML with inheritance",
		Long: heredoc.Docf(`
			Manage environment profiles defined in YAML or TOML, with inheritance and dotenv imports.

			The config file is chosen via --file (can be a list) or ENVPROF_FILE.
			By default, envprof searches in the current directory for:

			%s

			Use subcommands to list profiles, export variables, write dotenv files,
			spawn a subshell, or exec a command with a selected profile.
		`, " - "+strings.Join(options.EnvProf, "\n - ")),
		Example: heredoc.Doc(`
			# List variables for 'dev'
			envprof --profile dev list -v

			# Create a dotenv file from a profile
			envprof --profile dev write .env

			# Export to current shell
			eval "$(envprof --profile dev export)"

			# Enter a subshell with a profile
			envprof --profile dev shell --shell zsh

			# Execute a command with a profile
			envprof --profile dev exec -- ls -la
		`),
		Version:          version,
		SilenceErrors:    true,
		TraverseChildren: true,
		SilenceUsage:     true,
		RunE:             UnknownSubcommandAction,
	}

	root.SetVersionTemplate("{{ .Version }}\n")
	root.SetHelpCommand(&cobra.Command{Hidden: true})

	root.Flags().SortFlags = false
	root.PersistentFlags().SortFlags = false

	root.CompletionOptions.DisableDefaultCmd = true
	cobra.EnableCommandSorting = false

	if file := os.Getenv("ENVPROF_FILE"); file != "" {
		options.EnvProf = []string{file}
	}

	root.Flags().
		StringSliceVarP(&options.EnvProf, "file", "f", options.EnvProf, "Config file to use, in order of preference")
	root.Flags().
		StringVarP(&options.Profile, "profile", "p", "", "Profile to activate")
	root.PersistentFlags().
		BoolVarP(&options.Verbose, "verbose", "v", false, "Increase verbosity level")
	root.Flags().
		StringSliceVarP(&options.Overlay, "overlay", "o", nil, "Profiles to overlay on top of the current profile")

	root.AddCommand(
		Path(options),
		Profiles(options),
		List(options),
		Export(options),
		Write(options),
		Shell(options),
		Exec(options),
		Diff(options),
	)

	if err := root.Execute(); err != nil {
		return err
	}

	return nil
}

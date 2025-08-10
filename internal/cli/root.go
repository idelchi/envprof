package cli

import (
	"os"
	"strings"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"
)

// Execute runs the root command for the envprof CLI application.
func Execute(version string) error {
	envprof := &[]string{
		"envprof.yaml",
		"envprof.yml",
		"envprof.toml",
	}

	root := &cobra.Command{
		Use:   "envprof",
		Short: "Manage env profiles in YAML/TOML with inheritance",
		Long: heredoc.Docf(`
			Manage env profiles in YAML/TOML with inheritance.

			Profiles are loaded from a config file, which can be specified with the --file flag,
			or by setting the ENVPROF_FILE environment variable.

			The tool will by default search for the following files in the current directory:

			%s

			Profiles can be listed, exported, and used to spawn a new shell with the profile's environment.

			Profiles can be inherited from other profiles and dotenv files.
		`, " - "+strings.Join(*envprof, "\n - ")),
		Example: heredoc.Doc(`
			# List the variables for the 'dev' profile
			$ envprof list dev -v

			# Create a dotenv file from a given profile
			$ envprof env dev

			# Eval the profile in the current shell
			$ eval "$(envprof export dev)"

			# Enter a new shell with the profile's environment
			$ envprof shell dev --shell zsh
		`),
		Version:       version,
		SilenceErrors: true,
		PersistentPreRun: func(cmd *cobra.Command, _ []string) {
			// Do not print usage after basic validation has been done.
			cmd.SilenceUsage = true
		},
	}

	root.SetVersionTemplate("{{ .Version }}\n")
	root.SetHelpCommand(&cobra.Command{Hidden: true})

	root.Flags().SortFlags = false
	root.CompletionOptions.DisableDefaultCmd = true
	cobra.EnableCommandSorting = false

	if file := os.Getenv("ENVPROF_FILE"); file != "" {
		envprof = &[]string{file}
	}

	root.PersistentFlags().
		StringSliceVarP(envprof, "file", "f", *envprof, "config file to use, in order of preference")

	root.AddCommand(
		List(envprof),
		Export(envprof),
		Env(envprof),
		Shell(envprof),
	)

	if err := root.Execute(); err != nil {
		return err //nolint:wrapcheck	// Error does not need additional wrapping.
	}

	return nil
}

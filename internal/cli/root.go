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
			Manage environment profiles defined in YAML or TOML, with inheritance and dotenv imports.

			The config file is chosen via --file (can be a list) or ENVPROF_FILE.
			By default, envprof searches in the current directory for:

			%s

			Use subcommands to list profiles, export variables, write dotenv files,
			spawn a subshell, or exec a command with a selected profile.
		`, " - "+strings.Join(*envprof, "\n - ")),
		Example: heredoc.Doc(`
			# List variables for 'dev'
			envprof list dev -v

			# Create a dotenv file from a profile
			envprof write dev .env

			# Export to current shell
			eval "$(envprof export dev)"

			# Enter a subshell with a profile
			envprof shell dev --shell zsh

			# Execute a command with a profile
			envprof exec dev -- ls -la
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
	root.PersistentFlags().SortFlags = false

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
		Write(envprof),
		Shell(envprof),
		Exec(envprof),
	)

	if err := root.Execute(); err != nil {
		return err //nolint:wrapcheck	// Error does not need additional wrapping.
	}

	return nil
}

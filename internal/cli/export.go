package cli

import (
	"fmt"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"github.com/idelchi/envprof/internal/environment"
)

// Export defines the command for exporting a profile's variables.
func Export(options *Options) *cobra.Command {
	prefix := "export "

	cmd := &cobra.Command{
		Use:   "export",
		Short: "Emit '<prefix>KEY=VAL' lines",
		Long: heredoc.Doc(`
			Print lines suitable for eval in the current shell:

			 <prefix>KEY1=VAL1
			 <prefix>KEY2=VAL2
			 ...

			Default prefix is "export ". Override with --prefix (e.g., "$env:" on PowerShell).
		`),
		Example: heredoc.Doc(`
			# Emit 'export KEY=VAL' lines for 'dev'
			envprof --profile dev export

			# Use a custom prefix (PowerShell)
			envprof --profile dev export --prefix "$env:"
		`),
		Aliases: []string{"x"},
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			env, err := LoadProfile(options)
			if err != nil {
				return err
			}

			formatter := environment.Formatter{
				WithKey: true,
				Prefix:  prefix,
			}

			envs := formatter.All(env)

			//nolint:forbidigo	// Command prints out to the console.
			fmt.Println(envs)

			return nil
		},
	}

	cmd.Flags().SortFlags = false

	cmd.Flags().StringVarP(&prefix, "prefix", "p", prefix, "Prefix for the export command")

	return cmd
}

package cli

import (
	"fmt"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"
)

// Export defines the command for exporting a profile's variables.
//
//nolint:forbidigo	// Command prints out to the console.
func Export(envprof *[]string) *cobra.Command {
	prefix := "export "

	cmd := &cobra.Command{
		Use:   "export <profile> [file]",
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
			envprof export dev

			# Use a custom prefix (PowerShell)
			envprof export dev --prefix "$env:"
		`),
		Aliases: []string{"x"},
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			prof := args[0]

			vars, err := loadProfileVars(*envprof, prof)
			if err != nil {
				return err
			}

			envs := vars.FormatAll(prefix, false)

			fmt.Println(envs)

			return nil
		},
	}

	cmd.Flags().SortFlags = false

	cmd.Flags().StringVarP(&prefix, "prefix", "p", prefix, "Prefix for the export command")

	return cmd
}

package cli

import (
	"fmt"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"
)

// Export defines the command for exporting a profile's variables.
// It emits '<prefix> KEY=VAL' lines or writes them out as a dotenv file.
// By default, 'prefix' is set to "export".
//
//nolint:forbidigo	// Command prints out to the console.
func Export(envprof *[]string) *cobra.Command {
	prefix := "export "

	cmd := &cobra.Command{
		Use:   "export <profile> [file]",
		Short: "Emit 'export KEY=VAL' lines to stdout",
		Long: heredoc.Doc(`
			Emit 'export KEY=VAL' lines to stdout as:

			<prefix>KEY1=VAL1
			<prefix>KEY2=VAL2
			<prefix>KEY3=VAL3

			The default prefix is "export " and can be customized with the --prefix flag.
		`),
		Example: heredoc.Doc(`
			# Emit 'export KEY=VAL' lines
			$ envprof export dev

			# Emit '$env:KEY=VAL' lines with a custom prefix
			$ envprof export dev --prefix "$env:"
		`),
		Aliases: []string{"x"},
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			profiles, err := load(*envprof)
			if err != nil {
				return err
			}

			prof := args[0]

			vars, err := profiles.Environment(prof)
			if err != nil {
				return err //nolint:wrapcheck	// Error does not need additional wrapping.
			}

			envs := vars.FormatAll(prefix, false)

			fmt.Println(envs)

			return nil
		},
	}

	cmd.Flags().StringVarP(&prefix, "prefix", "p", prefix, "Prefix for the export command")

	return cmd
}

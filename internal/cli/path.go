package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Path returns the cobra command for displaying the path to the configuration file used.
func Path(options *Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "path",
		Short: "Display the path of the configuration file used",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			envprof, err := EnvProf(options)
			if err != nil {
				return err
			}

			//nolint:forbidigo	// Command prints out to the console.
			fmt.Println(envprof.File().Path())

			return nil
		},
	}

	cmd.Flags().SortFlags = false

	return cmd
}

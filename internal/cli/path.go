package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/idelchi/envprof/internal/envprof"
	"github.com/idelchi/godyl/pkg/path/files"
)

// Path returns the cobra command for displaying the path to the configuration file used.
func Path(options *Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "path",
		Short: "Display the path of the configuration file used",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			file, _, err := envprof.File(files.New("", options.EnvProf...))
			if err != nil {
				return err
			}

			fmt.Println(file)

			return nil
		},
	}

	cmd.Flags().SortFlags = false

	return cmd
}

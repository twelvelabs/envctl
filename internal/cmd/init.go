package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/twelvelabs/termite/fsutil"

	"github.com/twelvelabs/envctl/internal/core"
)

func NewInitCmd(app *core.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a new envctl config",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			path := core.ConfigPathDefault

			if fsutil.PathExists(path) {
				app.UI.Out("Existing config found at: %s\n", path)
				ok, err := app.UI.Confirm("Overwrite", false)
				if err != nil {
					return err
				}
				if !ok {
					app.UI.Out("Exiting.\n")
					return nil
				}
			}

			app.UI.Out("Initializing config: %s\n", path)
			return os.WriteFile(path, core.ConfigContentDefault, 0644) //nolint: gosec
		},
	}
	return cmd
}

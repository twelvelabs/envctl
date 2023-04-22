package cmd

import (
	"time"

	"github.com/spf13/cobra"

	"github.com/twelvelabs/envctl/internal/core"
)

func NewVersionCmd(app *core.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Show full version info",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			app.UI.Out("Version: %s\n", app.Meta.Version)
			app.UI.Out("GOOS: %s\n", app.Meta.GOOS)
			app.UI.Out("GOARCH: %s\n", app.Meta.GOARCH)
			app.UI.Out("\n")
			app.UI.Out("Build Time: %s\n", app.Meta.BuildTime.Format(time.RFC3339))
			app.UI.Out("Build Commit: %s\n", app.Meta.BuildCommit)
			app.UI.Out("Build Version: %s\n", app.Meta.BuildVersion)
			app.UI.Out("Build Checksum: %s\n", app.Meta.BuildChecksum)
			app.UI.Out("Build Go Version: %s\n", app.Meta.BuildGoVersion)
			return nil
		},
	}

	return cmd
}

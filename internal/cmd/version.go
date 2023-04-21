package cmd

import (
	"fmt"
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
			fmt.Fprintln(app.IO.Out, "Version:", app.Meta.Version)
			fmt.Fprintln(app.IO.Out, "GOOS:", app.Meta.GOOS)
			fmt.Fprintln(app.IO.Out, "GOARCH:", app.Meta.GOARCH)
			fmt.Fprintln(app.IO.Out, "")
			fmt.Fprintln(app.IO.Out, "Build Time:", app.Meta.BuildTime.Format(time.RFC3339))
			fmt.Fprintln(app.IO.Out, "Build Commit:", app.Meta.BuildCommit)
			fmt.Fprintln(app.IO.Out, "Build Version:", app.Meta.BuildVersion)
			fmt.Fprintln(app.IO.Out, "Build Checksum:", app.Meta.BuildChecksum)
			fmt.Fprintln(app.IO.Out, "Build Go Version:", app.Meta.BuildGoVersion)
			return nil
		},
	}

	return cmd
}

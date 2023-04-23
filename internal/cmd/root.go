package cmd

import (
	"time"

	"github.com/spf13/cobra"

	"github.com/twelvelabs/envctl/internal/core"
)

func NewRootCmd(app *core.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "envctl",
		Short:   "Manage project environment variables",
		Version: app.Meta.Version,
		Args:    cobra.NoArgs,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			app.Logger.Debug("App initialized",
				"config", app.Config.ConfigPath,
				"duration", time.Since(app.CreatedAt),
			)
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			app.UI.Out("Hello 👋 \n")
			return nil
		},
		SilenceUsage: true,
	}

	flags := cmd.PersistentFlags()
	flags.StringVarP(
		&app.Config.ConfigPath,
		core.ConfigPathLongFlag,
		core.ConfigPathShortFlag,
		app.Config.ConfigPath,
		"Config path",
	)

	// Hide the built in `completion` subcommand
	cmd.CompletionOptions.HiddenDefaultCmd = true

	cmd.AddCommand(NewManCmd(app))
	cmd.AddCommand(NewVersionCmd(app))

	return cmd
}

package cmd

import (
	"time"

	"github.com/spf13/cobra"

	"github.com/twelvelabs/envctl/internal/core"
)

func NewRootCmd(app *core.App) *cobra.Command {
	noPrompt := false
	verbosity := 0

	cmd := &cobra.Command{
		Use:     "envctl",
		Short:   "Manage project environment variables",
		Version: app.Meta.Version,
		Args:    cobra.NoArgs,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			app.SetVerbosity(verbosity)
			if noPrompt {
				app.IO.SetInteractive(false)
			}

			app.Logger.Debug("App initialized",
				"config", app.Config.ConfigPath,
				"duration", time.Since(app.CreatedAt),
			)
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			name := app.UI.Cyan("envctl")
			app.UI.Out("Hello from %s 👋 \n", name)
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
		"config path",
	)

	flags.CountVarP(&verbosity, "verbose", "v", "enable verbose logging (increase via -vv)")
	flags.BoolVar(&noPrompt, "no-prompt", noPrompt, "do not prompt for input")

	// Hide the built in `completion` subcommand
	cmd.CompletionOptions.HiddenDefaultCmd = true

	cmd.AddCommand(NewManCmd(app))
	cmd.AddCommand(NewVersionCmd(app))

	return cmd
}

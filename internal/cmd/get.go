package cmd

import (
	"errors"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/twelvelabs/envctl/internal/core"
)

func NewGetCmd(app *core.App) *cobra.Command {
	allEnvNames := app.Config.EnvironmentNames()

	cmd := &cobra.Command{
		Use:       "get ENVIRONMENT KEY",
		Short:     "Get an environment variable value",
		Args:      cobra.ArbitraryArgs,
		ValidArgs: allEnvNames,
		RunE: func(cmd *cobra.Command, args []string) error {
			switch len(args) {
			case 0:
				return errors.New("missing environment and key names")
			case 1:
				return errors.New("missing key name")
			default:
				return runGet(app, args[0], args[1])
			}
		},
	}

	return cmd
}

func runGet(app *core.App, envName string, keyName string) error {
	start := time.Now()
	app.Logger.Debug(
		"App running",
		"cmd", "get",
		"env", envName,
		"key", keyName,
	)

	// Get the environment.
	env, err := app.Environments.Get(envName)
	if err != nil {
		return err
	}

	// Get the unresolved value or URL for the key.
	unresolved, exists := env.Vars.Get(keyName)
	if !exists {
		return fmt.Errorf("unknown variable key: %s", keyName)
	}

	// Resolve it through the appropriate store.
	resolved, err := app.Stores.Get(app.Context(), unresolved)
	if err != nil {
		return err
	}

	_, _ = fmt.Fprint(app.IO.Out, resolved+"\n")

	app.Logger.Debug(
		"App finish",
		"cmd", "get",
		"duration", time.Since(start),
	)
	return err
}

package cmd

import (
	"errors"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/twelvelabs/envctl/internal/core"
)

func NewSetCmd(app *core.App) *cobra.Command {
	allEnvNames := app.Config.EnvironmentNames()

	cmd := &cobra.Command{
		Use:       "set ENVIRONMENT KEY VALUE",
		Short:     "Sets an environment variable value",
		Args:      cobra.ArbitraryArgs,
		ValidArgs: allEnvNames,
		RunE: func(cmd *cobra.Command, args []string) error {
			switch len(args) {
			case 0:
				return errors.New("missing arguments: environment, key, and value")
			case 1:
				return errors.New("missing arguments: key and value")
			case 2:
				return errors.New("missing arguments: value")
			default:
				return runSet(app, args[0], args[1], args[2])
			}
		},
	}

	return cmd
}

func runSet(app *core.App, envName, keyName, updated string) error {
	start := time.Now()
	app.Logger.Debug(
		"App running",
		"cmd", "set",
		"env", envName,
		"key", keyName,
		"value", updated,
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

	// FIXME: don't allow setting via CLI arg - only via file or prompt.

	// Set it to the updated value.
	err = app.Stores.Set(app.Context(), unresolved, updated)
	if err != nil {
		return err
	}

	app.Logger.Debug(
		"App finish",
		"cmd", "get",
		"duration", time.Since(start),
	)
	return err
}

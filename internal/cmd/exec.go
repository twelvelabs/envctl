package cmd

import (
	"errors"

	"github.com/spf13/cobra"

	"github.com/twelvelabs/envctl/internal/core"
)

func NewExecCmd(app *core.App) *cobra.Command {
	allEnvNames := app.Config.EnvironmentNames()

	cmd := &cobra.Command{
		Use:       "exec ENVIRONMENT -- COMMAND",
		Short:     "Exec a command in a given environment",
		Args:      cobra.ArbitraryArgs,
		ValidArgs: allEnvNames,
		RunE: func(cmd *cobra.Command, args []string) error {
			switch len(args) {
			case 0:
				return errors.New("missing env name")
			case 1:
				return errors.New("missing command")
			default:
				name := args[0]
				args := args[1:]

				app.Logger.Debug(
					"exec start",
					"name", name,
					"args", args,
				)
				return execEnv(app, name, args)
			}
		},
	}

	cmd.Flags().BoolVar(&app.Config.DotEnv, "dotenv", app.Config.DotEnv, "create a temporary dotenv file")

	return cmd
}

func execEnv(app *core.App, name string, args []string) error {
	envSvc := core.NewEnvironmentService(app.Config)
	env, err := envSvc.Get(name)
	if err != nil {
		return err
	}

	resSvc := core.NewResolverService(app.Config, core.Resolvers)
	vars, err := resSvc.ResolveVars(env.Vars)
	if err != nil {
		return err
	}

	if app.Config.DotEnv {
		dotEnvSvc := core.NewDotEnvService("")

		var cleanup core.CleanupFunc
		vars, args, cleanup, err = dotEnvSvc.Create(vars, args)
		if err != nil {
			return err
		}
		defer cleanup()
	}

	execSvc := core.NewExecService(app.Config, app.ExecClient)
	_, err = execSvc.Run(app.Context(), args, vars)
	return err
}

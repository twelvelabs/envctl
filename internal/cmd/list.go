package cmd

import (
	"fmt"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/twelvelabs/envctl/internal/core"
)

func NewListCmd(app *core.App) *cobra.Command {
	allEnvNames := app.Config.EnvironmentNames()
	cmd := &cobra.Command{
		Use:       "list [NAME]",
		Short:     "List environment vars",
		Args:      cobra.ArbitraryArgs,
		ValidArgs: allEnvNames,
		RunE: func(cmd *cobra.Command, args []string) error {
			envNames := args
			if len(envNames) == 0 {
				envNames = allEnvNames
			}

			for _, name := range envNames {
				if err := listEnv(app, name); err != nil {
					return err
				}
			}

			return nil
		},
	}
	return cmd
}

func listEnv(app *core.App, envName string) error {
	env, err := app.Environments.Get(envName)
	if err != nil {
		return err
	}

	vars, err := app.Stores.MultiGet(app.Context(), env.Vars)
	if err != nil {
		return err
	}

	app.UI.Out("# %s \n", envName)
	app.UI.Out("---------------------------------------- \n")
	w := tabwriter.NewWriter(app.IO.Out, 0, 0, 1, ' ', 0)
	for k, v := range vars {
		_, _ = fmt.Fprintf(w, "%s\t%s \n", k, v)
	}
	_ = w.Flush()
	app.UI.Out("---------------------------------------- \n")
	app.UI.Out("\n")

	return nil
}

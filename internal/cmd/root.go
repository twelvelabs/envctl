package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/twelvelabs/envctl/internal/envctl"
)

func NewRootCmd(app *envctl.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "envctl",
		Short:   "Manage project environment variables",
		Version: app.Meta.Version,
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("Hello 👋 \n")
			return nil
		},
	}

	cmd.AddCommand(NewVersionCmd(app))

	return cmd
}

package main

import (
	"fmt"
	"os"

	"github.com/twelvelabs/envctl/internal/cmd"
	"github.com/twelvelabs/envctl/internal/envctl"
)

var (
	version = "dev"
	commit  = ""
	date    = ""
)

// The actual `main` logic.
// Broken out so we can safely use defer (see [os.Exit] docs).
func run() error {
	app, err := envctl.NewApp(version, commit, date)
	if err != nil {
		return err
	}
	defer app.Close()

	return cmd.NewRootCmd(app).ExecuteContext(app.Context())
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

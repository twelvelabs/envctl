package core

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/twelvelabs/termite/run"
)

func TestExecService(t *testing.T) {
	app := NewTestApp()
	// defer app.ExecClient.VerifyStubs(t)

	app.ExecClient.RegisterStub(
		run.MatchRegexp(`echo`),
		run.StringResponse(""),
	)

	svc := NewExecService(app.Config, app.ExecClient)
	cmd, err := svc.Run(t.Context(), []string{"echo"}, EnvVars{"FOO": "bar"})
	require.NoError(t, err)
	require.Equal(t, []string{"FOO=bar"}, cmd.Env)
}

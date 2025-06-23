package exec

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/twelvelabs/termite/run"

	"github.com/twelvelabs/envctl/internal/models"
)

func TestExecService(t *testing.T) {
	client := run.NewClient().WithStubbing()
	// defer client.VerifyStubs(t)

	client.RegisterStub(
		run.MatchRegexp(`echo`),
		run.StringResponse(""),
	)

	svc := NewExecService(client)
	cmd, err := svc.Run(t.Context(), []string{"echo"}, models.Vars{"FOO": "bar"})
	require.NoError(t, err)
	require.Equal(t, []string{"FOO=bar"}, cmd.Env)
}

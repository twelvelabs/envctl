package exec

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/twelvelabs/termite/run"

	"github.com/twelvelabs/envctl/internal/models"
)

func TestExecService(t *testing.T) {
	client := run.NewClient()

	t.Setenv("FROM_PARENT", "parent")
	t.Setenv("FOR_CHILD", "parent")

	svc := NewExecService(client)
	cmd, err := svc.Run(t.Context(), []string{"echo"}, models.Vars{
		"FOR_CHILD": "child",
	})
	require.NoError(t, err)

	// Subprocess should have received both vars.
	require.Contains(t, cmd.Env, "FROM_PARENT=parent")
	require.Contains(t, cmd.Env, "FOR_CHILD=child")
}

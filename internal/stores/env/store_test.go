package env

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnvStore(t *testing.T) {
	var (
		val string
		err error
	)
	ctx := t.Context()
	store := NewEnvStore()
	defer store.Close()

	val, err = store.Get(ctx, "env://")
	require.ErrorIs(t, err, ErrEmptyKeyName)
	require.Equal(t, "", val)

	val, err = store.Get(ctx, "env://ENV_VAR_THAT_IS_NOT_SET")
	require.ErrorIs(t, err, ErrMissingKey)
	require.Equal(t, "", val)

	val, err = store.Get(ctx, "env://ENV_VAR_THAT_IS_NOT_SET?default=nope")
	require.NoError(t, err)
	require.Equal(t, "nope", val)

	t.Setenv("ENV_VAR_SET_TO_EMPTY_STRING", "")
	val, err = store.Get(ctx, "env://ENV_VAR_SET_TO_EMPTY_STRING?default=nope")
	require.NoError(t, err)
	require.Equal(t, "", val)

	t.Setenv("ENV_VAR_SET_TO_NON_EMPTY_STRING", "hello")
	val, err = store.Get(ctx, "env://ENV_VAR_SET_TO_NON_EMPTY_STRING")
	require.NoError(t, err)
	require.Equal(t, "hello", val)

	// These should be no-ops.
	err = store.Set(ctx, "", "")
	require.NoError(t, err)
	err = store.Delete(ctx, "")
	require.NoError(t, err)
}

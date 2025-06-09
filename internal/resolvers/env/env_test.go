package env

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnvResolver(t *testing.T) {
	var (
		u   *url.URL
		val string
		err error
	)
	resolver := NewEnvResolver()

	u, _ = url.Parse("env://")
	val, err = resolver.Resolve(u)
	require.ErrorIs(t, err, ErrEmptyKeyName)
	require.Equal(t, "", val)

	u, _ = url.Parse("env://ENV_VAR_THAT_IS_NOT_SET")
	val, err = resolver.Resolve(u)
	require.ErrorIs(t, err, ErrMissingKey)
	require.Equal(t, "", val)

	u, _ = url.Parse("env://ENV_VAR_THAT_IS_NOT_SET?default=nope")
	val, err = resolver.Resolve(u)
	require.NoError(t, err)
	require.Equal(t, "nope", val)

	t.Setenv("ENV_VAR_SET_TO_EMPTY_STRING", "")
	u, _ = url.Parse("env://ENV_VAR_SET_TO_EMPTY_STRING?default=nope")
	val, err = resolver.Resolve(u)
	require.NoError(t, err)
	require.Equal(t, "", val)

	t.Setenv("ENV_VAR_SET_TO_NON_EMPTY_STRING", "hello")
	u, _ = url.Parse("env://ENV_VAR_SET_TO_NON_EMPTY_STRING")
	val, err = resolver.Resolve(u)
	require.NoError(t, err)
	require.Equal(t, "hello", val)
}

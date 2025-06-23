package models

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValue_URL(t *testing.T) {
	require.Equal(t, "foo", Value("foo").URL().String())
	require.Equal(t, "foo/bar", Value("foo/bar").URL().String())
	require.Equal(t, "/foo/bar", Value("/foo/bar").URL().String())

	u := Value("env://FOO?default=bar").URL()
	require.Equal(t, "env", u.Scheme)
	require.Equal(t, "FOO", u.Host)
	require.Equal(t, "bar", u.Query().Get("default"))
}

func TestVars_Environ(t *testing.T) {
	vars := Vars{
		"AAA": "something",
		"BBB": "something with spaces",
		"CCC": "something \"quoted\"",
		"DDD": "something\nmultiline",
	}
	require.Equal(t, []string{
		"AAA=something",
		"BBB=something with spaces",
		"CCC=something \"quoted\"",
		"DDD=something\nmultiline",
	}, vars.Environ())
}

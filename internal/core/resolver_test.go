package core

import (
	"errors"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type mockResolver struct {
	Resolved string
	Err      error
}

func (r *mockResolver) Resolve(u *url.URL) (string, error) {
	return r.Resolved, r.Err
}

func TestResolverService(t *testing.T) {
	config, err := NewTestConfig()
	require.NoError(t, err)
	svc := NewResolverService(config, map[string]Resolver{
		"foo+one": &mockResolver{
			Resolved: "Foo One",
		},
		"foo+two": &mockResolver{
			Resolved: "Foo Two",
		},
		"err": &mockResolver{
			Err: errors.New("boom"),
		},
	})

	vars, err := svc.ResolveVars(EnvVars{
		"ONE":   "foo+one://something/one",
		"TWO":   "foo+two://something/two",
		"THREE": "not a URL",
		"FOUR":  strings.Join([]string{"one", "two", "three"}, "\n"),
		"FIVE":  "https://something/",
	})
	require.NoError(t, err)
	require.Equal(t, EnvVars{
		"ONE":   "Foo One",
		"TWO":   "Foo Two",
		"THREE": "not a URL",
		"FOUR":  strings.Join([]string{"one", "two", "three"}, "\n"),
		"FIVE":  "https://something/",
	}, vars)

	_, err = svc.ResolveVars(EnvVars{
		"ONE": "err://something/that/explodes",
	})
	require.ErrorContains(t, err, "boom")
}

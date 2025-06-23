package google

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGSMStore_names(t *testing.T) {
	store := NewGSMStore(nil)

	names, err := store.names("secret+google:///projects/my-project-id/secrets/my-secret-id/versions/latest")
	require.NoError(t, err)
	require.Equal(t, "projects/my-project-id", names.Parent)
	require.Equal(t, "projects/my-project-id/secrets/my-secret-id", names.Secret)
	require.Equal(t, "projects/my-project-id/secrets/my-secret-id/versions/latest", names.Version)

	names, err = store.names("secret+google:///projects/my-project-id/locations/global/secrets/my-secret-id/versions/latest") //nolint:lll
	require.NoError(t, err)
	require.Equal(t, "projects/my-project-id/locations/global", names.Parent)
	require.Equal(t, "projects/my-project-id/locations/global/secrets/my-secret-id", names.Secret)
	require.Equal(t, "projects/my-project-id/locations/global/secrets/my-secret-id/versions/latest", names.Version)

	names, err = store.names("secret+google:///lol/wat")
	require.ErrorContains(t, err, "invalid URL")
	require.Nil(t, names)

	names, err = store.names("secret+google:///projects//secrets//versions/latest")
	require.ErrorContains(t, err, "invalid URL")
	require.Nil(t, names)

	names, err = store.names("secret+google:///projects//locations//secrets//versions/latest")
	require.ErrorContains(t, err, "invalid URL")
	require.Nil(t, names)
}

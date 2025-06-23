package stores

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/twelvelabs/envctl/internal/models"
)

func TestStoreService_Close(t *testing.T) {
	ctx := t.Context()
	svc := NewStoreService(ctx, TestStoreFactories)

	// Should work when no stores have been initialized.
	require.NoError(t, svc.Close())

	// Now initialize the default store.
	_, _ = svc.Get(ctx, "foo")
	require.NoError(t, svc.Close())

	// The failure store should return an error on close.
	_, _ = svc.Get(ctx, "err+store://boom")
	require.ErrorContains(t, svc.Close(), "boom")
}

func TestStoreService_MultiGet(t *testing.T) {
	ctx := t.Context()
	svc := NewStoreService(ctx, TestStoreFactories)

	vars, err := svc.MultiGet(ctx, models.Vars{
		"AAA": "something",
		"BBB": "something with spaces",
		"CCC": "something \"quoted\"",
		"DDD": "something\nmultiline",
		"EEE": "https://some/normal/url",
	})
	require.NoError(t, err)
	require.Equal(t, models.Vars{
		"AAA": "something",
		"BBB": "something with spaces",
		"CCC": "something \"quoted\"",
		"DDD": "something\nmultiline",
		"EEE": "https://some/normal/url",
	}, vars)

	vars, err = svc.MultiGet(ctx, models.Vars{
		"AAA": "err+store://boom",
	})
	require.ErrorContains(t, err, "boom")
	require.Nil(t, vars)
}

func TestStoreService_ErrorInFactory(t *testing.T) {
	ctx := t.Context()
	svc := NewStoreService(ctx, TestStoreFactories)

	value, err := svc.Get(ctx, "err+factory://...")
	require.ErrorContains(t, err, "boom")
	require.Equal(t, "", value)
}

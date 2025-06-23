package fallback

import (
	"context"

	"github.com/twelvelabs/envctl/internal/models"
)

func FallbackStoreFactory(_ context.Context) (models.Store, error) { //nolint:ireturn
	store := NewFallbackStore()
	return store, nil
}

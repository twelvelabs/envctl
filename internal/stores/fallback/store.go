package fallback

import (
	"context"

	"github.com/twelvelabs/envctl/internal/models"
)

// FallbackStore is a noop store that simply returns values as-is.
type FallbackStore struct {
}

// NewFallbackStore returns a new [FallbackStore].
func NewFallbackStore() *FallbackStore {
	return &FallbackStore{}
}

func (s *FallbackStore) Close() error {
	return nil
}

func (s *FallbackStore) Get(ctx context.Context, value models.Value) (string, error) {
	return value.String(), nil
}

func (s *FallbackStore) Set(ctx context.Context, value models.Value, updated string) error {
	return nil
}

func (s *FallbackStore) Delete(ctx context.Context, value models.Value) error {
	return nil
}

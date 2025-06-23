package test

import (
	"context"

	"github.com/twelvelabs/envctl/internal/models"
)

// FailureStore is a store that always returns an error.
// Should only be used in tests.
type FailureStore struct {
	err error
}

// NewFailureStore returns a new [FailureStore].
func NewFailureStore(err error) *FailureStore {
	return &FailureStore{
		err: err,
	}
}

func (s *FailureStore) Close() error {
	return s.err
}

func (s *FailureStore) Get(ctx context.Context, value models.Value) (string, error) {
	return "", s.err
}

func (s *FailureStore) Set(ctx context.Context, value models.Value, updated string) error {
	return s.err
}

func (s *FailureStore) Delete(ctx context.Context, value models.Value) error {
	return s.err
}

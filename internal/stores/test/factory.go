package test

import (
	"context"
	"errors"

	"github.com/twelvelabs/envctl/internal/models"
)

func FailureFactory(ctx context.Context) (models.Store, error) { //nolint:ireturn
	return nil, errors.New("boom")
}

func FailureStoreFactory(ctx context.Context) (models.Store, error) { //nolint:ireturn
	store := NewFailureStore(errors.New("boom"))
	return store, nil
}

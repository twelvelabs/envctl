package env

import (
	"context"

	"github.com/twelvelabs/envctl/internal/models"
)

func EnvStoreFactory(_ context.Context) (models.Store, error) { //nolint:ireturn
	store := NewEnvStore()
	return store, nil
}

package env

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/twelvelabs/envctl/internal/models"
)

var (
	ErrEmptyKeyName = errors.New("empty key name")
	ErrMissingKey   = errors.New("missing key")
)

// EnvStore is a noop store that simply returns values as-is.
type EnvStore struct {
}

// NewEnvStore returns a new [EnvStore].
func NewEnvStore() *EnvStore {
	return &EnvStore{}
}

func (s *EnvStore) Close() error {
	return nil
}

func (s *EnvStore) Get(ctx context.Context, value models.Value) (string, error) {
	u := value.URL()
	// Resolve the key.
	key := strings.TrimSpace(u.Hostname())
	if key == "" {
		return "", fmt.Errorf("%w: %s", ErrEmptyKeyName, u.String())
	}

	// Resolve the value.
	val, ok := os.LookupEnv(key)
	if !ok {
		// No env var set. Check for a default value.
		params := u.Query()
		if !params.Has("default") {
			return "", fmt.Errorf("%w: %s", ErrMissingKey, u.String())
		}
		val = params.Get("default")
	}

	return val, nil
}

func (s *EnvStore) Set(ctx context.Context, value models.Value, updated string) error {
	return nil
}

func (s *EnvStore) Delete(ctx context.Context, value models.Value) error {
	return nil
}

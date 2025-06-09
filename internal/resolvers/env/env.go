package env

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"
)

var (
	ErrEmptyKeyName = errors.New("empty key name")
	ErrMissingKey   = errors.New("missing key")
)

func NewEnvResolver() *EnvResolver {
	return &EnvResolver{}
}

type EnvResolver struct {
}

func (receiver *EnvResolver) Resolve(u *url.URL) (string, error) {
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
		if params.Has("default") {
			return params.Get("default"), nil
		}
		return "", fmt.Errorf("%w: %s", ErrMissingKey, u.String())
	}

	return val, nil
}

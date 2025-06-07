package core

import (
	"context"
	"fmt"
	"os"

	getter "github.com/hashicorp/go-getter"
)

// Getter is a function that copies a package from `src` to `dst`.
type Getter func(ctx context.Context, src, dst string) error

// DefaultGetter uses hashicorp/go-getter to copy packages.
func DefaultGetter(ctx context.Context, src, dst string) error {
	pwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("unable to resolve working directory: %w", err)
	}

	client := getter.Client{
		Ctx:  ctx,
		Src:  src,
		Dst:  dst,
		Pwd:  pwd,
		Mode: getter.ClientModeFile,
	}

	if err := client.Get(); err != nil {
		return fmt.Errorf("unable to get: %w", err)
	}

	return nil
}

// MockGetter delegates to the supplied handler function.
type MockGetter struct {
	Called  bool
	Ctx     context.Context //nolint:containedctx
	Src     string
	Dst     string
	handler func() error
}

// NewMockGetter returns a new MockGetter struct.
func NewMockGetter(handler func() error) *MockGetter {
	return &MockGetter{handler: handler}
}

// Get implements the Getter type and is what should be passed to the store.
// It logs the arguments and delegates to the handler.
func (m *MockGetter) Get(ctx context.Context, src, dst string) error {
	m.Called = true
	m.Ctx = ctx
	m.Src = src
	m.Dst = dst
	return m.handler()
}

package models

import (
	"context"
	"io"
)

type Store interface {
	io.Closer

	// Get returns the value identified by the given URL.
	Get(ctx context.Context, value Value) (string, error)

	// Set updates the value identified by the given URL.
	Set(ctx context.Context, value Value, updated string) error

	// Delete destroys the value identified by the given URL.
	Delete(ctx context.Context, value Value) error
}

type StoreFactory func(ctx context.Context) (Store, error)

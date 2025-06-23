package stores

import (
	"context"
	"maps"

	"github.com/twelvelabs/envctl/internal/models"
	"github.com/twelvelabs/envctl/internal/stores/env"
	"github.com/twelvelabs/envctl/internal/stores/fallback"
	"github.com/twelvelabs/envctl/internal/stores/google"
	"github.com/twelvelabs/envctl/internal/stores/test"
)

var (
	DefaultProtocol       = ""
	DefaultStoreFactories = map[string]models.StoreFactory{
		"env":           env.EnvStoreFactory,
		"secret+google": google.GSMStoreFactory,
		DefaultProtocol: fallback.FallbackStoreFactory,
	}
	TestStoreFactories = map[string]models.StoreFactory{
		"err+factory":   test.FailureFactory,
		"err+store":     test.FailureStoreFactory,
		DefaultProtocol: fallback.FallbackStoreFactory,
	}

	// Ensure StoreService implements the Store interface.
	_ models.Store = &StoreService{}
)

// StoreService aggregates all stores into a single [models.Store].
// Individual stores are lazily constructed when needed using the
// [models.StoreFactory] registered for their URL protocol.
// Non-URL values will use the default/fallback store.
type StoreService struct {
	ctx       context.Context
	factories map[string]models.StoreFactory
	stores    map[string]models.Store
}

func NewStoreService(ctx context.Context, factories map[string]models.StoreFactory) *StoreService {
	stores := map[string]models.Store{}
	for k := range factories {
		stores[k] = nil
	}
	return &StoreService{
		ctx:       ctx,
		factories: factories,
		stores:    stores,
	}
}

// Close frees up any store resources that may have been opened.
func (s *StoreService) Close() error {
	for _, store := range s.stores {
		if store == nil {
			continue // store was never initialized
		}
		if err := store.Close(); err != nil {
			return err
		}
	}
	return nil
}

// MultiGet gets all vars from their respective stores and
// returns a new set of vars with the resolved values.
func (s *StoreService) MultiGet(ctx context.Context, vars models.Vars) (models.Vars, error) {
	batch := maps.Clone(vars)
	for key, val := range vars {
		val, err := s.Get(ctx, val)
		if err != nil {
			return nil, err
		}
		batch[key] = models.Value(val)
	}
	return batch, nil
}

// Get returns the value identified by the given URL.
func (s *StoreService) Get(ctx context.Context, value models.Value) (string, error) {
	store, err := s.storeFor(value)
	if err != nil {
		return "", err
	}
	return store.Get(ctx, value)
}

// Set updates the value identified by the given URL.
func (s *StoreService) Set(ctx context.Context, value models.Value, updated string) error {
	store, err := s.storeFor(value)
	if err != nil {
		return err
	}
	return store.Set(ctx, value, updated)
}

// Delete destroys the value identified by the given URL.
func (s *StoreService) Delete(ctx context.Context, value models.Value) error {
	store, err := s.storeFor(value)
	if err != nil {
		return err
	}
	return store.Delete(ctx, value)
}

func (s *StoreService) storeFor(value models.Value) (models.Store, error) { //nolint:ireturn
	var err error
	protocol := value.URL().Scheme

	// Fallback to the default if this isn't a known protocol.
	if _, ok := s.stores[protocol]; !ok {
		protocol = DefaultProtocol
	}

	// Stores are lazily initialized.
	// Doing this because some store factories require
	// valid cloud credentials, and we don't want every
	// CLI command to inherit that requirement.
	store := s.stores[protocol]
	if store == nil {
		factory := s.factories[protocol]
		store, err = factory(s.ctx)
		if err != nil {
			return nil, err
		}
		s.stores[protocol] = store
	}

	return store, nil
}

package google

import (
	"context"

	"github.com/twelvelabs/envctl/internal/models"
)

func GSMStoreFactory(ctx context.Context) (models.Store, error) { //nolint:ireturn
	client, err := NewGSMClient(ctx)
	if err != nil {
		return nil, err
	}
	store := NewGSMStore(client)
	return store, nil
}

package store

import (
	"context"
	"time"
)

func (s *StoreDatabase) SetWithContext(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return s.Store.Set(ctx, key, value, expiration).Err()
}

func (s *StoreDatabase) GetWithContext(ctx context.Context, key string) (string, error) {
	return s.Store.Get(ctx, key).Result()
}

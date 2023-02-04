package store

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

func (s *StoreDatabase) IsBlacklist(token string) (bool, error) {
	_, err := s.Store.Get(context.Background(), token).Result()
	if err == redis.Nil {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

func (s *StoreDatabase) SetBlacklist(token string, expiration time.Duration) error {
	err := s.Store.Set(context.Background(), token, 0, expiration).Err()
	return err
}

package store

import (
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type StoreDatabase struct {
	Store *redis.Client
}

func New(options *redis.Options) (*StoreDatabase, error) {
	fmt.Println("Connecting to Redis Store")
	rdb := redis.NewClient(options)
	if rdb == nil {
		return nil, errors.New("failed to connect redis")
	}
	return &StoreDatabase{Store: rdb}, nil
}

package auth

import (
	"time"

	"github.com/google/uuid"
)

type DB interface {
	SetUserTokenPair(uuid uuid.UUID, expireAt time.Time, tokenString string) error
}

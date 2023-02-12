package auth

import (
	"errors"

	"github.com/golang-jwt/jwt/v4"
)

var (
	ErrNoValidEmail = errors.New("auth: no valid email")
)

var (
	ErrTokenExpired   = jwt.ErrTokenExpired
	ErrTokenInValid   = errors.New("token is invalid")
	ErrTokenBlacklist = errors.New("token is blacklist")
)

package auth

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type Database interface {
}

type Store interface {
	IsBlacklist(token string) (bool, error)
	SetBlacklist(token string, experation time.Duration) error
}

type JwtAuth struct {
	DB     Database
	Store  Store
	Secret string
}

type JwtPayload struct {
	UserUUID uuid.UUID
	Username string
}

type AuthTokenClaims struct {
	UserUUID     uuid.UUID `json:"uid"`
	Username     string    `json:"username"`
	RefreshToken string    `json:"refresh_token"`
	jwt.RegisteredClaims
}

func (j *JwtAuth) GenerateJWT(jp *JwtPayload) (string, error) {
	refreshToken, err := createRandomToken()
	if err != nil {
		fmt.Println("jwt create random token error")
		return "", err
	}
	expirationTime := time.Now().Add(time.Hour)
	claims := &AuthTokenClaims{
		UserUUID:     jp.UserUUID,
		Username:     jp.Username,
		RefreshToken: refreshToken,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(j.Secret))
	if err != nil {
		fmt.Println("jwt signed string error")
		return "", err
	}
	return tokenString, nil
}

func (j *JwtAuth) ParseToken(c *gin.Context, claims *AuthTokenClaims) (*jwt.Token, error) {
	token, err := c.Request.Cookie("token")
	if err != nil {
		return nil, err
	}
	// parse cookie
	tokenString := token.Value
	jwtToken, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.Secret), nil
	})
	if err != nil {
		return nil, err
	}
	return jwtToken, nil
}

func (j *JwtAuth) AuthorizeToken(tokenString string) (*AuthTokenClaims, error) {
	claims := &AuthTokenClaims{}
	jwtToken, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.Secret), nil
	})
	if err != nil {
		return claims, err
	}

	// check token is valid
	if !jwtToken.Valid {
		return claims, jwt.ErrTokenMalformed
	}

	isExpired := claims.isExpired()
	if isExpired {
		return claims, jwt.ErrTokenExpired
	}

	// check blacklist
	isBlacklist, err := j.Store.IsBlacklist(tokenString)
	if err != nil {
		return claims, err
	}
	if isBlacklist {
		return claims, ErrTokenBlacklist
	}
	return claims, nil
}

func (j *JwtAuth) AuthorizeRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Request.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"message": "no token cookie",
				})
				return
			}
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "no valid token",
			})
			return
		}

		// parse cookie
		tokenString := token.Value

		claims, err := j.AuthorizeToken(tokenString)
		if err != nil {
			if err == ErrTokenInValid {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"message": "internal server error",
				})
				return
			} else if err == ErrTokenExpired {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"message": "token is expired",
				})
				return
			} else if err == ErrTokenBlacklist {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"message": "token is expired",
				})
				return
			} else {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"message": "internal server error",
				})
				return
			}
		}
		c.Set("claims", claims)
		c.Next()
	}
}

func (a *AuthTokenClaims) isExpired() bool {
	return time.Until(a.ExpiresAt.Time) < 0
}

// not length parameter, since this function only used for generate random token
func createRandomToken() (string, error) {
	bytes := make([]byte, 20)
	for i := range bytes {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(validLetters))))
		if err != nil {
			return "", err
		}
		bytes[i] = validLetters[num.Int64()]
	}
	return string(bytes), nil
}

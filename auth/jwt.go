package auth

import (
	"crypto/rand"
	"math/big"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/redis/go-redis/v9"
)

type Database interface {
}

type JwtAuth struct {
	DB      Database
	Session *redis.Client
	Secret  string
}

type JwtPayload struct {
	UserUUID string
	Username string
}

type AuthTokenClaims struct {
	UserUUID     string `json:"uid"`
	Username     string `json:"username"`
	RefreshToken string `json:"refresh_token"`
	jwt.RegisteredClaims
}

func (j *JwtAuth) GenerateJWT(jp *JwtPayload) (string, error) {
	refreshToken, err := createRandomToken()
	if err != nil {
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

	tokenString, err := token.SignedString(j.Secret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (j *JwtAuth) VerifyMiddleWare() gin.HandlerFunc {
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
		claims := &AuthTokenClaims{}
		jwtToken, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return j.Secret, nil
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "internal server error",
			})
			return
		}

		// check token is valid
		if !jwtToken.Valid {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "token is invalid",
			})
		}

		isExpired := claims.isExpired()
		if isExpired {
			// implement when access token is expired
			// need to implement check session to this token is in blacklist token
		}
		// implement when access token is expired
		c.Next()
	}
}

func (a *AuthTokenClaims) isExpired() bool {
	return time.Until(a.ExpiresAt.Time) > 0
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

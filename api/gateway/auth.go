package gateway

import (
	"github.com/PylonSchema/server/auth"
)

func (g *Gateway) Auth(tokenString string) (claims *auth.AuthTokenClaims, err error) {
	claims, err = g.JwtAuth.AuthorizeToken(tokenString)
	return
}

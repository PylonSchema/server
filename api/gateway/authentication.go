package gateway

import (
	"fmt"

	"github.com/PylonSchema/server/auth"
)

func (g *Gateway) Auth(tokenString string) (claims *auth.AuthTokenClaims, err error) {
	fmt.Println(g.JwtAuth)
	claims, err = g.JwtAuth.AuthorizeToken(tokenString)
	return
}

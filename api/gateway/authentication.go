package gateway

func (g *Gateway) Auth(tokenString string) error {
	_, err := g.jwtAuth.AuthorizeToken(tokenString)
	if err != nil {
		return err
	}
	return nil
}

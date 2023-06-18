package auth

import "fmt"

type JwkRes interface {
	JwksUrl() string
}

type CognitoJwks struct{}

func NewCognitoJwks() JwkRes {
	return &CognitoJwks{}
}

func (c *CognitoJwks) JwksUrl() string {
	jwkUrl := fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s/.well-known/jwks.json",
		Settings.Region,
		Settings.UserPool,
	)
	return jwkUrl
}

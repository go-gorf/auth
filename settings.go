package auth

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
)

var client *cognitoidentityprovider.Client
var cognitoCtx context.Context

type authSettings struct {
	ClientId       string
	UserPool       string
	Region         string
	UserObjId      string
	AuthMiddleware *JwtAuthMiddleware
}

var Settings = authSettings{
	ClientId:       "",
	UserPool:       "",
	Region:         "",
	UserObjId:      "user",
	AuthMiddleware: nil,
}

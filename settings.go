package auth

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
)

var client *cognitoidentityprovider.Client
var cognitoCtx context.Context

type authSettings struct {
	ClientId  string
	UserPool  string
	Region    string
	UserObjId string
}

var Settings = authSettings{
	ClientId:  "",
	UserPool:  "",
	Region:    "",
	UserObjId: "user",
}

package auth

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
)

var client *cognitoidentityprovider.Client
var ctx context.Context

type AuthSettings struct {
	ClientId string
	UserPool string
	Region   string
}

var Settings = AuthSettings{}

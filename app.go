package auth

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/go-gorf/gorf"
)

func setup() error {
	cognitoCtx = context.Background()
	defaultConfig, err := config.LoadDefaultConfig(cognitoCtx, config.WithRegion(Settings.Region))
	if err != nil {
		return err
	}
	client = cognitoidentityprovider.NewFromConfig(defaultConfig)
	if Settings.AuthMiddleware == nil {
		return errors.New("no auth middleware present")
	}
	return nil
}

var App = gorf.BaseApp{
	Title:        "Auth",
	Info:         "Gorf Auth handler app",
	RouteHandler: Urls,
	SetUpHandler: setup,
}

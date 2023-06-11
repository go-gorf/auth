package auth

import (
	"context"
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
	return nil
}

var App = gorf.BaseApp{
	Title:        "Auth",
	Info:         "Gorf Auth handler app",
	RouteHandler: Urls,
	SetUpHandler: setup,
}

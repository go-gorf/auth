package auth

import (
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/gin-gonic/gin"
	"github.com/go-gorf/gorf"
)

func UserLogin(ctx *gin.Context) {
	loginInput := &LoginInput{}
	err := ctx.Bind(loginInput)
	if err != nil {
		gorf.BadRequest(ctx, "unable to parse login input", err)
		return
	}

	authInput := &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: types.AuthFlowTypeUserPasswordAuth,
		ClientId: &Settings.ClientId,
		AuthParameters: map[string]string{
			"USERNAME": loginInput.Email,
			"PASSWORD": loginInput.Password,
		},
	}

	result, err := client.InitiateAuth(cognitoCtx, authInput)
	if err != nil {
		gorf.BadRequest(ctx, "failed to authenticate", err)
		return
	}

	gorf.Response(ctx, gin.H{
		"TokenType":    *result.AuthenticationResult.TokenType,
		"ExpiresIn":    result.AuthenticationResult.ExpiresIn,
		"AccessToken":  *result.AuthenticationResult.AccessToken,
		"RefreshToken": *result.AuthenticationResult.RefreshToken,
	})
}

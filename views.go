package auth

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/gin-gonic/gin"
	"github.com/go-gorf/gorf"
	"github.com/go-gorf/gorf/common"
)

func UserLogin(ctx *gin.Context) {
	loginInput := &common.LoginInput{}
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
	//client.RevokeToken()
	//client.GetUser()
	//client.VerifySoftwareToken()
	//client.AssociateSoftwareToken()
}

func ProtectedApi(ctx *gin.Context) {
	//todo fix format
	gorf.Response(ctx, gin.H{
		"Status": "ok",
	})
}

func UserLogout(ctx *gin.Context) {

	err := Settings.AuthMiddleware.ParseAuthHeader(ctx)
	if err != nil {
		gorf.BadRequest(ctx, "error on parsing header", err)
		return
	}

	tokenInput := &cognitoidentityprovider.RevokeTokenInput{
		ClientId: &Settings.ClientId,
		Token:    Settings.AuthMiddleware.GetTokenStr(),
	}

	result, err := client.RevokeToken(cognitoCtx, tokenInput)
	if err != nil {
		gorf.BadRequest(ctx, "failed to revoke token", err)
		return
	}

	fmt.Println(result.ResultMetadata)

	//fix format
	gorf.Response(ctx, gin.H{
		"Status": "logout successfully",
	})
}

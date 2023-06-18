package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/go-gorf/gorf"
	"net/http"
)

type Middleware interface {
	ParseAuthHeader(ctx *gin.Context) error
	ParseJwtToken(ctx *gin.Context) error
	GetUser(ctx *gin.Context) (*gorf.BaseUser, error)
	Authenticate(ctx *gin.Context) (*gorf.BaseUser, error)
	GetTokenStr() *string
}

func AuthenticationRequiredMiddleware(ctx *gin.Context) {
	user, err := Settings.AuthMiddleware.Authenticate(ctx)
	if err != nil {
		e := gorf.NewErr("failed to authenticate", http.StatusUnauthorized, err)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, e.Response())
		return
	}
	ctx.Set(gorf.Settings.UserObjKey, user)
	ctx.Next()
}

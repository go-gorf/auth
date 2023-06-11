package auth

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-gorf/gorf"
)

func UserLogin(ctx *gin.Context) {
	fmt.Println("UserLogin")
	loginInput := &LoginInput{}
	err := ctx.Bind(loginInput)
	if err != nil {
		gorf.BadRequest(ctx, "Unable to parse login input", err)
		return
	}

	gorf.Response(ctx, gin.H{
		"email":    loginInput.Email,
		"password": loginInput.Password,
	})
}

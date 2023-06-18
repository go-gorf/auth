package auth

import "github.com/gin-gonic/gin"

func Urls(r *gin.Engine) {
	r.POST("auth/login", UserLogin)
	r.GET("auth/logout", AuthenticationRequiredMiddleware, UserLogout)
	r.GET("auth/protected-api", AuthenticationRequiredMiddleware, ProtectedApi)
}

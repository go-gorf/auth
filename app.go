package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/go-gorf/gorf"
	"log"
)

type App struct {
	name         string
	routeHandler func(r *gin.Engine)
}

func (*App) Setup() {
	println("Configuring the Auth app")
	err := gorf.DB.AutoMigrate(&User{})
	if err != nil {
		log.Fatal("Failed to migrate User model")
		return
	}
}

func (app *App) Register(r *gin.Engine) {
	app.routeHandler(r)
}

var AuthApp = App{
	name:         "auth",
	routeHandler: Urls,
}

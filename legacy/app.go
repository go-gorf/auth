package auth

import (
	"github.com/go-gorf/gorf"
)

func setup() error {
	err := gorf.DB.AutoMigrate(&User{})
	if err != nil {
		return err
	}
	return nil
}

var AuthApp = gorf.GorfBaseApp{
	Name:         "auth",
	RouteHandler: Urls,
	SetUpHandler: setup,
}

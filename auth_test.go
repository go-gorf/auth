package auth

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-gorf/gorf"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

// add all the apps
var apps = []gorf.GorfApp{
	&AuthApp,
}

func LoadSettings() {
	// jwt secret key
	gorf.Settings.SecretKey = "GOo8Rs8ht7qdxv6uUAjkQuopRGnql2zWJu08YleBx6pEv0cQ09a"
}

// bootstrap server
func BootstrapRouter() *gin.Engine {
	gorf.Apps = append(apps)
	LoadSettings()
	dbConf := gorf.DbConf{
		"data.db",
	}
	gorf.InitializeDatabase(&dbConf)
	gorf.SetupApps()
	r := gin.Default()
	gorf.RegisterApps(r)
	return r
}

func TestNewUserHandler(t *testing.T) {
	r := BootstrapRouter()
	newUser := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{
		"tom@example.com",
		"asd",
	}

	jsonValue, _ := json.Marshal(newUser)
	req, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer(jsonValue))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
}

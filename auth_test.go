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

var router *gin.Engine = nil

// add all the apps
var apps = []gorf.GorfApp{
	&AuthApp,
}

func LoadSettings() {
	// jwt secret key
	gorf.Settings.SecretKey = "GOo8Rs8ht7qdxv6uUAjkQuopRGnql2zWJu08YleBx6pEv0cQ09a"
	gorf.Settings.DbConf = &gorf.SqliteBackend{
		Name: "db.sqlite",
	}
}

// bootstrap server
func BootstrapRouter() *gin.Engine {
	gorf.Apps = append(apps)
	LoadSettings()
	gorf.InitializeDatabase()
	gorf.SetupApps()
	router = gin.Default()
	gorf.RegisterApps(router)
	return router
}

func GetRouter() *gin.Engine {
	if router != nil {
		return router
	}
	router = BootstrapRouter()
	return router
}

const testMail = "test@example.com"
const testPassword = "asd123"

func TestNewUserHandler(t *testing.T) {
	r := GetRouter()
	newUser := struct {
		Email     string `json:"email"`
		Password  string `json:"password"`
		FirstName string `json:"first_name"`
	}{
		testMail,
		testPassword,
		"toms",
	}

	jsonValue, _ := json.Marshal(newUser)
	req, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer(jsonValue))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, testMail, response["email"])
}

func TestInvalidPassNewUserHandler(t *testing.T) {
	r := GetRouter()
	newUser := struct {
		Email string `json:"email"`
	}{
		testMail,
	}

	jsonValue, _ := json.Marshal(newUser)
	req, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer(jsonValue))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestInvalidMailNewUserHandler(t *testing.T) {
	r := GetRouter()
	newUser := struct {
		Password string `json:"password"`
	}{
		testPassword,
	}

	jsonValue, _ := json.Marshal(newUser)
	req, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer(jsonValue))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

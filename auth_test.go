package auth

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-gorf/gorf"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strconv"
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

type RandUser struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
}

func (u *RandUser) isEmpty() bool {
	if u.Email == "" {
		return true
	}
	if u.Password == "" {
		return true
	}
	return false
}

func RandonStr() string {
	return strconv.Itoa(int(rand.Uint32()))
}

func (u *RandUser) populate() {
	u.Email = RandonStr()
	u.Password = RandonStr()
	u.FirstName = RandonStr()
}

func NewTestUser() RandUser {
	user := RandUser{}
	user.populate()
	return user
}

var test_user RandUser = NewTestUser()

func TestNewUserHandler(t *testing.T) {
	r := GetRouter()

	jsonValue, _ := json.Marshal(test_user)
	req, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer(jsonValue))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, test_user.Email, response["email"])
}

func TestInvalidPassNewUserHandler(t *testing.T) {
	r := GetRouter()
	newUser := struct {
		Email string `json:"email"`
	}{
		test_user.Email,
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
		test_user.Password,
	}

	jsonValue, _ := json.Marshal(newUser)
	req, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer(jsonValue))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLoginHandler(t *testing.T) {
	r := GetRouter()
	newUser := NewTestUser()
	jsonValue, _ := json.Marshal(newUser)
	req, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer(jsonValue))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, newUser.Email, response["email"])

	// try to login

	req, _ = http.NewRequest("POST", "/login", bytes.NewBuffer(jsonValue))

	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestLoginInvalidPassHandler(t *testing.T) {
	r := GetRouter()
	newUser := NewTestUser()

	jsonValue, _ := json.Marshal(newUser)
	req, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer(jsonValue))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, newUser.Email, response["email"])

	// try to login
	newUser.Password = "invalid"
	jsonValue, _ = json.Marshal(newUser)
	req, _ = http.NewRequest("POST", "/login", bytes.NewBuffer(jsonValue))

	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUniqueUser(t *testing.T) {
	r := GetRouter()
	newUser := NewTestUser()

	// create new user
	jsonValue, _ := json.Marshal(newUser)
	req, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer(jsonValue))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// try to recreate same user
	req, _ = http.NewRequest("POST", "/signup", bytes.NewBuffer(jsonValue))

	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var result map[string]string

	err := json.Unmarshal(w.Body.Bytes(), &result)
	if err != nil {
		return
	}
	assert.Equal(t, "User with same email already exists", result["message"])
}

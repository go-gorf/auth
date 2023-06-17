package auth

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-gorf/gorf"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
	"time"
)

type Middleware interface {
	ParseAuthHeader(ctx *gin.Context) error
	ParseJwtToken() error
	GetUser(ctx *gin.Context) (*gorf.BaseUser, error)
	Authenticate(ctx *gin.Context) (*gorf.BaseUser, error)
}

type JwtAuthMiddleware struct {
	tokenString string
	claims      jwt.MapClaims
	db          gorf.Db
}

func (m *JwtAuthMiddleware) ParseAuthHeader(ctx *gin.Context) error {
	authHead := ctx.Request.Header.Get("Authorization")
	if authHead == "" {
		return errors.New("no valid header")
	}
	jwtArr := strings.Split(authHead, " ")
	if len(jwtArr) < 2 {
		return errors.New("no jwt provided")
	}
	m.tokenString = jwtArr[1]
	return nil
}

func (m *JwtAuthMiddleware) ParseJwtToken() error {
	token, _ := jwt.Parse(m.tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(gorf.Settings.SecretKey), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		m.claims = claims
		return nil
	}

	return errors.New("unable to parse token")
}

func (m *JwtAuthMiddleware) GetUser(ctx *gin.Context) (*gorf.BaseUser, error) {
	user := &gorf.BaseUser{}
	err := m.db.GetUser(user, m.claims[Settings.UserObjId].(string))
	if err != nil {
		//TODO: use gorf impl
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid User",
		})
		return nil, errors.New("fdfd")
	}
	return user, nil
}

func (m *JwtAuthMiddleware) Authenticate(ctx *gin.Context) (*gorf.BaseUser, error) {
	err := m.ParseAuthHeader(ctx)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	err = m.ParseJwtToken()
	if err != nil {
		return nil, errors.New(err.Error())
	}

	exp, err := m.claims.GetExpirationTime()
	if err != nil {
		return nil, errors.New(err.Error())
	}

	if time.Now().Unix() > exp.Unix() {
		return nil, errors.New("jwt token expired")
	}

	user, err := m.GetUser(ctx)
	if err != nil {
		return nil, errors.New("unable to get User")
	}
	return user, nil
}

func AuthenticationRequiredMiddleware(ctx *gin.Context) {
	jwtAuth := JwtAuthMiddleware{}
	user, err := jwtAuth.Authenticate(ctx)

	if err != nil {
		e := gorf.NewErr("failed to authenticate", http.StatusUnauthorized, err)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, e.Response())
		return
	}
	ctx.Set(gorf.Settings.UserObjKey, user)
	ctx.Next()
}

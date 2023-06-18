package auth

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-gorf/gorf"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"net/http"
	"strings"
	"time"
)

type JwtAuthMiddleware struct {
	tokenString string
	claims      jwt.MapClaims
	DB          gorf.Db
	jwkRes      JwkRes
}

func NewJwtMiddleware(database gorf.Db, jwkRes JwkRes) *JwtAuthMiddleware {
	return &JwtAuthMiddleware{DB: database, jwkRes: jwkRes}
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
func (m *JwtAuthMiddleware) GetTokenStr() *string {
	return &m.tokenString
}

func (m *JwtAuthMiddleware) ParseJwtToken(ctx *gin.Context) error {
	//todo move to new res
	keySet, err := jwk.Fetch(ctx, m.jwkRes.JwksUrl())
	if err != nil {
		return err
	}
	token, _ := jwt.Parse(m.tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, errors.New("kid header not found")
		}
		keys, ok := keySet.LookupKeyID(kid)
		if !ok {
			return nil, fmt.Errorf("key with specified kid is not present in jwks")
		}
		var publickey interface{}
		err = keys.Raw(&publickey)
		if err != nil {
			return nil, fmt.Errorf("could not parse pubkey")
		}
		return publickey, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		m.claims = claims
		return nil
	}

	return errors.New("unable to parse token")
}

func (m *JwtAuthMiddleware) GetUser(ctx *gin.Context) (*gorf.BaseUser, error) {
	user := &gorf.BaseUser{}
	err := m.DB.GetUser(user, m.claims[Settings.UserObjId].(string))
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

	err = m.ParseJwtToken(ctx)
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

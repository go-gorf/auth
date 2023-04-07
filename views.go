package auth

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-gorf/gorf"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

func UserSignUp(ctx *gin.Context) {
	var body struct {
		Email     string `binding:"required"`
		Password  string `binding:"required"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}
	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Error while creating password",
		})
		return
	}

	newUser := User{
		Email:     body.Email,
		Password:  string(passwordHash),
		FirstName: body.FirstName,
		LastName:  body.LastName,
		Active:    bool(AuthSettings.NewUserState),
		Admin:     bool(AuthSettings.NewUserAdminState),
	}
	result := gorf.DB.Create(&newUser)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err,
		})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"user_id": newUser.ID,
		"email":   newUser.Email,
		"name":    newUser.FullName(),
	})
}

func UserLogin(ctx *gin.Context) {
	var body struct {
		Email    string `binding:"required"`
		Password string `binding:"required"`
	}
	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	user := User{}
	gorf.DB.First(&user, "email = ?", body.Email)
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Password",
		})
		return
	}

	if !user.Active {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "User is not active",
		})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		gorf.Settings.UserObjKey: user.ID,
		"timestamp":              time.Now(),
		"exp":                    time.Now().Add(time.Hour * 24).Unix(),
	})

	token_string, err := token.SignedString([]byte(gorf.Settings.SecretKey))

	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Unable to generate jwt token",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "JWT token generated successfully",
		"jwt":     token_string,
	})
}

func UserProfile(ctx *gin.Context) {
	user, _ := GetUser(ctx)
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Success",
		"email":   user.Email,
	})
}

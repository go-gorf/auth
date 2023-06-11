package auth

import (
	"fmt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Email     string `json:"email" binding:"required" gorm:"unique;not null"`
	Password  string `json:"-" binding:"required"`
	Active    bool   `json:"active"`
	Admin     bool   `json:"admin"`
}

func (user *User) FullName() string {
	return fmt.Sprintf("%v %v", user.FirstName, user.LastName)
}

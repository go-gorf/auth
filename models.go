package auth

import "gorm.io/gorm"

type User struct {
	gorm.Model
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Email     string `json:"email" binding:"required"`
	Password  string `json:"-" binding:"required"`
	Active    bool   `json:"active"`
	Admin     bool   `json:"admin"`
}

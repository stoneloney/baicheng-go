package models

import (
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Model
	Name   		string  `form:"email" binding:"required"`
	Password	string  `form:"password" binding:"required"`
	Logintime   string  `from:"logintime"` 
}

type Login struct {
	Name      string   `form:"name" binding:"required"`
	Password  string   `form:"password" binding:"required"`
}

// gorm beforesave hook
func (u *User) BeforeSave() (err error) {
	var hash []byte
	hash, err = bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return
	}
	u.Password = string(hash)
	return
}

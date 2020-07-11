package main

import (
	"github.com/jinzhu/gorm"
	"github.com/rs/xid"
)

type RefreshToken struct {
	Username string
	RefreshTokenID xid.ID
}

type User struct {
	gorm.Model
	Name string
	Username string
	Password string
	Salt string
	Email string
}

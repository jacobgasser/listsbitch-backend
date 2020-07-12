package main

import (
	"github.com/jinzhu/gorm"
	"github.com/rs/xid"
	"time"
)

type RefreshToken struct {
	Username string
	RefreshTokenID xid.ID
	UpdatedAt time.Time
}

type User struct {
	gorm.Model
	Name string
	Username string
	Password string
	Salt string
	Email string
}

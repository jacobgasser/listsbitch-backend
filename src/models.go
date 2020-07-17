package main

import (
	"github.com/jinzhu/gorm"
	"time"
)

type RefreshToken struct {
	Username string
	RefreshTokenID string
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

type List struct {
	gorm.Model
	Author User
	Title string
}

type ListItem struct {
	gorm.Model
	ListId string
}
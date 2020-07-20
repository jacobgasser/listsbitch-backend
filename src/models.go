package main

import (
	"time"
)

type Model struct {
	ID        string `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

type RefreshToken struct {
	User string
	RefreshTokenID string
	UpdatedAt time.Time
}

type User struct {
	Model
	Name string
	Username string
	Password string
	Salt string
	Email string
}

type List struct {
	Model
	Author string
	Title string
}

type ListItem struct {
	Model
	ListID string `json:"list_id"`
	Content string `json:"content"`
}
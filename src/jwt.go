package main

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/rs/xid"
)

var users = map[string]string{
	"user1": "password1",
	"user2": "password2",
}

var jwtKey = []byte("asdf")

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email string `json:"email"`
	Name string `json:"name"`
}

type AuthClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type RefreshClaims struct {
	RefreshId xid.ID `json:"refresh_ID"`
	jwt.StandardClaims
}


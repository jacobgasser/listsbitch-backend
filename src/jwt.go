package main

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/rs/xid"
	"net/http"
	"time"
)

var JwtKey = []byte("asdf")

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

func SetAuthJWT(w http.ResponseWriter, creds Credentials) {
	expirationTime := time.Now().Add(5 * time.Minute)

	claims := &AuthClaims{
		Username: creds.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JwtKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	http.SetCookie(w, &http.Cookie{
		Name: "token",
		Value: tokenString,
		Expires: expirationTime,
	})
	SetRefreshJWT(w, creds)
	w.WriteHeader(http.StatusOK)
}

func SetRefreshJWT(w http.ResponseWriter, creds Credentials) {
	expirationTime := time.Now().Add(72 * time.Hour)
	refreshID := xid.New()

	claims := &RefreshClaims{
		RefreshId: refreshID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	refreshToken := &RefreshToken{Username: creds.Username, RefreshTokenID: refreshID, UpdatedAt: time.Now()}
	DB.Where("username=?", creds.Username).Delete(&RefreshToken{})
	DB.Create(refreshToken)

	http.SetCookie(w, &http.Cookie{
		Name: "refresh_token",
		Value: tokenString,
		Expires: expirationTime,
	})
}


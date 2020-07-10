package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/rs/xid"
	"net/http"
	"time"
)

func AuthHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	var creds Credentials

	creds.Username = r.PostFormValue("username")
	creds.Password = r.PostFormValue("password")

	expectedPassword, ok := users[creds.Username]
	if !ok || expectedPassword != creds.Password {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(5 * time.Minute)

	claims := &AuthClaims{
		Username: creds.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	http.SetCookie(w, &http.Cookie{
		Name: "token",
		Value: tokenString,
		Expires: expirationTime,
	})
	createRefreshToken(w, creds.Username)
	fmt.Fprintf(w, "Welcome, %s", creds.Username)
}

func createRefreshToken(w http.ResponseWriter, username string) {
	expirationTime := time.Now().Add(72 * time.Hour)
	refreshID := xid.New()

	claims := &RefreshClaims{
		RefreshId: refreshID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		println(err.Error())
		return
	}

	refreshToken := &RefreshToken{Username: username, RefreshTokenID: refreshID, }
	DB.Where("username=?", username).Delete(&RefreshToken{})
	DB.Create(refreshToken)

	http.SetCookie(w, &http.Cookie{
		Name: "refresh_token",
		Value: tokenString,
		Expires: expirationTime,
	})

}

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	
}

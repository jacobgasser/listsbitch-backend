package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/rs/xid"
	"net/http"
	"time"
)

func AuthHandler(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user := User{}

	DB.First(&user, "username = ?", creds.Username)

	hashedSuppliedPassword := append([]byte(creds.Password), []byte(user.Salt)...)
	shasher := sha256.New()
	shasher.Write(hashedSuppliedPassword)
	password := hex.EncodeToString(shasher.Sum(nil))
	if user.Password != password {
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
	w.WriteHeader(http.StatusOK)
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
		return
	}

	refreshToken := &RefreshToken{Username: username, RefreshTokenID: refreshID }
	DB.Where("username=?", username).Delete(&RefreshToken{})
	DB.Create(refreshToken)

	http.SetCookie(w, &http.Cookie{
		Name: "refresh_token",
		Value: tokenString,
		Expires: expirationTime,
	})
}

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user := &User{}
	user.Username = creds.Username
	user.Email = creds.Email
	user.Name = creds.Name
	user.CreatedAt = time.Now()

	salt := make([]byte, 8)
	 _, err = rand.Read(salt)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user.Salt = hex.EncodeToString(salt)

	password := sha256.New()
	password.Write(append([]byte(creds.Password), []byte(user.Salt)...))
	user.Password = hex.EncodeToString(password.Sum(nil))

	DB.Create(&user)

	w.WriteHeader(http.StatusOK)
}

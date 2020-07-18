package main

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"net/http"
	"os"
	"time"
)

var JwtKey = []byte(os.Getenv("JWTKEY"))

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
	RefreshId string `json:"refresh_ID"`
	jwt.StandardClaims
}

func SetAuthJWT(w http.ResponseWriter, creds Credentials) string {
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
		Path: "/",
	})
	return tokenString
}

func SetRefreshJWT(w http.ResponseWriter, creds Credentials) {
	expirationTime := time.Now().Add(72 * time.Hour)
	refreshID := uuid.New().String()

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
		Path: "/",
	})

}

func Refresh(w http.ResponseWriter, refreshTokenJWT string) (int, string) {
	refreshToken, err := jwt.ParseWithClaims(refreshTokenJWT, &RefreshClaims{}, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})
	claims, ok := refreshToken.Claims.(*RefreshClaims)
	if !ok || !refreshToken.Valid {
		return http.StatusUnauthorized, ""
	}
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return http.StatusUnauthorized, ""
		}
		return http.StatusBadRequest, ""
	}
	if !refreshToken.Valid {
		return http.StatusUnauthorized, ""
	}
	OldRefreshTokenEntry := &RefreshToken{}
	DB.First(&OldRefreshTokenEntry, "refresh_token_id = ?", claims.RefreshId)
	if OldRefreshTokenEntry.Username == "" {
		return http.StatusUnauthorized, ""
	}
	creds := Credentials{Username: OldRefreshTokenEntry.Username}
	SetRefreshJWT(w, creds)
	authJWT := SetAuthJWT(w, creds)
	return http.StatusOK, authJWT
}

func Authenticate(w http.ResponseWriter, authTokenJWT string, user *User) (int) {
	authToken, err := jwt.ParseWithClaims(authTokenJWT, &AuthClaims{}, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})
	claims, ok := authToken.Claims.(*AuthClaims)
	if !ok || !authToken.Valid {
		return http.StatusUnauthorized
	}
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return http.StatusUnauthorized
		}
		return http.StatusBadRequest
	}
	if !authToken.Valid {
		return http.StatusUnauthorized
	}
	user.Username = claims.Username
	return http.StatusOK
}
func Verify(w http.ResponseWriter, r *http.Request) (User, int) {
	user := User{}
	authTokenCookie, err := r.Cookie("token")
	authToken := ""
	if authTokenCookie != nil {
		authToken = authTokenCookie.Value
	}
	if err != nil {
		if err == http.ErrNoCookie {
			refreshTokenCookie, err2 := r.Cookie("refresh_token")
			if err2 != nil {
				if err2 == http.ErrNoCookie {
					return user, http.StatusUnauthorized
				}
				return user, http.StatusInternalServerError
			}
			status, authJWT := Refresh(w, refreshTokenCookie.Value)
			if status != http.StatusOK {
				return user, status
			}
			authToken = authJWT
		}
	}
	status := Authenticate(w, authToken, &user)
	return user, status
}


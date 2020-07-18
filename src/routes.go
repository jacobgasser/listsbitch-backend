package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/google/uuid"
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

	if creds.Username == "" || creds.Password == ""{
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
	SetAuthJWT(w, creds)
	SetRefreshJWT(w, creds)
	w.WriteHeader(http.StatusOK)
}

func RefreshHandler(w http.ResponseWriter, r *http.Request) {
	refreshTokenCookie, err := r.Cookie("refresh_token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	status, _ := Refresh(w, refreshTokenCookie.Value)
	w.WriteHeader(status)
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
	user.ID = uuid.New().String()

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

func CreateListHandler(w http.ResponseWriter, r *http.Request) {
	user, status := Verify(w, r)
	if status != http.StatusOK {
		w.WriteHeader(status)
		return
	}
	if user.Username == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	body := struct {
		Name string `json:"name"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := CreateList(user.Username, body.Name); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

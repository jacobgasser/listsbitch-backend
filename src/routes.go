package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
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
	SetAuthJWT(w, user)
	SetRefreshJWT(w, user)
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

	if DB.Create(&user).Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func CreateListHandler(w http.ResponseWriter, r *http.Request) {
	user, status := Verify(w, r)
	if status != http.StatusOK {
		w.WriteHeader(status)
		return
	}
	if user.ID == "" {
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
	if err := CreateList(user.ID, body.Name); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func CreateListItemHandler(w http.ResponseWriter, r *http.Request) {
	user, status := Verify(w, r)
	if status != http.StatusOK {
		w.WriteHeader(status)
		return
	}
	if user.ID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	listItem := ListItem{}
	err := json.NewDecoder(r.Body).Decode(&listItem)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = CreateListItem(user.ID, listItem.ListID, listItem.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func GetListItemsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	listItems, err := GetListItems(vars["listID"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	for i, v := range listItems {
		fmt.Printf("ID: %x; Content: %s\n", i, v.Content)
	}
	listItemsJSON, err := json.Marshal(listItems)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = fmt.Fprintf(w, "%s", listItemsJSON)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func RemoveListItemHandler(w http.ResponseWriter, r *http.Request) {
	user, status := Verify(w, r)
	if status != http.StatusOK {
		w.WriteHeader(status)
		return
	}
	vars := mux.Vars(r)
	listItemID := vars["listItemID"]
	listItemToDelete := ListItem{}
	DB.First(&listItemToDelete, "id = ?", listItemID)
	listHoldingItemToDelete := List{}
	DB.First(&listHoldingItemToDelete, "id = ?", listItemToDelete.ListID)
	if listHoldingItemToDelete.Author != user.ID {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	err := DeleteListItem(listItemID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func DeleteListHandler(w http.ResponseWriter, r *http.Request) {
	user, status := Verify(w, r)
	if status != http.StatusOK {
		w.WriteHeader(status)
		return
	}
	vars := mux.Vars(r)
	listID := vars["listID"]
	listToDelete := List{}
	if DB.First(&listToDelete, "id = ?", listID).RecordNotFound() {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if listToDelete.Author != user.ID {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	err := DeleteList(listID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func GetUserListsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userID"]
	lists, err := GetUserLists(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	listsJSON, err := json.Marshal(lists)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = fmt.Fprintf(w, "%s", listsJSON)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

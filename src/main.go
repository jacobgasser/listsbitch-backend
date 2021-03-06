package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

var DB *gorm.DB

func main() {
	fmt.Print("Starting Listsbit.ch . . .\n")
	r := mux.NewRouter()
	r.HandleFunc("/login", AuthHandler).Methods("POST")
	r.HandleFunc("/signup", SignUpHandler).Methods("POST")
	r.HandleFunc("/refresh", RefreshHandler).Methods("POST")
	r.HandleFunc("/api/lists", CreateListHandler).Methods("POST")
	r.HandleFunc("/api/listitems", CreateListItemHandler).Methods("POST")
	r.HandleFunc("/api/listitems/{listID}", GetListItemsHandler).Methods("GET")
	r.HandleFunc("/api/listitems/{listItemID}", RemoveListItemHandler).Methods("DELETE")
	r.HandleFunc("/api/lists/{listID}", DeleteListHandler).Methods("DELETE")
	r.HandleFunc("/api/users/{userID}/lists", GetUserListsHandler).Methods("GET")
	http.Handle("/", r)

	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open("mysql", os.Getenv("DB_USER")+":"+os.Getenv("DB_PASSWORD")+"@("+os.Getenv("DB_HOST")+")/"+os.Getenv("DB_NAME")+"?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	DB = db
	DB.AutoMigrate(&User{})
	DB.AutoMigrate(&RefreshToken{})
	DB.AutoMigrate(&List{})
	DB.AutoMigrate(&ListItem{})
	defer DB.Close()
	
	fmt.Print("Listsbit.ch is now online!")
	log.Fatal(http.ListenAndServe(":8080", handlers.CORS(handlers.AllowedOrigins([]string{"*"}), handlers.AllowCredentials())(r)))

}

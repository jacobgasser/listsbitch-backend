package main

import (
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

var DB *gorm.DB

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/login", AuthHandler).Methods("POST")
	r.HandleFunc("/signup", SignUpHandler).Methods("POST")
	http.Handle("/", r)

	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open("mysql", os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASSWORD") + "@(" + os.Getenv("DB_HOST") + ")/" + os.Getenv("DB_NAME") + "?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	DB = db
	DB.AutoMigrate(&User{})
	DB.AutoMigrate(&RefreshToken{})
	defer DB.Close()
	log.Fatal(http.ListenAndServe(":8080", nil))


}



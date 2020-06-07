package main

import (
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"

	auth "./controllers"
)

func main() {

	http.HandleFunc("/", auth.UserMain)
	http.HandleFunc("/signup", auth.UserSignUp)
	http.HandleFunc("/login", auth.UserLogIn)
	http.HandleFunc("/logout", auth.UserLogOut)

	log.Println("Wake up, Neo...")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal(err)
	}
}

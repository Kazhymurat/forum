package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"

	handlers "./cookies/handlers"
)

type user struct {
	ID       int
	email    string
	username string
	password string
}

var database *sql.DB

//UserSignUp is a handler to add user credentials to database
func UserSignUp(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {

		err := r.ParseForm()
		if err != nil {
			log.Println(err)
		}
		email := r.FormValue("email")
		username := r.FormValue("username")
		password := r.FormValue("password")

		_, err = database.Exec("insert into user (email, username, password) values ($1, $2, $3)",
			email, username, password)

		if err != nil {
			log.Println(err)
		}
		http.Redirect(w, r, "/login", 301)

	} else {
		http.ServeFile(w, r, "static/signup.html")
	}
}

//UserLogIn is a handler to log in user with its appropriate credentials
func UserLogIn(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {

		email := r.FormValue("email")

		row := database.QueryRow("select * from user where email = $1", email)
		user := user{}
		err := row.Scan(&user.ID, &user.email, &user.username, &user.password)

		if err != nil {
			fmt.Fprint(w, "WRONG EMAIL")
		}
		http.Redirect(w, r, "/", 301)
		http.HandleFunc("/", userSession)

	} else {
		http.ServeFile(w, r, "static/login.html")
	}
}

//Cookies for session
func userSession(w http.ResponseWriter, r *http.Request) {

	handler := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte("OK"))
	})

	loginHandler := handlers.Login{
		Name:   "session",
		Value:  "logged in",
		MaxAge: 60,
		Next:   handler,
	}

	logoutHandler := handlers.Logout{
		Name: "logout",
		Next: handler,
	}

	http.Handle("/", loginHandler)
	http.Handle("/logout", logoutHandler)

}

func main() {

	fs := http.FileServer(http.Dir("./static"))

	db, err := sql.Open("sqlite3", "database/forum.db")

	if err != nil {
		log.Println(err)
	}
	database = db
	defer db.Close()

	http.Handle("/", fs)
	http.HandleFunc("/signup", UserSignUp)
	http.HandleFunc("/login", UserLogIn)

	// http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
	// 	http.ServeFile(w, r, "static/login.html")
	// })

	log.Println("Listening on :8000...")
	err = http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type user struct {
	ID       int
	email    string
	username string
	password string
}

var database *sql.DB

func hashAndSalt(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}

func comparePasswords(hashedPwd string, plainPwd []byte) bool {
	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

//UserSignUp is a handler to add user credentials to database
func UserSignUp(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {

		err := r.ParseForm()
		if err != nil {
			log.Println(err)
		}
		email := r.FormValue("email")
		username := r.FormValue("username")
		password := hashAndSalt([]byte(r.FormValue("password")))

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

		password := []byte(r.FormValue("password"))

		if comparePasswords(user.password, password) {
			// http.Redirect(w, r, "/", 301)
			fmt.Fprint(w, "Mission completed")

		} else {
			//WP - Wrong Password
			WP := true
			tmpl, _ := template.ParseFiles("static/login.html")
			tmpl.Execute(w, WP)
		}

	} else {
		// http.ServeFile(w, r, "static/login.html")
		WP := false
		tmpl, _ := template.ParseFiles("static/login.html")
		tmpl.Execute(w, WP)
	}
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

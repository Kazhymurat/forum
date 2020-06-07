package auth

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"path"
	"text/template"
	"time"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

type user struct {
	ID       int
	email    string
	username string
	password string
}

type session struct {
	Username string
	Auth     bool
}

var cookies map[int]*http.Cookie

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

	db, err := sql.Open("sqlite3", "database/forum.db")

	if err != nil {
		log.Println(err)
	}
	database = db
	defer db.Close()

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

//UserMain is a handler to execute main page
func UserMain(w http.ResponseWriter, r *http.Request) {

	id, auth := IsAuthorized(r)

	if auth {
		db, err := sql.Open("sqlite3", "database/forum.db")

		if err != nil {
			log.Println(err)
		}
		database = db
		defer db.Close()

		row := database.QueryRow("select * from user where ID = $1", id)
		user := user{}
		err = row.Scan(&user.ID, &user.email, &user.username, &user.password)

		if err != nil {
			fmt.Println(err)
		}

		data := session{
			Username: user.username,
			Auth:     true,
		}

		var (
			tmpl = template.Must(template.ParseFiles(path.Join("./static", "index.html")))
		)
		err = tmpl.Execute(w, data)
		if err != nil {
			log.Fatalf("execution failed: %s", err)
		}

	} else {

		data := session{
			Username: "",
			Auth:     false,
		}

		var (
			tmpl = template.Must(template.ParseFiles(path.Join("./static", "index.html")))
		)
		err := tmpl.Execute(w, data)
		if err != nil {
			log.Fatalf("execution failed: %s", err)
		}
	}
}

//UserLogIn is a handler to log in user with its appropriate credentials
func UserLogIn(w http.ResponseWriter, r *http.Request) {

	if cookies == nil {
		cookies = map[int]*http.Cookie{}
	}

	_, auth := IsAuthorized(r)

	//for my purpose to see map of cookies
	// fmt.Println(cookies)

	if auth {
		fmt.Fprint(w, "YOU ARE ALREADY IN")
	} else {

		db, err := sql.Open("sqlite3", "database/forum.db")

		if err != nil {
			log.Println(err)
		}
		database = db
		defer db.Close()

		if r.Method == "POST" {

			email := r.FormValue("email")

			row := database.QueryRow("select * from user where email = $1", email)
			user := user{}
			err := row.Scan(&user.ID, &user.email, &user.username, &user.password)

			if err != nil {
				fmt.Println(err)
			}

			password := []byte(r.FormValue("password"))

			if comparePasswords(user.password, password) {

				u, _ := uuid.NewV4()
				sessionToken := u.String()

				cookie := &http.Cookie{
					Name:    "session_token",
					Value:   sessionToken, // Some encoded value
					Path:    "/",          // Otherwise it defaults to the /login if you create this on /login (standard cookie behaviour)
					Expires: time.Now().Add(7200 * time.Second),
				}

				cookies[user.ID] = cookie
				http.SetCookie(w, cookie)
				fmt.Println("this is your cookie:", cookie)

				data := session{
					Username: user.username,
					Auth:     true,
				}

				var (
					tmpl = template.Must(template.ParseFiles(path.Join("./static", "index.html")))
				)
				err := tmpl.Execute(w, data)
				if err != nil {
					log.Fatalf("execution failed: %s", err)
				}

				// tmpl, _ := template.ParseFiles("./static/index.html")
				// tmpl.Execute(w, User)
				// fmt.Println(Auth)

				// http.Redirect(w, r, "/", 301)
				// fmt.Fprint(w, "Mission completed")

			} else {
				//WP - Wrong Password
				WP := true
				tmpl, _ := template.ParseFiles("./static/login.html")
				tmpl.Execute(w, WP)
			}

		} else {
			// http.ServeFile(w, r, "static/login.html")
			WP := false
			tmpl, _ := template.ParseFiles("./static/login.html")
			tmpl.Execute(w, WP)
		}
	}
}

//UserLogOut ends session for user
func UserLogOut(w http.ResponseWriter, r *http.Request) {

	ClearSession(w)
	http.Redirect(w, r, "/", 302)

}

//IsAuthorized validates is this session already taken
func IsAuthorized(r *http.Request) (int, bool) {
	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			return 0, false
		}
		return 0, false
	}
	for ID, cookie := range cookies {
		if cookie.Value == c.Value {
			if cookie.Expires.Sub(time.Now()) <= 0 {

				return 0, false
			} else {

				return ID, true
			}
		}
	}

	return 0, false
}

//ClearSession ends the session
func ClearSession(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   "session_token",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
}

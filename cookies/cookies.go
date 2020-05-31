package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	handlers "./handlers"
)

func main() {
	cookieName := os.Getenv("COOKIE_NAME")
	cookieValue := os.Getenv("COOKIE_VALUE")
	cookieDomain := os.Getenv("COOKIE_DOMAIN")
	cookiePath := os.Getenv("COOKIE_PATH")
	cookieDuration := os.Getenv("COOKIE_DURATION")

	cookieDurationInteger, err := strconv.Atoi(cookieDuration)

	if err != nil {
		panic(err)
	}

	handler := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte("OK"))
	})

	loginHandler := handlers.Login{
		Name:   cookieName,
		Domain: cookieDomain,
		Path:   cookiePath,
		Value:  cookieValue,
		MaxAge: cookieDurationInteger,
		Next:   handler,
	}

	logoutHandler := handlers.Logout{
		Name:   cookieName,
		Domain: cookieDomain,
		Path:   cookiePath,
		Next:   handler,
	}

	http.Handle("/login", loginHandler)
	http.Handle("/logout", logoutHandler)

	address := fmt.Sprintf(":%s", os.Getenv("PORT"))

	if err := http.ListenAndServe(address, nil); err != nil {
		panic(err)
	}
}

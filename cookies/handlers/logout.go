package handlers

import (
	"net/http"
)

// Logout is a handler that clears a logged in cookie
type Logout struct {
	Name   string
	Path   string
	Domain string
	Next   http.Handler
}

func (l Logout) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name:     l.Name,
		Value:    "",
		MaxAge:   0,
		HttpOnly: true,
	}
	http.SetCookie(rw, &cookie)
	l.Next.ServeHTTP(rw, r)
}

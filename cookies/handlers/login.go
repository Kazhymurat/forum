package handlers

import (
	"net/http"
)

// Login is a handler that sets a logged in cookie
type Login struct {
	Name   string
	Value  string
	Path   string
	Domain string
	MaxAge int
	Next   http.Handler
}

func (l Login) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name:     l.Name,
		Value:    l.Value,
		Domain:   l.Domain,
		Path:     l.Path,
		MaxAge:   l.MaxAge,
		HttpOnly: true,
	}
	http.SetCookie(rw, &cookie)
	l.Next.ServeHTTP(rw, r)
}

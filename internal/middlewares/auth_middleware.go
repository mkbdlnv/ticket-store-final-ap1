package middleware

import (
	"github.com/golangcollege/sessions"
	"log"
	"net/http"
)

type AuthMiddleware struct {
	Session *sessions.Session
}

func (am *AuthMiddleware) IsNotAuthorized(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if am.CheckSession(r) {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	}
}

func (am *AuthMiddleware) IsAuthorized(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !am.CheckSession(r) {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	}
}

func (am *AuthMiddleware) IsAdmin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if am.CheckAdmin(r) {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	}
}

func (am *AuthMiddleware) CheckSession(r *http.Request) bool {
	id := am.Session.Get(r, "userId")
	log.Println("user id: ", id)
	if id == nil {
		return false
	}
	return true

}

func (am *AuthMiddleware) CheckAdmin(r *http.Request) bool {
	role := am.Session.GetString(r, "userRole")
	log.Println("User role:::", role)
	if role != "admin" {
		return false
	}
	return true
}

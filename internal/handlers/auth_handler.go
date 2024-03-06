package handlers

import (
	"finalProject/internal/repositories"
	"finalProject/internal/utils"
	"github.com/golangcollege/sessions"
	"log"
	"net/http"
	"strconv"
)

type AuthHandler struct {
	Session  *sessions.Session
	UserRepo *repositories.UserRepository
}

func (uh *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			log.Println(w, "ParseForm() err: %v", err)
			return
		}
		email := r.FormValue("email")
		password := r.FormValue("password")
		user, _ := uh.UserRepo.GetByEmail(email)
		log.Println("Email: " + email)
		log.Println("Password: " + password)

		if user == nil {
			http.Redirect(w, r, "/login?invaliddataerr=true", http.StatusFound)
			return
		}

		if !utils.VerifyPassword(user.Password, password) {
			http.Redirect(w, r, "/login?invaliddataerr=true", http.StatusInternalServerError)
			return
		}

		uh.Session.Put(r, "userEmail", user.Email)
		uh.Session.Put(r, "userRole", user.Role)
		uh.Session.Put(r, "userName", user.Name)
		uh.Session.Put(r, "userLastname", user.Lastname)
		userId := strconv.Itoa(user.ID)
		uh.Session.Put(r, "userId", userId)
		if user.Role == "admin" {
			http.Redirect(w, r, "/admin", http.StatusFound)
		} else {
			http.Redirect(w, r, "/", http.StatusFound)
		}

		http.Redirect(w, r, "/", http.StatusFound)

	}
}

func (uh *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if uh.Session.Exists(r, "userEmail") {
		uh.Session.Remove(r, "userEmail")
	}
	if uh.Session.Exists(r, "userRole") {
		uh.Session.Remove(r, "userRole")
	}
	if uh.Session.Exists(r, "userId") {
		uh.Session.Remove(r, "userId")
	}
	if uh.Session.Exists(r, "userName") {
		uh.Session.Remove(r, "userName")
	}
	if uh.Session.Exists(r, "userLastname") {
		uh.Session.Remove(r, "userLastname")
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

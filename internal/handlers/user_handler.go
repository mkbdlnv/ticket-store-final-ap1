// File: internal/handlers/user_handler.go

package handlers

import (
	"finalProject/internal/repositories"
	"github.com/golangcollege/sessions"
	"log"
	"net/http"
)

type UserHandler struct {
	UserRepo *repositories.UserRepository
	Session  *sessions.Session
}

func (uh *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			log.Println(w, "ParseForm() err: %v", err)
			return
		}
		name := r.FormValue("name")
		lastname := r.FormValue("lastname")
		email := r.FormValue("email")
		password := r.FormValue("password")
		
		user, err := uh.UserRepo.GetByEmail(email)
		if user != nil {
			http.Redirect(w, r, "/registration?emailerr=true", http.StatusFound)
			return
		}
		
		id, err := uh.UserRepo.Insert(name, lastname, email, password)
		if checkError(err, &w, "User not created", 500) {
			http.Redirect(w, r, "/registration", http.StatusFound)
		}
		log.Printf("User with id: %d was created.\n", id)
		http.Redirect(w, r, "/login", http.StatusFound)
	}
}

func checkError(err error, w *http.ResponseWriter, msg string, status int) bool {
	if err != nil {
		log.Println(err.Error())
		http.Error(*w, msg, status)
	}
	return err != nil
}

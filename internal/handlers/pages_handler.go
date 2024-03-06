package handlers

import (
	"finalProject/internal/repositories"
	"github.com/golangcollege/sessions"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

type PageHandler struct {
	Session    *sessions.Session
	EventRepo  *repositories.EventRepository
	TicketRepo *repositories.TicketRepository
	UserRepo   *repositories.UserRepository
}

func (ph *PageHandler) Home(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	userIdStr := ph.Session.GetString(r, "userId")
	userId, _ := strconv.Atoi(userIdStr)

	events, err := ph.EventRepo.GetAll()
	if err != nil {
		return
	}

	userRole := ph.Session.GetString(r, "userRole")

	data := TemplateData{
		Events:   events,
		LoggedIn: userId != 0,
		IsAdmin:  userRole == "admin",
	}

	files := []string{

		"../../templates/main.page.tmpl",

		"../../templates/base.layout.tmpl",

		"../../templates/header.partial.tmpl",
		"../../templates/footer.partial.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if checkError(err, &w, "Internal Server Error", 500) {
		return
	}

	err = ts.Execute(w, data)
	checkError(err, &w, "Internal Server Error", 500)
}

func (ph *PageHandler) Event(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if checkError(err, &w, "Not found", 404) {
		return
	}

	event, err := ph.EventRepo.GetByID(id)

	if checkError(err, &w, "Not found", 404) {
		return
	}
	userIdStr := ph.Session.GetString(r, "userId")
	userId, _ := strconv.Atoi(userIdStr)
	userRole := ph.Session.GetString(r, "userRole")
	data := TemplateData{
		Event:      event,
		FormatDate: event.DateTime.Format("2006-01-02"),
		FormatTime: event.DateTime.Format("15:04:05"),
		LoggedIn:   userId != 0,
		IsAdmin:    userRole == "admin",
	}

	files := []string{

		"../../templates/event.page.tmpl",

		"../../templates/base.layout.tmpl",

		"../../templates/header.partial.tmpl",
		"../../templates/footer.partial.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if checkError(err, &w, "Internal Server Error", 500) {
		return
	}

	err = ts.Execute(w, data)
	checkError(err, &w, "Internal Server Error", 500)

}

func (ph *PageHandler) Admin(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"../../templates/add_event.page.tmpl",
		"../../templates/admin.layout.tmpl",

		"../../templates/header.partial.tmpl",
		"../../templates/footer.partial.tmpl",
	}
	userIdStr := ph.Session.GetString(r, "userId")
	userId, _ := strconv.Atoi(userIdStr)
	userRole := ph.Session.GetString(r, "userRole")
	data := TemplateData{
		LoggedIn: userId != 0,
		IsAdmin:  userRole == "admin",
	}

	ts, err := template.ParseFiles(files...)
	if checkError(err, &w, "Internal Server Error", 500) {
		return
	}

	err = ts.Execute(w, data)
	checkError(err, &w, "Internal Server Error", 500)
}

func (ph *PageHandler) User(w http.ResponseWriter, r *http.Request) {

	userIdStr := ph.Session.GetString(r, "userId")
	userId, _ := strconv.Atoi(userIdStr)

	tickets, _ := ph.TicketRepo.GetAllByUserId(userId)
	user, _ := ph.UserRepo.GetById(userId)

	userRole := ph.Session.GetString(r, "userRole")

	data := TemplateData{
		Tickets:  tickets,
		User:     user,
		LoggedIn: userId != 0,
		IsAdmin:  userRole == "admin",
	}

	files := []string{
		"../../templates/profile.page.tmpl",
		"../../templates/base.layout.tmpl",

		"../../templates/header.partial.tmpl",
		"../../templates/footer.partial.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if checkError(err, &w, "Internal Server Error", 500) {
		return
	}

	err = ts.Execute(w, data)
	checkError(err, &w, "Internal Server Error", 500)
}

func (ph *PageHandler) Registration(w http.ResponseWriter, r *http.Request) {
	files := []string{

		"../../templates/registration.page.tmpl",

		"../../templates/base.layout.tmpl",

		"../../templates/header.partial.tmpl",
		"../../templates/footer.partial.tmpl",
	}

	userIdStr := ph.Session.GetString(r, "userId")
	userId, _ := strconv.Atoi(userIdStr)
	userRole := ph.Session.GetString(r, "userRole")

	ts, err := template.ParseFiles(files...)

	if checkError(err, &w, "Internal Server Error", 500) {
		return
	}
	queryParams := r.URL.Query()
	emailError := queryParams.Get("emailerr")

	if emailError == "true" {
		data := ErrorData{
			ErrorMsg: "Email already registered",
		}
		log.Println("ERROR EMAIL")
		err = ts.Execute(w, data)
		return
	}
	data := TemplateData{
		LoggedIn: userId != 0,
		IsAdmin:  userRole == "admin",
	}
	err = ts.Execute(w, data)
	checkError(err, &w, "Internal Server Error", 500)
}

func (ph *PageHandler) Login(w http.ResponseWriter, r *http.Request) {
	files := []string{

		"../../templates/login.page.tmpl",

		"../../templates/base.layout.tmpl",

		"../../templates/header.partial.tmpl",
		"../../templates/footer.partial.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if checkError(err, &w, "Internal Server Error", 500) {
		return
	}

	queryParams := r.URL.Query()
	invaliDataError := queryParams.Get("invaliddataerr")

	if invaliDataError == "true" {
		data := ErrorData{
			ErrorMsg: "Incorrect email or password",
		}
		err = ts.Execute(w, data)
		return
	}

	userIdStr := ph.Session.GetString(r, "userId")
	userId, _ := strconv.Atoi(userIdStr)
	userRole := ph.Session.GetString(r, "userRole")

	data := TemplateData{
		LoggedIn: userId != 0,
		IsAdmin:  userRole == "admin",
	}

	err = ts.Execute(w, data)
	checkError(err, &w, "Internal Server Error", 500)
}

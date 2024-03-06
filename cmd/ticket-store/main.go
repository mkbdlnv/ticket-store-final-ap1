package main

import (
	conf "finalProject/internal/config"
	"finalProject/internal/db"
	"finalProject/internal/handlers"
	middleware "finalProject/internal/middlewares"
	"finalProject/internal/repositories"
	"fmt"
	"github.com/golangcollege/sessions"
	"log"
	"net/http"
	"os"
	"time"
)

type application struct {
	config         *conf.Config
	ErrorLog       *log.Logger
	InfoLog        *log.Logger
	userRepo       *repositories.UserRepository
	pageHandler    *handlers.PageHandler
	userHandler    *handlers.UserHandler
	authHandler    *handlers.AuthHandler
	eventHandler   *handlers.EventHandler
	ticketHandler  *handlers.TicketHandler
	authMiddleware *middleware.AuthMiddleware
	session        *sessions.Session
}

var App *application

func main() {

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	cnf := conf.GetConfig()
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable ",
		cnf.DB.Host, cnf.DB.Port, cnf.DB.User, cnf.DB.Password, cnf.DB.DbName)

	dbcon, err := db.OpenDB(psqlInfo)
	if err != nil {
		errorLog.Fatal(err)
	}

	session := sessions.New([]byte(cnf.SERVER.SessionSecret))
	session.Lifetime = 12 * time.Hour
	userRepo := &repositories.UserRepository{DB: dbcon}
	eventRepo := &repositories.EventRepository{DB: dbcon}
	ticketRepo := &repositories.TicketRepository{DB: dbcon}
	App = &application{
		config:         conf.GetConfig(),
		ErrorLog:       errorLog,
		InfoLog:        infoLog,
		userRepo:       userRepo,
		pageHandler:    &handlers.PageHandler{EventRepo: eventRepo, TicketRepo: ticketRepo, UserRepo: userRepo, Session: session},
		userHandler:    &handlers.UserHandler{UserRepo: userRepo},
		authHandler:    &handlers.AuthHandler{UserRepo: userRepo, Session: session},
		eventHandler:   &handlers.EventHandler{EventRepo: eventRepo, Session: session},
		ticketHandler:  &handlers.TicketHandler{TicketRepo: ticketRepo},
		authMiddleware: &middleware.AuthMiddleware{Session: session},
		session:        session,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", App.pageHandler.Home)
	mux.HandleFunc("/login", App.authMiddleware.IsNotAuthorized(App.pageHandler.Login))
	mux.HandleFunc("/logout", App.authHandler.Logout)
	mux.HandleFunc("/registration", App.pageHandler.Registration)
	mux.HandleFunc("/event", App.pageHandler.Event)
	mux.HandleFunc("/user/create", App.userHandler.Create)
	mux.HandleFunc("/auth/login", App.authHandler.Login)

	mux.HandleFunc("/ticket/create", App.ticketHandler.Create)

	mux.HandleFunc("/admin", App.authMiddleware.IsAdmin(App.pageHandler.Admin))

	mux.HandleFunc("/user", App.authMiddleware.IsAuthorized(App.pageHandler.User))

	mux.HandleFunc("/event/create", App.eventHandler.Create)
	fileServer := http.FileServer(http.Dir(cnf.SERVER.StaticDir))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	infoLog.Println("Starting server on :" + cnf.SERVER.Port)
	http.ListenAndServe(":"+cnf.SERVER.Port, session.Enable(mux))
	errorLog.Fatal(err)
}

package handlers

import (
	"finalProject/internal/repositories"
	"github.com/golangcollege/sessions"
	"log"
	"net/http"
	"strconv"
	"time"
)

type TicketHandler struct {
	Session    *sessions.Session
	TicketRepo *repositories.TicketRepository
}

func (th *TicketHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			log.Println(w, "ParseForm() err: %v", err)
			return
		}
		eventIdStr := r.FormValue("event_id")
		userIdStr := th.Session.GetString(r, "userId")

		currentTime := time.Now()
		unixTimestamp := currentTime.Unix()
		ticketNumber := strconv.Itoa(int(unixTimestamp))

		eventId, _ := strconv.Atoi(eventIdStr)
		userId, _ := strconv.Atoi(userIdStr)

		insertedId, err := th.TicketRepo.Insert(ticketNumber, eventId, userId)
		if err != nil {
			log.Println("Cannot create ticket")
			return
		}
		log.Println("Ticket with id ", insertedId, " Created")
		http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
	}

}

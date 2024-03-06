package handlers

import (
	"finalProject/internal/repositories"
	"fmt"
	"github.com/golangcollege/sessions"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type EventHandler struct {
	Session   *sessions.Session
	EventRepo *repositories.EventRepository
}

func (eh *EventHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(2 << 20)
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusInternalServerError)
		return
	}

	name := r.FormValue("name")
	date := r.FormValue("date")
	time_e := r.FormValue("time")
	venue := r.FormValue("venue")
	description := r.FormValue("description")
	price := r.FormValue("price")
	ticketQuantity := r.FormValue("ticket_quantity")
	file, _, err := r.FormFile("picture")

	combinedString := date + " " + time_e + ":00"
	layout := "2006-01-02 15:04:05"
	dateTime, err := time.Parse(layout, combinedString)
	if err != nil {
		fmt.Println("Error parsing date and time:", err)
		return
	}

	if err == nil {
		defer file.Close()

		currentTime := time.Now()

		unixTimestamp := currentTime.Unix()

		filename := strconv.Itoa(int(unixTimestamp)) + ".jpg"
		f, err := os.OpenFile("../../static/img/"+filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			http.Error(w, "Error saving file", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		defer f.Close()

		io.Copy(f, file)

		floatPrice, err := strconv.ParseFloat(price, 64)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		intQuantity, err := strconv.ParseInt(ticketQuantity, 10, 64)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		insertedId, err := eh.EventRepo.Insert(name, dateTime, description, venue, floatPrice, int(intQuantity), filename)
		if checkError(err, &w, "Article not created.", 404) {
			http.Error(w, "Error saving file", http.StatusInternalServerError)
			return
		}

		log.Println("New event with id ", insertedId, " created")

	}
}

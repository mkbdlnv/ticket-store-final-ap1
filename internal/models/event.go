package models

import "time"

type Event struct {
	ID             int
	Name           string
	DateTime       time.Time
	Description    string
	Venue          string
	TicketPrice    float64
	TicketQuantity int
	ImagePath      string
}

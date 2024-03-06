package handlers

import "finalProject/internal/models"

type TemplateData struct {
	Events     []models.Event
	Event      *models.Event
	FormatDate string
	FormatTime string
	Tickets    []*models.Ticket
	User       *models.User
	LoggedIn   bool
	IsAdmin    bool
}

type ErrorData struct {
	ErrorMsg string
}

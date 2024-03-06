package repositories

import (
	"database/sql"
	"finalProject/internal/models"
	"log"
	"time"
)

type EventRepository struct {
	DB *sql.DB
}

func (er *EventRepository) Insert(name string, datetime time.Time, description string, venue string, ticketPrice float64, ticketQuantity int, imagePath string) (int, error) {
	stmt := `INSERT INTO events (name, datetime, description, venue, ticket_price, image_path, ticket_quantity) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	var insertedID int
	err := er.DB.QueryRow(stmt, name, datetime, description, venue, ticketPrice, imagePath, ticketQuantity).Scan(&insertedID)
	if err != nil {
		return 0, err
	}
	return insertedID, nil
}

func (er *EventRepository) GetByID(id int) (*models.Event, error) {
	var event models.Event
	stmt := `SELECT * FROM events WHERE id = $1`
	err := er.DB.QueryRow(stmt, id).Scan(&event.ID, &event.Name, &event.DateTime, &event.Description, &event.Venue, &event.TicketPrice, &event.ImagePath, &event.TicketQuantity)
	if err != nil {
		log.Println("Error getting event by ID:", err)
		return nil, err
	}
	return &event, nil
}

func (er *EventRepository) GetAll() ([]models.Event, error) {
	var events []models.Event
	stmt := `SELECT * FROM events`
	rows, err := er.DB.Query(stmt)
	if err != nil {
		log.Println("Error getting all events:", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var event models.Event
		err := rows.Scan(&event.ID, &event.Name, &event.DateTime, &event.Description, &event.Venue, &event.TicketPrice, &event.ImagePath, &event.TicketQuantity)
		if err != nil {
			log.Println("Error scanning events:", err)
			return nil, err
		}
		events = append(events, event)
	}
	return events, nil
}

func (er *EventRepository) Delete(id int) error {
	stmt := `DELETE FROM events WHERE id = $1`
	_, err := er.DB.Exec(stmt, id)
	if err != nil {
		log.Println("Error deleting event:", err)
		return err
	}
	return nil
}

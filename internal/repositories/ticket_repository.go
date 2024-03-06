package repositories

import (
	"database/sql"
	"finalProject/internal/models"
)

type TicketRepository struct {
	DB *sql.DB
}

func (tr *TicketRepository) Insert(ticketNumber string, eventId, userId int) (int, error) {
	stmt := `INSERT INTO tickets(ticket_number, event_id, user_id) VALUES ($1, $2, $3) RETURNING id`
	var insertedID int
	err := tr.DB.QueryRow(stmt, ticketNumber, eventId, userId).Scan(&insertedID)
	if err != nil {
		return 0, err
	}
	return insertedID, nil
}

func (tr *TicketRepository) GetById(id int) (*models.Ticket, error) {
	stmt := `SELECT * FROM tickets WHERE id = $1`
	t := &models.Ticket{}
	err := tr.DB.QueryRow(stmt, id).Scan(&t.ID, &t.TicketNumber, &t.EventId, &t.UserId)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (tr *TicketRepository) GetAllByUserId(userId int) ([]*models.Ticket, error) {
	stmt := `SELECT * FROM tickets WHERE user_id=$1`
	rows, err := tr.DB.Query(stmt, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	tickets := []*models.Ticket{}
	for rows.Next() {
		t := &models.Ticket{}
		err := rows.Scan(&t.ID, &t.TicketNumber, &t.EventId, &t.UserId)
		if err != nil {
			return nil, err
		}
		tickets = append(tickets, t)
	}
	return tickets, nil
}

func (tr *TicketRepository) GetAll() ([]*models.Ticket, error) {
	stmt := `SELECT * FROM tickets`
	rows, err := tr.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	tickets := []*models.Ticket{}
	for rows.Next() {
		t := &models.Ticket{}
		err := rows.Scan(&t.ID, &t.TicketNumber, &t.EventId, &t.UserId)
		if err != nil {
			return nil, err
		}
		tickets = append(tickets, t)
	}
	return tickets, nil
}

func (tr *TicketRepository) Delete(id int) error {
	stmt := `DELETE FROM tickets WHERE id =$1`
	_, err := tr.DB.Exec(stmt, id)
	if err != nil {
		return err
	}
	return nil
}

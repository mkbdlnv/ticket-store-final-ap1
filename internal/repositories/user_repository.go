package repositories

import (
	"database/sql"
	"finalProject/internal/models"
	"finalProject/internal/utils"
)

type UserRepository struct {
	DB *sql.DB
}

func (ur *UserRepository) Insert(name, lastname, email, password string) (int, error) {
	stmt := `INSERT INTO users(name, lastname, email, password, role) VALUES ($1, $2, $3, $4, 'user') RETURNING id`
	hashedPass, err := utils.HashPassword(password)
	if err != nil {
		return 0, err
	}
	insertedId := 0
	err = ur.DB.QueryRow(stmt, name, lastname, email, hashedPass).Scan(&insertedId)
	if err != nil {
		return 0, err
	}
	return int(insertedId), nil
}

func (ur *UserRepository) GetByEmail(email string) (*models.User, error) {
	stmt := `SELECT * FROM users WHERE email=$1`
	u := &models.User{}
	err := ur.DB.QueryRow(stmt, email).Scan(&u.ID, &u.Name, &u.Lastname, &u.Email, &u.Password, &u.Role)
	if err == sql.ErrNoRows {
		return nil, err
	}
	return u, nil
}

func (ur *UserRepository) GetById(id int) (*models.User, error) {
	stmt := `SELECT * FROM users WHERE id=$1`
	u := &models.User{}
	err := ur.DB.QueryRow(stmt, id).Scan(&u.ID, &u.Name, &u.Lastname, &u.Email, &u.Password, &u.Role)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (ur *UserRepository) GetAll() ([]*models.User, error) {
	stmt := `SELECT * FROM users`
	rows, err := ur.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	users := []*models.User{}
	for rows.Next() {
		user := &models.User{}
		err = rows.Scan(&user.ID, &user.Name, &user.Lastname, &user.Email, &user.Password, &user.Role)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (ur *UserRepository) Delete(id int) error {
	stmt := `DELETE FROM users WHERE id = $1`
	_, err := ur.DB.Exec(stmt, id)
	if err != nil {
		return err
	}
	return nil
}

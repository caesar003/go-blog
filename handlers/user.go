package handlers

import (
	"database/sql"
	"errors"

	"github.com/go-sql-driver/mysql"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func CreateUser(db *sql.DB, name, email string) error {
	query := "INSERT INTO users (name, email) VALUES(?, ?)"
	_, err := db.Exec(query, name, email)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			// MySQL error code 1062 corresponds to a duplicate entry error
			return errors.New("email already in use")
		}
		return err
	}
	return nil
}

func GetAllUsers(db *sql.DB) ([]User, error) {
	query := "SELECT * FROM users"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func GetUser(db *sql.DB, id int) (*User, error) {
	query := "SELECT * FROM users where id = ?"
	row := db.QueryRow(query, id)

	user := &User{}
	err := row.Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func UpdateUser(db *sql.DB, id int, name, email string) error {
	query := "UPDATE users set name = ?, email = ? WHERE id = ?"
	_, err := db.Exec(query, name, email, id)
	if err != nil {
		return err
	}
	return nil
}

func DeleteUser(db *sql.DB, id int) error {
	query := "DELETE FROM users WHERE id = ?"
	_, err := db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

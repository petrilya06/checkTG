package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type User struct {
	PhoneNumber string
	TgID        int
	HasText     bool
	HasPic      bool
}

type Database struct {
	conn *sql.DB
}

// NewDatabase создает новое соединение с базой данных
func NewDatabase(connStr string) (*Database, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("can`t connect to db: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error in connection: %w", err)
	}

	return &Database{conn: db}, nil
}

// Close закрывает соединение с базой данных
func (d *Database) Close() {
	d.conn.Close()
}

func (d *Database) AddUser(phoneNumber string, tgID int, hasText bool, hasPic bool) error {
	insertSQL := `INSERT INTO users (phone_number, tg_id, has_text, has_pic) VALUES ($1, $2, $3, $4)`
	_, err := d.conn.Exec(insertSQL, phoneNumber, tgID, hasText, hasPic)
	if err != nil {
		return fmt.Errorf("error in insert: %w", err)
	}

	return nil
}

func (d *Database) CreateTable() error {
	createTableSQL := `
    CREATE TABLE IF NOT EXISTS users (
        phone_number VARCHAR(15) NOT NULL,
        tg_id INT PRIMARY KEY,
        has_text BOOLEAN,
        has_pic BOOLEAN
    );`

	_, err := d.conn.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("error in create: %w", err)
	}
	return nil
}

func (d *Database) GetAllUsers() ([]User, error) {
	rows, err := d.conn.Query("SELECT phone_number, tg_id, has_text, has_pic FROM users")

	if err != nil {
		return nil, fmt.Errorf("error in data: %w", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.PhoneNumber, &user.TgID, &user.HasText, &user.HasPic); err != nil {
			return nil, fmt.Errorf("error in scan: %w", err)
		}
		users = append(users, user)
	}

	// Проверка на ошибки после завершения выборки
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error in processing: %w", err)
	}

	return users, nil
}

package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type User struct {
	TgID          int64
	SelectPic     int
	HasText       bool
	HasPic        bool
	CountCheck    int
	Status        string
	LastPhotoID   int
	LastMessageID int
}

type Database struct {
	conn *sql.DB
}

func NewDatabase() (*Database, error) {
	connStr := "user=postgres password=qwerty123456 host=localhost sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("can't connect to db: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error in connection: %w", err)
	}

	return &Database{conn: db}, nil
}

func (d *Database) CreateTable() error {
	createTableSQL := `
    CREATE TABLE IF NOT EXISTS users (
        tg_id BIGINT PRIMARY KEY,
		select_pic INT,
        has_text BOOLEAN,
        has_pic BOOLEAN,
        count_check INT,
        status VARCHAR(30),
        last_photo_id INT,
        last_message_id INT
    );`

	_, err := d.conn.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("error in create: %w", err)
	}
	return nil
}

func (d *Database) AddUser(user User) error {
	deleteSQL := `DELETE FROM users WHERE tg_id = $1;`
	_, err := d.conn.Exec(deleteSQL, user.TgID)
	if err != nil {
		return fmt.Errorf("error in delete: %w", err)
	}

	insertSQL := `
		INSERT INTO users (tg_id, select_pic, has_text, has_pic, count_check, status, last_photo_id, last_message_id) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8);`
	_, err = d.conn.Exec(insertSQL, user.TgID, user.SelectPic, user.HasText, user.HasPic, user.CountCheck,
		user.Status, user.LastPhotoID, user.LastMessageID)
	if err != nil {
		return fmt.Errorf("error in insert: %w", err)
	}

	return nil
}

func (d *Database) UpdateUser(user User) error {
	updateSQL := `
		UPDATE users 
		SET select_pic = $1, has_text = $2, has_pic = $3, count_check = $4, 
		    status = $5, last_photo_id = $6, last_message_id = $7 
		WHERE tg_id = $8`
	_, err := d.conn.Exec(updateSQL, user.SelectPic, user.HasText, user.HasPic,
		user.CountCheck, user.Status, user.LastPhotoID, user.LastMessageID, user.TgID)
	if err != nil {
		return fmt.Errorf("error in update: %w", err)
	}

	return nil
}

func (d *Database) GetUserByID(tgID int64) (*User, error) {
	query := `
		SELECT tg_id, select_pic, has_text, has_pic, count_check, 
		       status, last_photo_id, last_message_id 
		FROM users WHERE tg_id = $1`

	var user User
	err := d.conn.QueryRow(query, tgID).Scan(&user.TgID, &user.SelectPic, &user.HasText, &user.HasPic,
		&user.CountCheck, &user.Status, &user.LastPhotoID, &user.LastMessageID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user with tg_id %d not found", tgID)
		}
		return nil, fmt.Errorf("error in query: %w", err)
	}

	return &user, nil
}

func (d *Database) Close() {
	if err := d.conn.Close(); err != nil {
		fmt.Printf("error closing database: %v\n", err)
	}
}

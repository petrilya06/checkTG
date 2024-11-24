package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type User struct {
	TgID      int64
	SelectPic int
	HasText   bool
	HasPic    bool
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
        has_pic BOOLEAN
    );`

	_, err := d.conn.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("error in create: %w", err)
	}
	return nil
}

func (d *Database) AddUser(user User) error {
	insertSQL := `
		INSERT INTO users (tg_id, select_pic, has_text, has_pic) 
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (tg_id) DO UPDATE 
		SET select_pic = EXCLUDED.select_pic, 
			has_text = EXCLUDED.has_text, 
			has_pic = EXCLUDED.has_pic;`

	_, err := d.conn.Exec(insertSQL, user.TgID, user.SelectPic, user.HasText, user.HasPic)
	if err != nil {
		return fmt.Errorf("error in insert: %w", err)
	}

	return nil
}

func (d *Database) UpdateUser(user User) error {
	updateSQL := `
		UPDATE users 
		SET select_pic = $1, has_text = $2, has_pic = $3 
		WHERE tg_id = $4`
	_, err := d.conn.Exec(updateSQL, user.SelectPic, user.HasText, user.HasPic, user.TgID)
	if err != nil {
		return fmt.Errorf("error in update: %w", err)
	}

	return nil
}

func (d *Database) GetAllUsers() ([]User, error) {
	rows, err := d.conn.Query("SELECT tg_id, select_pic, has_text, has_pic FROM users")
	if err != nil {
		return nil, fmt.Errorf("error in data: %w", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.TgID, &user.SelectPic, &user.HasText, &user.HasPic); err != nil {
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

func (d *Database) Close() {
	if err := d.conn.Close(); err != nil {
		fmt.Printf("error closing database: %v\n", err)
	}
}

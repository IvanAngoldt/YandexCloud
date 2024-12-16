package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

// Connect подключается к базе данных и возвращает соединение
func Connect() (*sql.DB, error) {
	host := os.Getenv("POSTGRES_HOST")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", host, user, password, dbname)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

// User представляет данные пользователя
type User struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
}

// GetUserByEmail выполняет запрос к базе данных для получения пользователя по email
func GetUserByEmail(db *sql.DB, email string) (*User, error) {
	var user User
	err := db.QueryRow("SELECT id, username, email, password_hash FROM users WHERE email = $1", email).
		Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash)

	if err == sql.ErrNoRows {
		return nil, nil // Пользователь не найден
	} else if err != nil {
		return nil, err // Ошибка базы данных
	}
	return &user, nil // Пользователь найден
}

// SaveUser сохраняет нового пользователя в базе данных
func SaveUser(db *sql.DB, username, email, passwordHash string) error {
	_, err := db.Exec(`
		INSERT INTO users (username, email, password_hash)
		VALUES ($1, $2, $3)
	`, username, email, passwordHash)
	return err
}

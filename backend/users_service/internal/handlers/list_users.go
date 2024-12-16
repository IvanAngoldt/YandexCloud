package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

func ListUsers(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, username, email FROM users ORDER BY id ASC")
		if err != nil {
			log.Println("ListUsers: Query error:", err)
			http.Error(w, "Failed to list users", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var users []map[string]interface{}
		for rows.Next() {
			var (
				id       int
				username string
				email    string
			)
			if err := rows.Scan(&id, &username, &email); err != nil {
				log.Println("ListUsers: Scan error:", err)
				http.Error(w, "Failed to parse users", http.StatusInternalServerError)
				return
			}
			users = append(users, map[string]interface{}{
				"id":       id,
				"username": username,
				"email":    email,
			})
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	}
}

package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

func GetUserByUsername(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := r.URL.Query().Get("username")
		if username == "" {
			http.Error(w, "Username is required", http.StatusBadRequest)
			return
		}

		var user struct {
			ID           int    `json:"id"`
			Username     string `json:"username"`
			Email        string `json:"email"`
			PasswordHash string `json:"-"`
		}

		err := db.QueryRow("SELECT id, username, email, password_hash FROM users WHERE username = $1", username).
			Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash)
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		} else if err != nil {
			http.Error(w, "Failed to fetch user", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	}
}

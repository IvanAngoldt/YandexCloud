package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

type UpdateUserRequest struct {
	Username *string `json:"username,omitempty"`
	Email    *string `json:"email,omitempty"`
}

func UpdateUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userIDStr := vars["id"]
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		var req UpdateUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// Собираем динамический запрос
		var setParts []string
		var args []interface{}
		argIndex := 1

		if req.Username != nil && strings.TrimSpace(*req.Username) != "" {
			setParts = append(setParts, "username = $"+strconv.Itoa(argIndex))
			args = append(args, strings.TrimSpace(*req.Username))
			argIndex++
		}

		if req.Email != nil && strings.TrimSpace(*req.Email) != "" {
			setParts = append(setParts, "email = $"+strconv.Itoa(argIndex))
			args = append(args, strings.TrimSpace(*req.Email))
			argIndex++
		}

		if len(setParts) == 0 {
			http.Error(w, "No fields to update", http.StatusBadRequest)
			return
		}

		query := "UPDATE users SET " + strings.Join(setParts, ", ") + " WHERE id = $" + strconv.Itoa(argIndex)
		args = append(args, userID)

		_, err = db.Exec(query, args...)
		if err != nil {
			http.Error(w, "Failed to update user", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "User updated successfully"})
	}
}

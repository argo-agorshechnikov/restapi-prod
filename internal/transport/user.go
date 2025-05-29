package transport

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/argo-agorshechnikov/restapi-prod/internal/models"
)

func usersHandler(db *sql.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			// Get users list
			handleGetUsers(w, r, db)

		case http.MethodPost:
			// Create new user
			handleCreateUser(w, r, db)

		//case http.MethodPut:

		//case http.MethodDelete:

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}

}

func handleCreateUser(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	var user models.User

	// Stuct fulling
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid user data: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// New user insert in db
	err := db.QueryRow(
		"INSERT INTO users (name, email) VALUES($1, $2) RETURNING id",
		user.Name, user.Email,
	).Scan(&user.Id)
	if err != nil {
		http.Error(w, "DB insert error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("User struct: ", user)

	// json format
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	// user to json
	if err := json.NewEncoder(w).Encode(user); err != nil {

		http.Error(w, "Error to convert user to json", http.StatusBadRequest)
		return
	}

}

func handleGetUsers(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	// get users from db
	rows, err := db.Query("SELECT id, name, email FROM users")
	if err != nil {
		http.Error(w, "DB error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []models.User

	// Each rows iteration
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.Id, &user.Name, &user.Email); err != nil {
			http.Error(w, "Scan error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		users = append(users, user)
	}

	log.Printf("Users struct constraint: ", users)

	// Check iteration err
	if err = rows.Err(); err != nil {

		http.Error(w, "Rows error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(users); err != nil {

		http.Error(w, "Error send users list", http.StatusBadRequest)
		return
	}
}

package transport

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/argo-agorshechnikov/restapi-prod/internal/models"
)

func HandleCreateUser(w http.ResponseWriter, r *http.Request, db *sql.DB) {

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

func HandleGetUsers(w http.ResponseWriter, r *http.Request, db *sql.DB) {

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

func HandleUpdateUser(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	// Get value "id" from URL user/1 -> idStd == 1
	idStr := r.PathValue("id")
	if idStr == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	// idStr(string) -> int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Decoder reads r.Body, decode parsing json to user struct
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close() // allows free up resources

	// Exec return result sql-query and err
	result, err := db.Exec(
		"UPDATE users SET name = $1, email = $2 WHERE id = $3",
		user.Name, user.Email, id,
	)
	if err != nil {
		http.Error(w, "DB update error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// get count string which been change, it's allows check change user by id
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Error check update result: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		http.Error(w, "User not found", http.StatusNotFound)
	}

	// Allows returned full user object
	user.Id = id

	// Set header response json
	w.Header().Set("Content-Type", "application/json")
	// Convert User to json
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "Error user update", http.StatusBadRequest)
		return
	}

}

func HandleDeleteUser(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	// Get string value id from URL
	idStr := r.PathValue("id")
	if idStr == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	// Convert string to int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Result sql-query
	relust, err := db.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		http.Error(w, "DB delete error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Check rows been changed(Deleted)
	rowsAffected, err := relust.RowsAffected()
	if err != nil {
		http.Error(w, "Error check delete result: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// 204 status
	w.WriteHeader(http.StatusNoContent)
}

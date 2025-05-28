package transport

import (
	"encoding/json"
	"net/http"

	"github.com/argo-agorshechnikov/restapi-prod/internal/models"
)

func usersHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		// Get users list
		handleGetUsers(w, r)

	case http.MethodPost:
		// Create new user
		handleCreateUser(w, r)

	//case http.MethodPut:

	//case http.MethodDelete:

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

}

func handleCreateUser(w http.ResponseWriter, r *http.Request) {

	var user models.User

	// Stuct fulling
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid user data: "+err.Error(), http.StatusBadRequest)
		return
	}

	defer r.Body.Close()
	// json format
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	// user to json
	if err := json.NewEncoder(w).Encode(user); err != nil {

		http.Error(w, "Error to convert user to json", http.StatusBadRequest)
		return
	}

}

func handleGetUsers(w http.ResponseWriter, r *http.Request) {

	users := []models.User{}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(users); err != nil {

		http.Error(w, "Error send users list", http.StatusBadRequest)
		return
	}
}

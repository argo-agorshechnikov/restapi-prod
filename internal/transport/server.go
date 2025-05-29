package transport

import (
	"database/sql"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/argo-agorshechnikov/restapi-prod/internal/database"
)

type Server struct {
	port   string
	logger *slog.Logger
	db     *sql.DB
}

func NewServer(port string) *Server {
	handler := slog.NewJSONHandler(os.Stdout, nil)
	logger := slog.New(handler)
	return &Server{
		port:   port,
		logger: logger,
	}
}

func (s *Server) StartServer() error {

	// connect to db
	db, err := database.ConnectionDB("argo", "argo", "restapi_db", "localhost", 5432)
	if err != nil {
		log.Fatalf("Failed to connect to db: %v", err)
	}
	s.db = db
	defer db.Close()
	s.logger.Info("Successfully connection to restapi_db")

	mux := http.NewServeMux()

	mux.HandleFunc("GET /users", func(w http.ResponseWriter, r *http.Request) {
		HandleGetUsers(w, r, s.db)
	})
	mux.HandleFunc("POST /users", func(w http.ResponseWriter, r *http.Request) {
		HandleCreateUser(w, r, s.db)
	})
	mux.HandleFunc("PUT /users/{id}", func(w http.ResponseWriter, r *http.Request) {
		HandleUpdateUser(w, r, s.db)
	})

	mux.HandleFunc("/", homeHandler)

	s.logger.Info("Server start on port: " + s.port)
	return http.ListenAndServe(s.port, mux)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("Home"))
}

package transport

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/argo-agorshechnikov/restapi-prod/internal/database"
)

type Server struct {
	port   string
	logger *slog.Logger
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

	defer db.Close()
	s.logger.Info("Successfully connection to restapi_db")

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/users", usersHandler(db))

	s.logger.Info("Server start")
	return http.ListenAndServe(s.port, nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("Home"))
}

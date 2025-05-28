package transport

import (
	"log/slog"
	"net/http"
	"os"
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
	http.HandleFunc("/", helloHandler)

	s.logger.Info("Server start")
	return http.ListenAndServe(":"+s.port, nil)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("Hello server!"))
}

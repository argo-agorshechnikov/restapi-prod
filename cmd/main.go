package main

import (
	"log"

	"github.com/argo-agorshechnikov/restapi-prod/internal/transport"
)

func main() {
	s := transport.NewServer(":8080")
	if err := s.StartServer(); err != nil {
		log.Fatal(err)
	}
}

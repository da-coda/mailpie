package main

import (
	"log"
	"mailmon-go/pkg/backend"
	"mailmon-go/pkg/server"
)

func main() {
	be := backend.New()
	s := server.NewDefaultServer(be)

	log.Println("Starting server at", s.Addr)
	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
package main

import (
	"log"
	"os"

	server "github.com/utkarsh-pro/RapidoDB/server"
)

const (
	// Default port for the TCP server
	defaultPort = "2310"
)

func main() {
	// Get the port from the environment
	PORT := os.Getenv("PORT")
	// If the port isn't provided then fallback to
	// the default port
	if PORT == "" {
		PORT = defaultPort
	}

	// Instantiate the TCP server
	s := server.New(log.New(os.Stdout, "[SERVER]: ", log.LstdFlags))
	// Setup and spin up the TCP server
	s.Setup(PORT)
}

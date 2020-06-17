package main

import (
	"log"
	"os"

	server "github.com/utkarsh-pro/RapidoDB/server"
)

const (
	// Default port for the TCP server
	defaultPort = "2310"
	// Default password for the server
	defaultPass = "pass"
	// Default username for the server
	defaultUser = "admin"
)

func main() {
	// Get the port from the environment
	PORT := os.Getenv("PORT")
	PASS := os.Getenv("PASS")
	USER := os.Getenv("USER")

	// If the port isn't provided then fallback to
	// the default port
	if PORT == "" {
		PORT = defaultPort
	}

	// If the pass isn't provided then fallback to
	// the default pass
	if PASS == "" {
		PASS = defaultPass
	}

	// If the user isn't provided then fallback to
	// the default user
	if USER == "" {
		USER = defaultUser
	}

	// Instantiate the TCP server
	s := server.New(log.New(os.Stdout, "[SERVER]: ", log.LstdFlags), PORT, USER, PASS)
	// Setup and spin up the TCP server
	s.Setup()
}

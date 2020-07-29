package main

import (
	"log"
	"os"

	db "github.com/utkarsh-pro/RapidoDB/DB"
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
	PORT := getEnv("RAPIDO_PORT", defaultPort)
	// PASS := getEnv("RAPIDO_PASS", defaultPass)
	// USER := getEnv("RAPIDO_USER", defaultUser)

	// Instantiate the TCP server
	// s := server.New(log.New(os.Stdout, "[SERVER]: ", log.LstdFlags), PORT, USER, PASS)
	// Setup and spin up the TCP server
	// s.Setup()

	database := db.New(log.New(os.Stdout, "[RAPIDO DB]: ", log.LstdFlags), PORT)

	database.Run()
}

// getEnv is a thin wrapper over os.GetEnv. It replaces read
// value with a fallback if the env var is an empty string
func getEnv(env string, fallback string) string {
	param := os.Getenv(env)
	if param == "" {
		param = fallback
	}
	return param
}

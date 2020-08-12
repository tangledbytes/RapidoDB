package main

import (
	"fmt"
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
	PASS := getEnv("RAPIDO_PASS", defaultPass)
	USER := getEnv("RAPIDO_USER", defaultUser)
	BACKUP := getEnv("HOME", "")

	// Print the RapidoDB logo
	fmt.Println(db.RapidoMSG)

	database := db.New(log.New(os.Stdout, "[RAPIDO DB]: ", log.LstdFlags), PORT, USER, PASS, BACKUP)

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

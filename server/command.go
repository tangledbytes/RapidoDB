package server

import (
	"encoding/json"
	"errors"
	"strings"
)

type commandID int

const (
	cmdSet commandID = iota
	cmdGet
	cmdDelete
	cmdQuit
	cmdAuth
)

type command struct {
	id     commandID
	client *client
	args   []string
}

// Parse the commands sent by the client
func parse(args []string, cl *client) (command, error) {
	cmd := strings.ToLower(args[0])

	switch cmd {
	case "set":
		return command{
			id:     cmdSet,
			args:   args,
			client: cl,
		}, nil
	case "get":
		return command{
			id:     cmdGet,
			args:   args,
			client: cl,
		}, nil
	case "delete":
		return command{
			id:     cmdDelete,
			args:   args,
			client: cl,
		}, nil
	case "quit":
		return command{
			id:     cmdQuit,
			args:   args,
			client: cl,
		}, nil
	case "auth":
		return command{
			id:     cmdAuth,
			args:   args,
			client: cl,
		}, nil
	default:
		return command{}, errors.New("Unknown command: " + cmd)
	}
}

// deserialise deserialises the passed json into a map.
// It also ensures that the passed json has an "input" field
func deserialise(input string) (map[string]interface{}, error) {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(input), &data); err != nil {
		return data, err
	}

	if _, ok := data["input"]; !ok {
		return data, errors.New("No input field in the passed JSON")
	}
	return data, nil
}

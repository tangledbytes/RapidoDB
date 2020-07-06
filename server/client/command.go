package client

import (
	"encoding/json"
	"errors"
	"strings"
)

type commandID int

const (
	// CmdSet represents the SET command
	CmdSet commandID = iota
	// CmdGet represents the GET command
	CmdGet
	// CmdDelete represents the DELETE command
	CmdDelete
	// CmdQuit represents the QUIT command
	CmdQuit
	// CmdAuth represents the AUTH command
	CmdAuth
)

// Command represents the command sent by the client
type Command struct {
	id     commandID
	client *Client
	args   []string
}

// Parse parses the commands sent by the client
func Parse(args []string, cl *Client) (Command, error) {
	cmd := strings.ToLower(args[0])

	switch cmd {
	case "set":
		return Command{
			id:     CmdSet,
			args:   args,
			client: cl,
		}, nil
	case "get":
		return Command{
			id:     CmdGet,
			args:   args,
			client: cl,
		}, nil
	case "delete":
		return Command{
			id:     CmdDelete,
			args:   args,
			client: cl,
		}, nil
	case "quit":
		return Command{
			id:     CmdQuit,
			args:   args,
			client: cl,
		}, nil
	case "auth":
		return Command{
			id:     CmdAuth,
			args:   args,
			client: cl,
		}, nil
	default:
		return Command{}, errors.New("Unknown command: " + cmd)
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

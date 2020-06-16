package server

import "errors"

type commandID int

const (
	cmdSet commandID = iota
	cmdGet
	cmdDelete
	cmdQuit
)

type command struct {
	id     commandID
	client *client
	args   []string
}

// Parse the commands sent by the client
func parse(args []string, cl *client) (command, error) {
	cmd := args[0]

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
	default:
		return command{}, errors.New("Unknown command: " + cmd)
	}
}

package server

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/utkarsh-pro/RapidoDB/db"
)

type client struct {
	conn            net.Conn
	commands        chan<- command
	log             *log.Logger
	isAuthenticated bool
	db              *db.DB
}

// newClient returns a new client
func newClient(conn net.Conn, cmd chan<- command, log *log.Logger, isAuthenticated bool, db *db.DB) *client {
	return &client{
		conn:            conn,
		commands:        cmd,
		log:             log,
		isAuthenticated: isAuthenticated,
		db:              db,
	}
}

// initRead reads the input of the TCP clients
func (c *client) initRead(auth Auth) {
	for {
		// Read data from TCP client and parse it
		data, err := bufio.NewReader(c.conn).ReadString('\n')

		// Check for errors
		if err != nil {
			// If error is io.EOF then it indicates that the client has
			// disconnected and hence closing the connection here
			if err == io.EOF {
				c.log.Printf("Client %s disconnected", c.conn.RemoteAddr().String())
				c.conn.Close()
				return
			}

			// Log the error
			c.log.Printf("Error from client %s: %v", c.conn.RemoteAddr().String(), err)
			return
		}

		// Trim the data
		data = strings.Trim(data, "\n")

		// Parse the commands
		cmd, err := parse(strings.Split(data, " "), c)

		if err != nil {
			// Send error to the client if the command is not known
			c.err(err)
			continue
		}

		// Send the command to the server
		// c.commands <- cmd
		c.handleCommand(auth, cmd)
	}
}

// msg sends a message to the client
func (c *client) msg(msg string) {
	c.conn.Write([]byte(msg + "\n"))
}

// err sends an error message to the client
func (c *client) err(err error) {
	c.conn.Write([]byte("ERR: " + err.Error() + "\n"))
}

// handleCommand checks for authentication of the client
//
// It exclusively handles "AUTH" commands and passes on other commands
// to "authorizedCommandHandler"
func (c *client) handleCommand(auth Auth, cmd command) {
	// Check if the AUTH type command is sent
	if cmd.id == cmdAuth && !c.isAuthenticated {

		// Handle authentication
		isAuthenticated, err := auth.HandleAuth(cmd.args)

		if err != nil {
			cmd.client.err(err)
			return
		}

		// Set the authentication
		cmd.client.isAuthenticated = isAuthenticated
		// Respond with success
		cmd.client.msg("Success")
		return
	}

	// The commands will only be handled if the client is authenticated
	if cmd.client.isAuthenticated {
		c.log.Printf("Received: %v from %v", cmd.id, cmd.client.conn.RemoteAddr().String())
		c.authorizedCommandHandler(cmd)
	} else {
		cmd.client.err(errors.New("Not authorized"))
	}
}

// authorizedCommandHandler handles authorized commands
// hence this method should only be called when the authorization
// of the client is guaranteed
func (c *client) authorizedCommandHandler(cmd command) {
	switch cmd.id {
	case cmdSet:
		c.set(cmd)
	case cmdGet:
		c.get(cmd)
	default:
		c.err(errors.New("Invalid command"))
	}
}

func (c *client) set(cmd command) {
	// set command looks like
	// SET <key> <value> [expiry in ms]
	argLen := len(cmd.args)

	if argLen < 3 {
		c.err(errors.New("Invalid SET command\n\tSYNTAX: SET <key> <value> [expiry in milliseconds]"))
		return
	}

	key := cmd.args[1]
	value := cmd.args[2]
	expire := db.NeverExpire

	// Check if a expiry is given
	if argLen == 4 {
		e, err := strconv.Atoi(cmd.args[3])

		if err != nil {
			c.log.Println("ERROR: Invalid expiry provided")
			c.err(err)
			return
		}

		item := db.NewItem(value, time.Duration(e)*time.Millisecond)
		c.db.Set(key, item)
		return
	}

	// Create item for insertion
	item := db.NewItem(value, time.Duration(expire))
	c.db.Set(key, item)
	return
}

func (c *client) get(cmd command) {
	// get command looks like
	// GET <key>
	argLen := len(cmd.args)

	if argLen < 2 {
		c.err(errors.New("Invalid GET command\n\tSYNTAX: GET <key>"))
		return
	}

	c.msg(fmt.Sprintf("Value: %v", c.db.Get(cmd.args[1])))
}

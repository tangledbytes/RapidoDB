package client

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

// Client represents an active TCP client communicating
// with the server
type Client struct {
	conn            net.Conn
	commands        chan<- Command
	log             *log.Logger
	isAuthenticated bool
	db              *db.DB
}

// NewClient returns a new client
func NewClient(conn net.Conn, cmd chan<- Command, log *log.Logger, isAuthenticated bool, db *db.DB) *Client {
	return &Client{
		conn:            conn,
		commands:        cmd,
		log:             log,
		isAuthenticated: isAuthenticated,
		db:              db,
	}
}

// InitRead reads the input of the TCP clients
func (c *Client) InitRead(auth Auth) {
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
		cmd, err := Parse(strings.Split(data, " "), c)

		if err != nil {
			// Send error to the client if the command is not known
			c.Err(err)
			continue
		}

		// Send the command to the server
		// c.commands <- cmd
		c.handleCommand(auth, cmd)
	}
}

// Msg sends a message to the client
func (c *Client) Msg(msg string) {
	c.conn.Write([]byte(msg + "\n"))
}

// Err sends an error message to the client
func (c *Client) Err(err error) {
	c.conn.Write([]byte("ERR: " + err.Error() + "\n"))
}

// handleCommand checks for authentication of the client
//
// It exclusively handles "AUTH" commands and passes on other commands
// to "authorizedCommandHandler"
func (c *Client) handleCommand(auth Auth, cmd Command) {
	// Check if the AUTH type command is sent
	if cmd.id == CmdAuth && !c.isAuthenticated {

		// Handle authentication
		isAuthenticated, err := auth.HandleAuth(cmd.args)

		if err != nil {
			cmd.client.Err(err)
			return
		}

		// Set the authentication
		cmd.client.isAuthenticated = isAuthenticated
		// Respond with success
		cmd.client.Msg("Success")
		return
	}

	// The commands will only be handled if the client is authenticated
	if cmd.client.isAuthenticated {
		c.log.Printf("Received: %v from %v", cmd.id, cmd.client.conn.RemoteAddr().String())
		c.authorizedCommandHandler(cmd)
	} else {
		cmd.client.Err(errors.New("Not authorized"))
	}
}

// authorizedCommandHandler handles authorized commands
// hence this method should only be called when the authorization
// of the client is guaranteed
func (c *Client) authorizedCommandHandler(cmd Command) {
	switch cmd.id {
	case CmdSet:
		c.set(cmd)
	case CmdGet:
		c.get(cmd)
	default:
		c.Err(errors.New("Invalid command"))
	}
}

func (c *Client) set(cmd Command) {
	// set command looks like
	// SET <key> <value> [expiry in ms]
	argLen := len(cmd.args)

	if argLen < 3 {
		c.Err(errors.New("Invalid SET command\n\tSYNTAX: SET <key> <value> [expiry in milliseconds]"))
		return
	}

	key := cmd.args[1]
	value := cmd.args[2]
	expire := time.Duration(db.NeverExpire)

	d, err := deserialise(value)

	if err != nil {
		c.log.Printf("ERROR: %s", err)
		c.Err(errors.New("Invalid input"))
		return
	}

	// Check if a expiry is given
	if argLen == 4 {
		e, err := strconv.Atoi(cmd.args[3])

		if err != nil {
			c.log.Println("ERROR: Invalid expiry provided")
			c.Err(err)
			return
		}

		expire = time.Duration(e) * time.Millisecond
		return
	}

	// Create item for insertion
	item := db.NewItem(d["input"], expire)
	c.db.Set(key, item)
	return
}

func (c *Client) get(cmd Command) {
	// get command looks like
	// GET <key>
	argLen := len(cmd.args)

	if argLen < 2 {
		c.Err(errors.New("Invalid GET command\n\tSYNTAX: GET <key>"))
		return
	}

	c.Msg(fmt.Sprintf("Value: %v", c.db.Get(cmd.args[1])))
}

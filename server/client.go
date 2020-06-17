package server

import (
	"bufio"
	"io"
	"log"
	"net"
	"strings"
)

type client struct {
	conn            net.Conn
	commands        chan<- command
	log             *log.Logger
	isAuthenticated bool
}

// read reads the input of the TCP clients
func (c *client) read() {
	for {
		// Read data from TCP client and parse it
		data, err := bufio.NewReader(c.conn).ReadString('\n')

		// Check for errors
		if err != nil {
			// If error is io.EOF then it indicates that the client has
			// exited and hence closing the connection here
			if err == io.EOF {
				c.log.Printf("Client %s exited", c.conn.RemoteAddr().String())
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
		c.commands <- cmd
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

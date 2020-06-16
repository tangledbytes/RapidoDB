package server

import (
	"bufio"
	"io"
	"log"
	"net"
	"strings"
)

type client struct {
	conn     net.Conn
	commands chan<- command
	log      *log.Logger
}

// read reads the input of the TCP clients
func (c *client) read() {
	for {
		// Read data from TCP client and parse it
		data, err := bufio.NewReader(c.conn).ReadString(';')

		// Log the received data
		c.log.Printf("Received from %s: %v", c.conn.RemoteAddr().String(), data)

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

		// Parse the command
		cmd, err := parse(strings.Split(data[:len(data)-1], " "), c)

		if err != nil {
			// Send error to the client of the message is not known
			c.err(err)
			continue
		}

		// Send the parsed command to the channel
		c.commands <- cmd
	}
}

// msg sends a message to the client
func (c *client) msg(msg string) {
	c.conn.Write([]byte(msg))
}

// err sends an error message to the client
func (c *client) err(err error) {
	c.conn.Write([]byte("ERR: " + err.Error() + "\n"))
}

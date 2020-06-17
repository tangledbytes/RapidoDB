package server

import (
	"errors"
	"log"
	"net"
)

// Server is a struct
type Server struct {
	commands chan command
	log      *log.Logger
	PORT     string
	user     string
	pass     string
}

// New returns a an instance of server
func New(l *log.Logger, PORT string, user string, pass string) *Server {
	return &Server{
		commands: make(chan command),
		log:      l,
		PORT:     PORT,
		user:     user,
		pass:     pass,
	}
}

// Setup setups a TCP server
func (s *Server) Setup() {
	// Setup the TCP server
	listener, err := net.Listen("tcp", ":"+s.PORT)
	if err != nil {
		s.log.Fatalln(err)
		return
	}

	defer listener.Close()
	s.log.Println("Started server on PORT", s.PORT)

	// Start listening to the commands sent by the clients
	go s.listenCommands()

	// An infinite loop to listen for any number of TCP clients
	for {
		// Accept WAITS for and returns the next connection
		// to the listener. This is a blocking call.
		conn, err := listener.Accept()
		if err != nil {
			s.log.Println("Unable to accept connection: ", err.Error())
			continue
		}

		// Handle the client
		go s.clientHandler(conn)
	}
}

// clientHandler handles the client connecting to the server
func (s *Server) clientHandler(c net.Conn) {
	// Print the address of the client
	s.log.Println("Connected: ", c.RemoteAddr().String())

	// Create a client
	cl := &client{
		conn:            c,
		commands:        s.commands,
		log:             s.log,
		isAuthenticated: false,
	}

	// Inform the client about the connection
	cl.msg("Successfully connected to RapidoDB. Please run AUTH <username> <password> to access the DB")

	// Initialise the reader for the client
	cl.read()
}

// listenCommands listens for the commands passed on to the
// client channels and passes it onto the commandHandler
func (s *Server) listenCommands() {
	for cmd := range s.commands {
		// Check if the AUTH type command is sent
		if cmd.id == cmdAuth && !cmd.client.isAuthenticated {

			// Check if the command length is 3
			if len(cmd.args) != 3 {
				cmd.client.err(errors.New(("Invalid AUTH command")))
				return
			}

			// Get the user and password
			user := cmd.args[1]
			pass := cmd.args[2]

			// Check if the creds matches
			if user == s.user && pass == s.pass {
				// Set the authentication to true
				cmd.client.isAuthenticated = true
				// Respond with success
				cmd.client.msg("Success")
			} else {
				cmd.client.err(errors.New("Invalid credentials"))
			}

			continue
		}

		// The commands will only be handled if the client is authenticated
		if cmd.client.isAuthenticated {
			s.log.Printf("Received: %v from %v", cmd.id, cmd.client.conn.RemoteAddr().String())
		} else {
			cmd.client.err(errors.New(" Not authorized"))
		}
	}
}

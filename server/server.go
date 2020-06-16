package server

import (
	"log"
	"net"
)

// Server is a struct
type Server struct {
	commands chan command
	log      *log.Logger
}

// New returns a an instance of server
func New(l *log.Logger) *Server {
	return &Server{
		commands: make(chan command),
		log:      l,
	}
}

// Setup setups a TCP server
func (s *Server) Setup(PORT string) {
	// Setup the TCP server
	listener, err := net.Listen("tcp", ":"+PORT)
	if err != nil {
		s.log.Fatalln(err)
		return
	}

	defer listener.Close()
	s.log.Println("Started server on PORT", PORT)

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
		conn:     c,
		commands: s.commands,
		log:      s.log,
	}

	// Initialise the reader for the client
	cl.read()
}

// listenCommands listens for the commands passed on to the
// client channels
func (s *Server) listenCommands() {
	for cmd := range s.commands {
		s.log.Printf("Received: %v from %v", cmd.id, cmd.client.conn.RemoteAddr().String())
	}
}

package transport

import (
	"io"
	"log"
	"net"
)

// Driver interface demands an object which
// has an Operate method, it takes in string as input
// and returns nothing
type Driver interface {
	Operate(cmd string, w io.Writer)
}

// DriverFactory interface demands an object which
// can generate new Drivers
type DriverFactory interface {
	New() Driver
}

// Server struct represents the server
type Server struct {
	log           *log.Logger
	PORT          string
	DriverFactory DriverFactory
}

// New returns an instance of the Server object
func New(log *log.Logger, PORT string, df DriverFactory) *Server {
	return &Server{log, PORT, df}
}

// Run method starts the TCP server and sets up the TCP client handlers
func (s *Server) Run() {
	listener := s.setupTCPServer()
	defer listener.Close()

	s.setupTCPClientHandler(listener)
}

// setupTCPServer starts a TCP server and returns the listener
func (s *Server) setupTCPServer() net.Listener {
	listener, err := net.Listen("tcp", ":"+s.PORT)
	if err != nil {
		s.log.Fatalf("Listen setup failed: %s", err)
	}

	s.log.Println("Started server on PORT", s.PORT)

	return listener
}

// setupTCPClientHandler sets up the TCP client handler via an infinite loop
func (s *Server) setupTCPClientHandler(l net.Listener) {
	// An infinite loop to listen for any number of TCP clients
	for {
		// Accept WAITS for and returns the next connection
		// to the listener. This is a blocking call.
		conn, err := l.Accept()
		if err != nil {
			s.log.Println("Unable to accept connection: ", err.Error())
			continue
		}

		// Handle the client
		go s.clientHandler(conn)
	}
}

func (s *Server) clientHandler(c net.Conn) {
	// Print the address of the client
	s.log.Println("Connected: ", c.RemoteAddr().String())

	// Create a client
	cl := NewClient(c, s.log, s.DriverFactory.New())

	// Initialise the reader for the client
	cl.InitRead()
}

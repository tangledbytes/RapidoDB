package server

import (
	"log"
	"net"

	"github.com/utkarsh-pro/RapidoDB/db"

	client "github.com/utkarsh-pro/RapidoDB/server/client"
)

// Server is a struct
type Server struct {
	commands chan client.Command
	log      *log.Logger
	PORT     string
	auth     client.Auth
	db       *db.DB
}

// New returns a an instance of server
func New(l *log.Logger, PORT, user, pass string) *Server {
	return &Server{
		commands: make(chan client.Command),
		log:      l,
		PORT:     PORT,
		auth: client.Auth{
			User: user,
			Pass: pass,
		},
		db: db.New(db.NeverExpire),
	}
}

// Setup setups a TCP server and starts accepting connections
func (s *Server) Setup() {

	listener := s.setupTCPServer()

	// Start listening to the commands sent by the clients
	go s.listenCommand()

	s.setupTCPClientHandler(listener)
}

// setupTCPServer sets up a tcp server and returns a listener
func (s *Server) setupTCPServer() net.Listener {
	listener, err := net.Listen("tcp", ":"+s.PORT)
	if err != nil {
		s.log.Fatalf("Listen setup failed: %s", err)
	}

	defer listener.Close()
	s.log.Println("Started server on PORT", s.PORT)

	return listener
}

// setupTCPCkientHandler sets up the TCP client handlers as the
// connection requests comes in to the server
func (s *Server) setupTCPClientHandler(listener net.Listener) {
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
	cl := client.NewClient(c, s.commands, s.log, client.NotAuthenticated, s.db)

	// Inform the client about the connection
	cl.Msg("Successfully connected to RapidoDB. Please run AUTH <user> <pass> to access the DB")

	// Initialise the reader for the client
	cl.InitRead(s.auth)
}

// listenCommand listens for the commands passed on to the
// client channels to the server
//
// DEPRECATED
func (s *Server) listenCommand() {
	for cmd := range s.commands {
		s.log.Println(cmd)
	}
}

package server

import (
	"log"
	"net"

	"github.com/utkarsh-pro/RapidoDB/db"
)

// Server is a struct
type Server struct {
	commands chan command
	log      *log.Logger
	PORT     string
	auth     Auth
	db       *db.DB
}

// New returns a an instance of server
func New(l *log.Logger, PORT string, user string, pass string) *Server {
	return &Server{
		commands: make(chan command),
		log:      l,
		PORT:     PORT,
		auth: Auth{
			user: user,
			pass: pass,
		},
		db: db.New(db.NeverExpire),
	}
}

// Setup setups a TCP server and starts accepting connections
func (s *Server) Setup() {
	// Setup the TCP server
	listener, err := net.Listen("tcp", ":"+s.PORT)
	if err != nil {
		s.log.Fatalf("Listen setup failed: %s", err)
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
	cl := newClient(c, s.commands, s.log, false, s.db)

	// Inform the client about the connection
	cl.msg("Successfully connected to RapidoDB. Please run AUTH <user> <pass> to access the DB")

	// Initialise the reader for the client
	cl.initRead(s.auth)
}

// listenCommands listens for the commands passed on to the
// client channels to the server
func (s *Server) listenCommands() {
	for cmd := range s.commands {
		s.log.Println(cmd)
	}
}

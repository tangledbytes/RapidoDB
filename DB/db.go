/*
   db package will glue all the different layers of the database by providing
   necessary abstractions to each layer. The following is the architecture of RapidoDB

                  TCP SERVER
   ----------------------------------------
   |              TCP CLIENT              |    <==== TRANSPORT LAYER
   ----------------------------------------
   | RQL LEXER | RQL PARSER | RQL DRIVER  |    <==== TRANSLATION LAYER
   ----------------------------------------
   |               OBSERVER               |    <==== OBSERVER LAYER
   ----------------------------------------
   |        SECURITY  |  MANAGER          |    <==== CLIENT MANAGEMENT LAYER
   ----------------------------------------
   |        STORE API | RAW DATA          |    <==== STORAGE LAYER
   ----------------------------------------

   Each layer here is completey independent of the implementation of another layer

*/

package db

import (
	"log"
	"net"

	"github.com/utkarsh-pro/RapidoDB/manage"
	"github.com/utkarsh-pro/RapidoDB/store"
)

// RapidoMSG is the ascii logo for rapidoDB
const RapidoMSG = `
************************************************
   ____             _     _       ____  ____  
  |  _ \ __ _ _ __ (_) __| | ___ |  _ \| __ ) 
  | |_) / _  |  _ \| |/ _  |/ _ \| | | |  _ \ 
  |  _ < (_| | |_) | | (_| | (_) | |_| | |_) |
  |_| \_\__,_| .__/|_|\__,_|\___/|____/|____/ 
             |_|                              

************************************************
`

// RapidoDB struct represents the server
type RapidoDB struct {
	// log will be used internally for logging
	log *log.Logger

	// PORT on which the server should run
	PORT string

	// Store that the RapidoDB will be using internally
	store *store.Store

	// Store that RapidoDB uses to store the DB users info
	usersStore *store.Store
}

// New returns an instance of the Server object
func New(log *log.Logger, PORT, username, password, bckpath string) *RapidoDB {
	// Create a new store for the database
	storage := prepareStorageLayer(log, bckpath+"/rapido.db")

	// Create a new store for the users
	usersDB := store.New(store.NeverExpire, log, bckpath+"/rapido_user.db")

	usersDB.Set(username,
		manage.NewDBUser(username, password, manage.AdminAccess, manage.Events{}), usersDB.DefaultExpiry(),
	)

	return &RapidoDB{log, PORT, storage, usersDB}
}

// Run method starts the TCP server and sets up the TCP client handlers
func (s *RapidoDB) Run() {
	listener := s.setupTCPServer()
	defer listener.Close()

	s.setupTCPClientHandler(listener)
}

// setupTCPServer starts a TCP server and returns the listener
func (s *RapidoDB) setupTCPServer() net.Listener {
	listener, err := net.Listen("tcp", ":"+s.PORT)
	if err != nil {
		s.log.Fatalf("Listen setup failed: %s", err)
	}

	s.log.Println("Started server on PORT", s.PORT)
	s.log.Println("Accepting Connections")

	return listener
}

// setupTCPClientHandler sets up the TCP client handler via an infinite loop
func (s *RapidoDB) setupTCPClientHandler(l net.Listener) {
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

func (s *RapidoDB) clientHandler(c net.Conn) {
	// Print the address of the client
	s.log.Println("Connected: ", c.RemoteAddr().String())

	// get the client manager layer
	sl := prepareClientManagerLayer(s.store, s.usersStore)

	// get the observer layer and the private event bus
	ol, eb := prepareObserverLayer(sl)

	// get the translation layer
	tl := prepareTranslationLayer(ol)

	// get the transporter
	trl := prepareTransportLayer(c, s.log, tl)

	// setup transport extension using the private event bus
	prepareTransportExt(trl, eb)

	// Initialise the reader for the client
	trl.InitRead()
}

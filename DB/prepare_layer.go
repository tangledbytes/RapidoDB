package db

import (
	"log"
	"net"

	"github.com/utkarsh-pro/RapidoDB/manage"
	"github.com/utkarsh-pro/RapidoDB/rql"
	"github.com/utkarsh-pro/RapidoDB/store"
	"github.com/utkarsh-pro/RapidoDB/transport"
)

// prepareStorageLayer prepares the storage layer
func prepareStorageLayer() *store.Store {
	return store.New(store.NeverExpire)
}

// prepareClientManagerLayer takes in a store and a userdb which it uses
// to prepare the client manager layer which also adds security to the database
func prepareClientManagerLayer(store *store.Store, userdb *store.Store) *manage.SecureDB {
	return manage.New(store, userdb)
}

// prepareTranslationLayer takes in a securedb and creates a translation
// driver for the database which enables the database to understand RQL
func prepareTranslationLayer(store *manage.SecureDB) *rql.Driver {
	return rql.New(store)
}

// preparetransportLayer takes in the connection parameter, logger and a translation driver
// to create a transport layer which takes in the remote commands and returns the results
func prepareTransportLayer(c net.Conn, l *log.Logger, d *rql.Driver) *transport.Client {
	return transport.New(c, l, d)
}

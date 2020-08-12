package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"time"
)

type persistor struct {
	interval time.Duration
	bckup    string
	sigStop  chan bool
}

// newPersistor returns a pointer to a new instance of the persistor
func newPersistor(interval time.Duration, bckup string) *persistor {
	return &persistor{interval, bckup, make(chan bool)}
}

// setupPersistor sets up the persistor and a mechanism to
// also setup the function to retrieve data from the disk
func setupPersistor(store *Store) {
	if store.persistor.interval > 0 {
		llog(store.log, "Setting up persistor...")
		runPersistor(store)
	}
}

// runPersistor spins up the persistor in a goroutine to
// make it non blocking while the retrieveData method is started
// up on the main thread to block the startup until the data has
// loaded into the memory
func runPersistor(store *Store) {
	// Run the persistor in a goroutine
	go store.persistor.persist(store)

	// Load the data in the main thread
	store.persistor.retrieve(store)
}

// stopPersistor simply sends a stop signal to the persistor
// func stopPersistor(store *Store) {
// 	store.persistor.sigStop <- true
// }

// persist stores the data onto the disk at regular
// intervals
func (p *persistor) persist(store *Store) error {
	ticker := time.NewTicker(p.interval)
	for {
		select {
		case <-ticker.C:
			// If no backup file is provided
			// then don't persist
			if p.bckup == "" {
				// Stop the ticker
				llog(store.log, "No backup location provided, skipping persistence")
				ticker.Stop()

				return nil
			}

			osf, err := os.Create(p.bckup)
			if err != nil {
				llog(store.log, err)
				return err
			}

			if err = save(store, osf); err != nil {
				llog(store.log, err)
				osf.Close()
				return err
			}

			if err = osf.Close(); err != nil {
				llog(store.log, err)
				return err
			}

		case <-p.sigStop:
			ticker.Stop()
			return nil
		}
	}
}

// retrieve retreives data from the disk and
// store in the main memory of the program
func (p *persistor) retrieve(store *Store) error {
	if p.bckup == "" {
		llog(store.log, "No backup location provided, skipping data retrieval")
		return nil
	}

	llog(store.log, "Started retrieving data from", p.bckup)

	osf, err := os.Open(p.bckup)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		llog(store.log, err)
		return err
	}

	if err = load(store, osf); err != nil {
		llog(store.log, err)
		osf.Close()
		return err
	}

	if err = osf.Close(); err != nil {
		llog(store.log, err)
	}

	llog(store.log, "Completed data retrieval from", p.bckup)
	return err
}

/////////////////////////// HELPER FUNCTIONS //////////////////////////

// save is the internal function which encodes the
// store data as json
func save(store *Store, w io.Writer) (err error) {
	defer func() {
		if x := recover(); x != nil {
			err = fmt.Errorf("Failed to store data on the disk")
		}
	}()

	store.RLock()
	b, err := json.Marshal(store.data)
	store.RUnlock()

	_, err = w.Write(b)
	return err
}

// load decodes the json data
func load(store *Store, r io.Reader) error {
	data := make(map[string]Item)

	// Read into the bytes
	b, err := read(r)
	if err != nil {
		return err
	}
	// Unmarshal the data
	err = json.Unmarshal(b, &data)

	store.Lock()
	store.data = data
	store.Unlock()

	return err
}

// read reads the data from the file
func read(r io.Reader) ([]byte, error) {
	var b []byte
	buf := make([]byte, 1024)

	for {
		// Read the chunk
		n, err := r.Read(buf)
		if err != nil && err != io.EOF {
			return []byte{}, err
		}
		if n == 0 {
			break
		}
		// Write the chunk
		b = append(b, buf[:n]...)
	}

	return b, nil
}

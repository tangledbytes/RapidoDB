package store

import (
	"runtime"
	"time"
)

// janitor is responsible for cleaning up the
// expired items at a regular interval
type janitor struct {
	interval time.Duration
	sigStop  chan bool
}

// newJanitor creates a new janitor and sets the interval,
// it returns a pointer to the janitor
func newJanitor(interval time.Duration) *janitor {
	return &janitor{
		interval: interval,
		sigStop:  make(chan bool),
	}
}

// run methods starts the janitor and executes the
// "DeleteExpired" method on the store at the regular intervals.
// This method can be stopped by the garbage collector though
func (j *janitor) run(store *Store) {
	ticker := time.NewTicker(j.interval)
	for {
		select {
		case <-ticker.C:
			store.DeleteExpired()
		case <-j.sigStop:
			ticker.Stop()
			return
		}
	}
}

// setupJanitor takes in the store and sets up a janitor for that
// store. It also adds a finalizer function to stop the janitor
// whenever required
func setupJanitor(store *Store) {
	if store.janitor.interval > 0 {
		runJanitor(store)
		runtime.SetFinalizer(store, stopJanitor)
	}
}

// runJanitor starts the janitor in a goroutine
func runJanitor(store *Store) {
	go store.janitor.run(store)
}

// stopJanitor stops the janitor by sending a stop signal
// to the janitor. This function is intended to be used by
// the go runtime as a finalizer function
func stopJanitor(store *Store) {
	store.janitor.sigStop <- true
}

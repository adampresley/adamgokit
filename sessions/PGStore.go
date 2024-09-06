package sessions

import (
	"fmt"
	"time"

	"github.com/antonlindstrom/pgstore"
)

/*
NewPGStore is a convenience method around the "pgstore" library
that initializes a session store and returns a method for cleanup.
*/
func NewPGStore(dsn, sessionKey string) (*pgstore.PGStore, func(), error) {
	store, err := pgstore.NewPGStore(dsn, []byte(sessionKey))

	if err != nil {
		return store, func() {}, fmt.Errorf("could not initialize postges session storage: %w", err)
	}

	cleaner := func() {
		store.Close()
		store.StopCleanup(store.Cleanup(time.Minute * 5))
	}

	return store, cleaner, nil
}

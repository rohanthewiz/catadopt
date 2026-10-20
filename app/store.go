// Store: bytdb-backed persistence for adoption requests.
//
// The shape follows grmob's examples/todoapp store, which is the documented
// GrMob persistence pattern. The store opens lazily on the first render pass
// (openStore is cheap to call repeatedly), and the set of requested cat IDs it
// read at open seeds App's NewState, so a request made in an earlier launch is
// already in the very first tree rather than popping in a patch later. Every
// later request goes through profileScreen's ask, the single mutation point,
// as a synchronous write-through. A bytdb single-row insert is microseconds
// (WAL, fsync before ack), so riding the tap's event path costs nothing
// perceptible and leaves no dirty flag or debounce to get wrong.
//
//	first render ──▶ openStore ──▶ SELECT cat_id ──▶ NewState initial
//	"Ask to adopt" ──▶ ask ──▶ requested.Set (UI)
//	                      └──▶ store.addRequest (disk)
//
// Persistence is strictly optional. With no data directory registered
// (mobile.SetDataDir unset: the browser build, bare unit tests) openStore
// returns nil, and every method below is nil-receiver-safe, so callers use the
// store unconditionally instead of branching at each site. The Android and iOS
// shells register their app-private files directory before the first render.
//
// This only remembers requests on the device. Sending them to the shelter
// (email, a form backend, its API) is separate work; requested_at is kept so
// that a later sync can report when each request was made.

package app

import (
	"log"
	"path/filepath"
	"sync"
	"time"

	"github.com/rohanthewiz/bytdb"
	bsql "github.com/rohanthewiz/bytdb/sql"
	"github.com/rohanthewiz/grmob/mobile"
)

// storeFile is the database's name inside the shell's data directory.
const storeFile = "catadopt.bytdb"

// adoptionStore wraps one open bytdb engine plus the snapshot read at open.
// The snapshot is a one-time read, not a live mirror: after the first render
// App's state is the in-process source of truth, and the table only has to be
// right for the next launch. (Tests simulate a relaunch with closeStore, which
// forces a fresh read.)
type adoptionStore struct {
	eng  *bytdb.Engine
	db   *bsql.DB
	path string

	requested map[string]bool // cat IDs with a request, as of open
}

// The open store is a package-level singleton because the mobile bridge runs
// one app per process, and bytdb holds an exclusive lock on its file, so a
// second Open of the same path would fail. The mutex covers open/close races
// between the render goroutine and tests.
var (
	storeMu sync.Mutex
	store   *adoptionStore
)

// openStore returns the store for the current data directory, opening it on
// first use. It is called on every render pass, and after the first call it is
// a mutex acquire and a string compare. A changed data directory (a test
// moving to a fresh t.TempDir) closes the old engine and opens the new path.
func openStore() *adoptionStore {
	storeMu.Lock()
	defer storeMu.Unlock()

	dir := mobile.DataDir()
	if dir == "" {
		return nil // no shell-provided directory: run in memory
	}
	path := filepath.Join(dir, storeFile)
	if store != nil && store.path == path {
		return store
	}
	if store != nil {
		if err := store.eng.Close(); err != nil {
			log.Printf("catadopt: closing previous store: %v", err)
		}
		store = nil
	}

	eng, err := bytdb.Open(path)
	if err != nil {
		// A failed open degrades to in-memory rather than crashing the app:
		// losing persistence beats losing the UI. Nothing is cached, so the
		// next render retries the Open; on a bad path that fails fast.
		log.Printf("catadopt: opening store at %s: %v", path, err)
		return nil
	}
	db := bsql.New(eng)

	// One row per cat: the primary key makes a repeat request a no-op (see
	// addRequest), so the table is a set of cat IDs, the same shape as the
	// requested map it seeds.
	if _, err := db.Exec(
		`CREATE TABLE IF NOT EXISTS adoption_requests (cat_id text PRIMARY KEY, requested_at timestamptz)`,
	); err != nil {
		log.Printf("catadopt: creating adoption_requests table: %v", err)
		_ = eng.Close()
		return nil
	}

	res, err := db.Exec(`SELECT cat_id FROM adoption_requests`)
	if err != nil {
		log.Printf("catadopt: loading adoption requests: %v", err)
		_ = eng.Close()
		return nil
	}
	s := &adoptionStore{eng: eng, db: db, path: path, requested: make(map[string]bool, len(res.Rows))}
	for _, row := range res.Rows {
		// A request for a cat no longer in the catalogue is kept, not
		// pruned: the map is only ever read by Cat.ID, so a stray key is
		// harmless, and the row still matters to whoever processes requests.
		if id, ok := row[0].(string); ok {
			s.requested[id] = true
		}
	}

	store = s
	return s
}

// closeStore releases the engine and forgets the singleton. The app never
// calls it, since the engine lives for the process and bytdb's WAL makes a
// hard kill safe, but tests use it to simulate a relaunch: close, render
// again, and the next snapshot comes from disk instead of memory.
func closeStore() {
	storeMu.Lock()
	defer storeMu.Unlock()
	if store == nil {
		return
	}
	if err := store.eng.Close(); err != nil {
		log.Printf("catadopt: closing store: %v", err)
	}
	store = nil
}

// snapshot returns the requests read at open, as a fresh map. It is copied
// because NewState keeps a reference to its initial value, and ask treats the
// state's map as immutable (copy, change, Set); two App mounts over one open
// store must not share a map.
func (s *adoptionStore) snapshot() map[string]bool {
	out := map[string]bool{}
	if s == nil {
		return out
	}
	for id := range s.requested {
		out[id] = true
	}
	return out
}

// addRequest records a request for one cat. ON CONFLICT DO NOTHING keeps the
// first request's timestamp if the same cat is asked for twice; the UI
// disables the button once asked, but a second App mount (or a future undo)
// shouldn't be able to turn a duplicate into an error.
//
// A failure is logged and swallowed. The in-memory map has already updated
// and stays authoritative for this process; the worst case is that this one
// request is missing after the next launch.
func (s *adoptionStore) addRequest(catID string, at time.Time) {
	if s == nil {
		return
	}
	const q = `INSERT INTO adoption_requests (cat_id, requested_at) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	if _, err := s.db.Exec(q, catID, at.UTC()); err != nil {
		log.Printf("catadopt: recording adoption request for %s: %v", catID, err)
	}
}

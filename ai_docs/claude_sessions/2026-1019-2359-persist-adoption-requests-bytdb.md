# Persist adoption requests with bytdb

Session: `10dcabb8-0280-4a9a-b570-d5553105b8c0`
**Date:** 2026-10-19 23:59 · **Repo:** https://github.com/rohanthewiz/catadopt (ed1a969 → this commit) · **grmob:** v0.7.2 · **bytdb:** v0.22.0 (new direct dependency)

## Goal

Follows `2026-1019-2347-warm-theme-android-and-ios`. The user asked to
"persist adoption requests with bytdb", the first half of N-002.

## Approach

The pattern is grmob's `examples/todoapp/store.go`, which
`docs/concepts/state-and-hooks.md` (Persistence) names as the recommended
shape:

- The native shells register a writable directory through
  `mobile.SetDataDir`. Both vendored shells already do this
  (`android/.../GomobileBridge.kt`, `ios/GrMob/App/GomobileBridge.swift`),
  so no shell change was needed.
- The store opens lazily on the first render. The snapshot it reads at open
  seeds `NewState`, so saved requests are in the first tree, not patched in
  later.
- Writes go to memory first, then disk, synchronously on the tap. There is
  no dirty flag or debounce.
- With no data directory (browser build, bare tests) the store is nil and
  every method is nil-receiver-safe. The app runs in memory as before.

## Changes

- **`app/store.go` (new):** an `adoptionStore` singleton, guarded by a mutex,
  over `catadopt.bytdb` in `mobile.DataDir()`.
  - Table: `adoption_requests (cat_id text PRIMARY KEY, requested_at timestamptz)`.
  - `openStore` reopens if the data directory changes (tests).
  - A failed open or write is logged (`catadopt: …`) and the app degrades
    to in-memory.
  - `snapshot()` returns a fresh `map[string]bool`, because NewState keeps a
    reference to its initial value.
  - `addRequest` uses `ON CONFLICT DO NOTHING`, which keeps the first
    request's timestamp.
  - `closeStore` is used by tests only, to simulate a relaunch.
- **`app/app.go`:** `requested` is seeded from `openStore().snapshot()`
  instead of an empty map.
- **`app/profile_screen.go`:** `ask` calls
  `openStore().addRequest(c.ID, time.Now())` after `requested.Set`.
- **`app/app_test.go`:** a `withDataDir(t)` helper (temp directory, with a
  cleanup that closes the store and unregisters the directory) and three new
  tests:
  - `TestAdoptionSurvivesRelaunch`: asks for Marmalade, calls
    `closeStore`, mounts a fresh App, and expects "Request sent". It also
    checks that the snapshot from disk holds only `marmalade`. Only the
    Gallery's front card has a "Meet" button, so the test can't open other
    cats' profiles without paging.
  - `TestAddRequestTwice`: a duplicate request leaves one row.
  - `TestNoDataDirRunsInMemory`: the nil store is safe.
- **`README.md`:** the feature list notes that requests are saved on the
  natives but not in the browser. The layout lists `store.go`.
- **`go.mod`:** `github.com/rohanthewiz/bytdb v0.22.0` (btypedb v0.9.0
  indirect).

## Verification

- `go test ./app`: all 8 tests pass. No `catadopt:` store errors were
  logged, so binding a `time.Time` to `timestamptz` works.
- `GOOS=js GOARCH=wasm go build ./wasm` succeeds: bytdb compiles for the
  browser, though it is never opened there.
- `go vet ./app` passes for `GOOS=android GOARCH=arm64` and
  `GOOS=ios GOARCH=arm64`.
- Not run on a device or simulator yet (N-010).

## Next

Closed: None. Declined: None. Raised: N-010.
Deferred: None. Promoted: None. Moved: None.
Updated: N-002 (persistence done; sending requests to the shelter remains).
Full list: `ai_docs/todo/next-list.md`.

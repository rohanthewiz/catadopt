package app

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/rohanthewiz/grmob/core"
	"github.com/rohanthewiz/grmob/mobile"
	"github.com/rohanthewiz/grmob/render"
)

// The fastest loop there is: render.Manager is the same engine every target
// runs, and it needs no browser, simulator or device. Taps go through the
// same dispatch path a finger does.

// TestMain turns on debug mode, which audits every render pass for hook-order
// drift and duplicate keys; assertNoConcerns fails the test on any finding.
func TestMain(m *testing.M) {
	core.SetDebugMode(true)
	m.Run()
}

// TestGalleryShowsEveryCat checks the home screen: the app bar count, one
// card picture per cat, and the short description on the front card.
func TestGalleryShowsEveryCat(t *testing.T) {
	core.ClearConcerns()
	mgr := render.New(core.NewContext(), App)
	defer mgr.Close()

	cats := catalogue()
	root := tree(t, mgr)

	if !shows(root, "9 cats looking for a home") {
		t.Fatal("app bar does not count the cats")
	}
	for _, c := range cats {
		// Each photo appears at least twice: on its card and in the
		// thumbnail strip.
		if n := count(root, func(n *node) bool { return n.Type == "Image" && n.Props["src"] == photoURL(c.Photo) }); n < 2 {
			t.Errorf("%s: photo drawn %d times, want card + thumbnail", c.Name, n)
		}
	}
	if !shows(root, cats[0].Short) {
		t.Error("front card does not show its short description")
	}
	assertNoConcerns(t)
}

// TestAppTheme checks the screen is painted with appTheme's Surface on the
// safe area (the status-bar strip on the natives), which proves the theme
// reaches the routes: they read it from the frame context, not the App's.
func TestAppTheme(t *testing.T) {
	core.ClearConcerns()
	mgr := render.New(core.NewContext(), App)
	defer mgr.Close()

	want := appTheme.Colors.Surface
	root := tree(t, mgr)
	sa := find(root, func(n *node) bool { return n.Type == "SafeArea" })
	if sa == nil || sa.Style["Background"] != want {
		t.Fatalf("safe area background = %v, want appTheme's Surface %s", sa.Style["Background"], want)
	}
	assertNoConcerns(t)
}

// TestTapCardForInfo taps the middle of the picture (the Gallery's centre tap
// zone) and expects the card to flip to its details: the Info control then
// reports itself expanded, which is the flipped state screen readers hear.
func TestTapCardForInfo(t *testing.T) {
	core.ClearConcerns()
	mgr := render.New(core.NewContext(), App)
	defer mgr.Close()

	if infoExpanded(t, mgr) {
		t.Fatal("card starts flipped")
	}
	// The tap zones are the Boxes with an onClick in the stage's overlay
	// Row: previous, flip, next, in that order.
	zones := tapZones(tree(t, mgr))
	if len(zones) != 3 {
		t.Fatalf("found %d tap zones over the card, want 3", len(zones))
	}
	mgr.DispatchCallback(zones[1])
	if !infoExpanded(t, mgr) {
		t.Fatal("tapping the middle of the picture did not show the cat's info")
	}
	assertNoConcerns(t)
}

// TestMeetAndAdopt walks the whole adoption path: open a profile from the
// card's action, ask to adopt, go back, reopen, and expect the request to
// have survived the round trip (it lives above the Navigator).
func TestMeetAndAdopt(t *testing.T) {
	core.ClearConcerns()
	mgr := render.New(core.NewContext(), App)
	defer mgr.Close()

	tap(t, mgr, "Meet Marmalade")
	root := tree(t, mgr)
	if !shows(root, catalogue()[0].Story) {
		t.Fatal("profile screen did not open")
	}

	tap(t, mgr, "Ask to adopt Marmalade")
	if !shows(tree(t, mgr), "Thanks! The shelter will contact you about Marmalade.") {
		t.Fatal("no confirmation after asking to adopt")
	}

	tap(t, mgr, "‹") // AppBar's automatic back control
	if !shows(tree(t, mgr), "9 cats looking for a home") {
		t.Fatal("back did not return to the gallery")
	}

	tap(t, mgr, "Meet Marmalade")
	if find(tree(t, mgr), buttonLabeled("Request sent")) == nil {
		t.Fatal("adoption request was lost after navigating away")
	}
	assertNoConcerns(t)
}

// TestAdoptionSurvivesRelaunch points the store at a temp data directory, as
// the native shells do with their files directory, asks to adopt, then
// simulates a relaunch: closeStore drops the engine and a fresh App mount
// reopens the file, so "Request sent" can only come from disk.
func TestAdoptionSurvivesRelaunch(t *testing.T) {
	core.ClearConcerns()
	withDataDir(t)

	// Marmalade, because only the Gallery's front card has a "Meet" action.
	mgr := render.New(core.NewContext(), App)
	tap(t, mgr, "Meet Marmalade")
	tap(t, mgr, "Ask to adopt Marmalade")
	mgr.Close()

	closeStore()
	mgr = render.New(core.NewContext(), App)
	defer mgr.Close()

	tap(t, mgr, "Meet Marmalade")
	if find(tree(t, mgr), buttonLabeled("Request sent")) == nil {
		t.Fatal("adoption request was lost across a relaunch")
	}
	// Only the cat asked for: the snapshot read from disk must not mark the
	// others. Checked on the store, since the other cats' profiles are only
	// reachable by paging the Gallery.
	if got := openStore().snapshot(); len(got) != 1 || !got["marmalade"] {
		t.Fatalf("snapshot after relaunch = %v, want only marmalade", got)
	}
	assertNoConcerns(t)
}

// TestAddRequestTwice checks a repeat request for the same cat is absorbed
// (ON CONFLICT DO NOTHING) rather than logged as a failed write, and still
// reads back as one request.
func TestAddRequestTwice(t *testing.T) {
	withDataDir(t)

	s := openStore()
	if s == nil {
		t.Fatal("store did not open in the temp data directory")
	}
	s.addRequest("rusty", time.Now())
	s.addRequest("rusty", time.Now())

	res, err := s.db.Exec(`SELECT count(*) FROM adoption_requests WHERE cat_id = $1`, "rusty")
	if err != nil {
		t.Fatalf("counting requests: %v", err)
	}
	if n := fmt.Sprint(res.Rows[0][0]); n != "1" {
		t.Fatalf("rusty has %s request rows, want 1", n)
	}
}

// TestNoDataDirRunsInMemory: with no directory registered (the browser build)
// openStore is nil, and the nil store's methods are safe no-ops.
func TestNoDataDirRunsInMemory(t *testing.T) {
	mobile.SetDataDir("")
	s := openStore()
	if s != nil {
		t.Fatal("store opened without a data directory")
	}
	s.addRequest("frost", time.Now()) // must not panic
	if got := s.snapshot(); len(got) != 0 {
		t.Fatalf("nil store snapshot = %v, want empty", got)
	}
}

// TestPhotoURL pins the Wikimedia thumbnail layout photoURL relies on.
func TestPhotoURL(t *testing.T) {
	got := photoURL("3/3a/Cat03.jpg")
	want := "https://upload.wikimedia.org/wikipedia/commons/thumb/3/3a/Cat03.jpg/960px-Cat03.jpg"
	if got != want {
		t.Fatalf("photoURL = %s, want %s", got, want)
	}
}

// --- helpers ---------------------------------------------------------------

// withDataDir registers a fresh temp directory as the shell's data directory
// for one test. Cleanups run last-in first-out, so the store is closed (and
// the directory unregistered) before t.TempDir removes the files; the other
// tests then run with no directory, in memory, as they always have.
func withDataDir(t *testing.T) {
	t.Helper()
	mobile.SetDataDir(t.TempDir())
	t.Cleanup(func() {
		closeStore()
		mobile.SetDataDir("")
	})
}

type node struct {
	Type     string
	Props    map[string]any
	Style    map[string]any
	Children []*node
}

// tree renders and parses the current tree. Callback IDs are assigned per
// pass, so each tap reads its ID from a fresh tree rather than reusing one.
func tree(t *testing.T, mgr *render.Manager) *node {
	t.Helper()
	var root node
	if err := json.Unmarshal([]byte(mgr.RenderInitial()), &root); err != nil {
		t.Fatalf("tree is not valid JSON: %v", err)
	}
	return &root
}

func find(n *node, pred func(*node) bool) *node {
	if n == nil || pred(n) {
		return n
	}
	for _, c := range n.Children {
		if f := find(c, pred); f != nil {
			return f
		}
	}
	return nil
}

func count(n *node, pred func(*node) bool) int {
	if n == nil {
		return 0
	}
	total := 0
	if pred(n) {
		total++
	}
	for _, c := range n.Children {
		total += count(c, pred)
	}
	return total
}

func buttonLabeled(label string) func(*node) bool {
	return func(n *node) bool { return n.Type == "Button" && n.Props["label"] == label }
}

func tap(t *testing.T, mgr *render.Manager, label string) {
	t.Helper()
	b := find(tree(t, mgr), buttonLabeled(label))
	if b == nil {
		t.Fatalf("no Button labeled %q", label)
	}
	mgr.DispatchCallback(b.Props["onClick"].(string))
}

func shows(root *node, text string) bool {
	return find(root, func(n *node) bool { return n.Type == "Text" && n.Props["content"] == text }) != nil
}

// tapZones returns the onClick IDs of the Gallery's overlay tap zones: the
// first Row whose children are all clickable Boxes.
func tapZones(root *node) []string {
	row := find(root, func(n *node) bool {
		if n.Type != "Row" || len(n.Children) == 0 {
			return false
		}
		for _, c := range n.Children {
			if c.Type != "Box" || c.Props["onClick"] == nil {
				return false
			}
		}
		return true
	})
	if row == nil {
		return nil
	}
	ids := make([]string, len(row.Children))
	for i, c := range row.Children {
		ids[i] = c.Props["onClick"].(string)
	}
	return ids
}

// infoExpanded reads whether the Gallery's Info (flip) control reports the
// current card as turned over. The control's label carries a glyph prefix,
// so it is matched by suffix rather than exactly.
func infoExpanded(t *testing.T, mgr *render.Manager) bool {
	t.Helper()
	b := find(tree(t, mgr), func(n *node) bool {
		l, _ := n.Props["label"].(string)
		return n.Type == "Button" && strings.HasSuffix(l, "Info")
	})
	if b == nil {
		t.Fatal("no Info button under the gallery")
	}
	v, _ := b.Style["AccessibilityExpanded"].(string)
	return v == "true"
}

func assertNoConcerns(t *testing.T) {
	t.Helper()
	if cs := core.Concerns(); len(cs) != 0 {
		t.Fatalf("debug concerns raised:\n%s", core.DumpConcerns())
	}
}

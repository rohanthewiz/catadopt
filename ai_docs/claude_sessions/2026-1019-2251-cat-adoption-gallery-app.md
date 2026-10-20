# Cat adoption app on comps.Gallery, dark mode

Session: `ada40c4b-9595-4b2b-bcbc-e2e6ecae2cf7`
**Date:** 2026-10-19 22:51 · **Repo:** https://github.com/rohanthewiz/catadopt (first commit) · **grmob:** v0.7.0

## Goal

The request came from the user's cats-todo backlog: use the GrMob skill to
build an app showing cats for adoption. The requirements were:

- use the Gallery widget, in dark mode;
- give a short description of each cat;
- let the user click somewhere on the picture card for more info.

Mid-session the user added: "Also show me some sample cat pics".

## What was built

- **Scaffold:** `grmob new catadopt -name "Cat Adoption" -id com.rohanthewiz.catadopt`.
  `comps.Gallery` was only on grmob master at first (it is not in v0.6.0),
  so the app tracked master, then moved to **v0.7.0** once that was cut.
- **`app/app.go`:** `App` wraps a `core.Navigator` in
  `core.WithTheme(core.DarkTheme, …)`. Dark mode lives in the view tree, so
  every target (Android, iOS, browser) gets it, not only the browser host.
  An adoption-request map (`State[map[string]bool]`, keyed by `Cat.ID`) sits
  above the Navigator so it survives the profile route being popped. Routes
  get the *frame* context: it inherits DarkTheme, while the App's own ctx
  sits outside `WithTheme`.
- **`app/cats.go`:** the `Cat` record and a nine-cat sample catalogue.
  - The photos are real Wikimedia Commons files, requested as 960px
    thumbnails. Commons only serves fixed widths, and arbitrary widths 404.
  - Each photo was downloaded and looked at so its description matches the
    cat: ginger tabbies, brown tabbies, an Abyssinian, a ginger kitten and a
    blue-eyed tabby-and-white.
  - The names, ages and personalities are invented.
  - `details()` feeds both the back of the card and the profile.
- **`app/gallery_screen.go`:** an AppBar ("9 cats looking for a home"), a
  hint line, and the Gallery.
  - Thumbnail indicator, so the strip previews every cat (the "sample cat
    pics").
  - Fullscreen zoom, `Wrap`, three description lines.
  - Labels: Details → "Info". Action: "Meet <name>", which pushes the
    profile.
  - The card's back is the details list. Tapping the middle of the photo
    flips to it, and tapping the edges steps between cats.
- **`app/profile_screen.go`:** a hero image (`ContentModeFill`), name, coat
  and age, the story, a `KeyValueList`, an "Ask to adopt <name>" button
  (which then reads "Request sent", disabled), a success Banner, and a
  Commons credit link.
- **`app/app_test.go`:** five `render.Manager` tests in debug mode, each
  asserting that no debug concerns were raised:
  - every cat's photo appears on its card and in the thumbnails;
  - DarkTheme's background is on the safe area;
  - tapping the middle zone flips the card (Info reports expanded);
  - the Meet → Ask → back → reopen flow keeps the request;
  - the `photoURL` layout.
- **`wasm/index.html`:**
  - Re-rendered from grmob master with `grmob web -refresh`, then darkened:
    `color-scheme: dark`, a near-black page, and a `#1C1C1E` screen behind
    the first frame.
  - Added `<meta name="darkreader-lock">`.
- **README:** a section on what the app does, and the new files under
  Layout.

## Found along the way

- **grmob bug, fixed upstream as N-115 (dcf9274, released in v0.7.0).**
  Gallery's card layers carry `core.ZIndex` (up to n+2), but its tap-zone
  Row carried none. On the web the cards painted over the zones, so clicking
  the photo did nothing.
  - catadopt first carried a CSS workaround in `index.html`.
  - The fix was made in grmob, with a regression test.
  - grmob was wrapped and v0.7.0 cut (see grmob's
    `2026-1019-2246-gallery-tapzone-zindex`).
  - catadopt was bumped and the workaround deleted. It was then seen working
    in Chrome: the edge stepped to the next cat and the middle flipped.
- **Dark Reader in the user's Chrome** runs in filter mode on localhost and
  ignored `darkreader-lock`. It inverts the dark page to cream and hides
  Gallery's caption band: its reverse filter on `<img>` makes the photo
  paint above the caption. The screenshots used for checks were taken with
  `filter: none !important` injected into the tab only.

## Verification

- `go test ./app` passes; `go vet ./...` and gofmt are clean.
- Seen in Chrome at `./dev.sh -addr :8091`:
  - gallery front, caption and thumbnails;
  - middle-tap flip to info;
  - edge and arrow stepping;
  - the profile screen.
- Not built for Android or iOS.

## Next

- Real listings: replace the sample `catalogue()` with a shelter feed, or a
  bytdb store opened lazily per the GrMob persistence pattern.
- Persist adoption requests (bytdb), and send them somewhere: email, a form
  backend, or the shelter's API. Today they live in memory only.
- Per-photo author and licence on the profile. It links to each Commons
  file page now, which carries them, but doesn't state them.
- Photos are hotlinked from upload.wikimedia.org. Consider bundling or
  caching them for offline use on the natives.
- Filters (kittens, good with kids or dogs). Gallery items are keyed by ID,
  so a filtered list keeps each card's node.
- Non-goal: working around Dark Reader beyond `darkreader-lock`. It is a
  per-user extension setting (turn it off for localhost).

**Validate**

- Build and run on Android (`grmob android -install`) and iOS
  (`grmob ios -run`). Check the dark status-bar strip, the photo tap zones
  (grmob's N-111 has never tapped one on the natives), the fullscreen
  Lightbox, and the profile's back control.
- Screen readers: the card names ("3 of 9, Juniper"), the Info and Zoom
  controls, and the profile heading and Banner.

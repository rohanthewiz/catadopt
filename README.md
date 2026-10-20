# Cat Adoption

A [GrMob](https://github.com/rohanthewiz/grmob) app: views, state and logic in
Go, rendered in the browser through WebAssembly and natively on Android and
iOS from the same package.

Browse cats up for adoption, one large card at a time (`comps.Gallery`, dark
mode). Each card shows the cat's name, age and a short description:

- tap the middle of the photo for the cat's info (age, coat, good with,
  vaccinated, spayed/neutered, fee); tap its edges to step between cats;
- the thumbnail strip previews every cat; Zoom shows the uncropped photo;
- "Meet <name>" opens a profile with the full story and an "Ask to adopt"
  button.

The cats are sample data (`app/cats.go`): the photos are real Wikimedia
Commons images, credited on each profile; the names and stories are invented.

## Develop in the browser

```bash
./dev.sh          # http://localhost:8080 — rebuilds and hot-swaps on every save
go test ./app     # drive the app from a test: no browser, no device
```

| You changed | What happens |
|---|---|
| anything under `app/`, or `wasm/main.go` | `./dev.sh` rebuilds and swaps the module in place |
| `wasm/index.html` | the page reloads |
| `go.mod` (e.g. a grmob upgrade) | rebuild; the runtime JS is re-copied to match |

`./build.sh` is the one-off build. It writes `wasm/main.wasm` and copies the
runtime JS from the grmob version in `go.mod`, so the browser runtime can
never be older or newer than the Go side. Deploy the `wasm/` directory to any
static host that serves `.wasm` as `application/wasm`.

Don't edit `wasm/grmob-runtime.js`: it is replaced on every build. A runtime
bug is a grmob bug — fix it there and upgrade:

```bash
go get github.com/rohanthewiz/grmob@<version> && go mod tidy
```

## Native targets

Install the SDKs when you need them; the browser target never does.

```bash
go run github.com/rohanthewiz/grmob/cmd/grmob doctor            # what's installed, what's missing
go run github.com/rohanthewiz/grmob/cmd/grmob android -install   # APK, installed on a device/emulator
go run github.com/rohanthewiz/grmob/cmd/grmob ios -open          # simulator build, then open Xcode
```

The first native build copies grmob's shell into `android/` or `ios/`. After
that the shell is part of this app: change the icon, permissions or manifest
there. The version it was copied from is recorded beside it, and a build warns
when `go.mod` has moved past it; pass `-refresh` to copy the new shell over it.

The launcher name and application ID come from `grmob.json`.

## Layout

```
app/          the app (every target mounts this):
  app.go              root view: DarkTheme + Navigator, adoption-request state
  cats.go             the Cat record and the sample catalogue
  gallery_screen.go   home screen: the Gallery of cats
  profile_screen.go   one cat's profile and the adoption request
  app_test.go         render.Manager tests in debug mode
ai_docs/      session docs
wasm/         browser host: main.go (webhost.Run) and index.html
build.sh      browser build
dev.sh        dev server with hot reload
grmob.json    name and application ID for the native shells
```

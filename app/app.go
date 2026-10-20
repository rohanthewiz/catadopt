// Package app is Cat Adoption's root view. Every target mounts this package:
// wasm/main.go for the browser, and `grmob android` / `grmob ios`, which bind
// it into the native shells.
package app

import (
	"github.com/rohanthewiz/grmob/core"
	"github.com/rohanthewiz/grmob/mobile"
)

// init registers the root view with the mobile bridge, which is the whole
// integration contract for the native shells: they call the bridge, and the
// bridge renders whatever registered here. The browser host mounts App
// directly (see wasm/main.go), so on that target this registration is unused.
func init() {
	mobile.Register(core.NewContext(), App)
}

// AppName gives gomobile a bindable symbol. gomobile links a bound package
// only if it exports something it can bind, and App — a function-typed value
// — is not, so without this the package and the init above would be dropped
// from the native library and the shell would start with no app.
func AppName() string { return "Cat Adoption" }

// appTheme is the whole app's palette, kept in one place so a change of look
// is one line. AmberTheme is grmob's bundled warm theme: an amber brand
// colour, amber-tinted surfaces and brown ink on its filled controls, over a
// white page. (The app shipped on core.DarkTheme first.) The screens read
// every colour from ctx.Theme(), so nothing else names a palette, apart from
// wasm/index.html's frame, which can't read Go and restates the page colour.
var appTheme = core.AmberTheme

// App is called again on every render pass: it reads state, returns a tree,
// and never mutates the screen itself.
//
//	App ── WithTheme(appTheme)
//	        └ Navigator
//	            ├ galleryScreen   (root route: the Gallery of cats)
//	            └ profileScreen   (pushed by a card's "Meet <name>" action)
//
// The theme is applied here, in the view tree, rather than by handing the
// host a themed context (wasm/main.go's webhost.Run(nil, …)). The tree is
// what every target shares, so this one line themes Android, iOS and the
// browser alike; a host-side theme would have to be repeated per shell.
func App(ctx *core.Context) core.View {
	// State that outlives a screen lives above the Navigator: the profile
	// route is discarded on Pop, but a request made there should still show
	// when the cat's profile is opened again. Keyed by Cat.ID.
	//
	// The initial value comes from the bytdb store (store.go), so requests
	// made in an earlier launch are in the first tree. openStore is called on
	// every pass but opens only once; NewState ignores its initial value after
	// the first pass, so the later snapshots are unused copies of a small map.
	// With no data directory (browser, bare tests) the store is nil and the
	// snapshot is empty: in-memory only, as before.
	requested := core.NewState(ctx, openStore().snapshot())
	cats := catalogue()

	return core.WithTheme(appTheme,
		core.Navigator(func(frame *core.Context) core.View {
			// The route captures the App's state rather than declaring its
			// own, so nothing is lost when frames come and go. It is handed
			// the frame context, not the App's: the frame is a scope of the
			// themed context, so frame.Theme() is appTheme, whereas the
			// App's ctx sits outside WithTheme and still reports the
			// default theme. Navigation works from either (they share one
			// route stack).
			return galleryScreen(frame, cats, requested)
		}),
	)
}

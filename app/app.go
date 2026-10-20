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

// App is called again on every render pass: it reads state, returns a tree,
// and never mutates the screen itself.
//
//	App ── WithTheme(DarkTheme)
//	        └ Navigator
//	            ├ galleryScreen   (root route: the Gallery of cats)
//	            └ profileScreen   (pushed by a card's "Meet <name>" action)
//
// Dark mode is applied here, in the view tree, rather than by handing the
// host a themed context (wasm/main.go's webhost.Run(nil, …)). The tree is
// what every target shares, so this one line makes Android, iOS and the
// browser all dark; a host-side theme would have to be repeated per shell.
func App(ctx *core.Context) core.View {
	// State that outlives a screen lives above the Navigator: the profile
	// route is discarded on Pop, but a request made there should still show
	// when the cat's profile is opened again. Keyed by Cat.ID.
	requested := core.NewState(ctx, map[string]bool{})
	cats := catalogue()

	return core.WithTheme(core.DarkTheme,
		core.Navigator(func(frame *core.Context) core.View {
			// The route captures the App's state rather than declaring its
			// own, so nothing is lost when frames come and go. It is handed
			// the frame context, not the App's: the frame is a scope of the
			// themed context, so frame.Theme() is DarkTheme, whereas the
			// App's ctx sits outside WithTheme and still reports the
			// default theme. Navigation works from either (they share one
			// route stack).
			return galleryScreen(frame, cats, requested)
		}),
	)
}

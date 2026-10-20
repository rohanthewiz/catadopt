package app

import (
	"fmt"

	"github.com/rohanthewiz/grmob/comps"
	"github.com/rohanthewiz/grmob/core"
)

// galleryHeight is the stage height in px. Tall enough that a portrait photo
// reads as a portrait and the caption band (name, age, two-to-three lines of
// description) doesn't cover most of the cat; short enough that the arrows,
// thumbnail strip and buttons still fit under it on a ~800pt phone screen.
const galleryHeight = 440

// galleryScreen is the root route: an app bar, a one-line hint, and the
// Gallery of adoptable cats.
//
// How "click the card for info" works — the Gallery lays three invisible tap
// zones over the current card, so no extra wiring is needed:
//
//	┌──────────── card ────────────┐
//	│  ‹    │                │   › │   leading / trailing quarter: previous / next
//	│ prev  │  tap: Info     │ next│   middle half: flip to the back (Details)
//	│       │  hold: zoom    │     │
//	├───────┴────────────────┴─────┤
//	│ Marmalade                    │   caption band: Title, Subtitle,
//	│ 4 yrs · Male                 │   Description (the short blurb)
//	│ A gentle ginger gentleman…   │
//	└──────────────────────────────┘
//	  ‹   [thumb][thumb][thumb]…  ›     the thumbnail strip doubles as a
//	  [Info] [Zoom]   [Meet Marmalade]  preview of every cat up for adoption
//
// The back of the card is the cat's details list; "Meet <name>" pushes the
// full profile. The tap zones are pointer shortcuts only — every action also
// has a named button, which is what screen readers use.
func galleryScreen(ctx *core.Context, cats []Cat, requested core.State[map[string]bool]) core.View {
	th := ctx.Theme()

	items := make([]comps.GalleryItem, len(cats))
	for i, c := range cats {
		items[i] = comps.GalleryItem{
			// Keyed by ID so a future filter (e.g. "kittens only") that
			// removes cards keeps each remaining card's node and animation.
			Key:         c.ID,
			Src:         photoURL(c.Photo),
			Alt:         c.Alt,
			Title:       c.Name,
			Subtitle:    c.subtitle(),
			Description: c.Short,
			Details:     c.details(),
			ActionLabel: "Meet " + c.Name,
		}
	}

	return comps.Screen{
		// Scroll so a short viewport (landscape, small phones) can still
		// reach the buttons under the stage. The Gallery's only inner scroll
		// is the horizontal thumbnail strip, which doesn't fight a vertical
		// screen scroll for the same drag.
		Scroll: true,
		Gap:    float64(th.Spacing.SM),
		// Screen paints no background of its own; stating the theme's here
		// also paints the safe-area inset, so the status-bar strip is dark
		// on the natives too.
		Style: []core.StyleProp{core.BackgroundColor(th.Colors.Background)},
		Children: []core.View{
			comps.AppBar{
				Title:    "Cat Adoption",
				Subtitle: fmt.Sprintf("%d cats looking for a home", len(cats)),
			},
			core.Text("Tap a photo for info · tap its edges to browse",
				core.UseStyle(th.Typography.Caption),
				core.TextColor(th.Colors.TextSecondary),
			),
			// Rendered unconditionally on every pass: the Gallery owns two
			// hook slots (its state and the autoplay timer).
			comps.Gallery{
				Items:  items,
				Height: galleryHeight,
				// Wrap so browsing feels endless; with only nine cats,
				// stopping dead at either end reads as "that's all" too soon.
				Wrap: true,
				// Thumbnails rather than dots: the strip is a contact sheet of
				// every cat, which is the point of an adoption page.
				Indicator: comps.GalleryIndicatorThumbnails,
				// Fullscreen zoom opens the uncropped photo in a Lightbox.
				// The card crops to fill; adopters want to see the whole cat.
				Zoom: comps.GalleryZoomFullscreen,
				// Three lines fit every blurb in the catalogue without an
				// ellipsis; the back of the card shows the facts anyway.
				DescriptionLines: 3,
				Labels:           comps.GalleryLabels{Gallery: "Cats for adoption", Details: "Info", Front: "Show photo"},
				OnAction: func(_ comps.GalleryItem, i int) {
					// i indexes Items, which was built 1:1 from cats.
					c := cats[i]
					core.Push(ctx, func(frame *core.Context) core.View {
						return profileScreen(frame, c, requested)
					})
				},
			},
		},
	}
}

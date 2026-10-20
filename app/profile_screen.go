package app

import (
	"fmt"

	"github.com/rohanthewiz/grmob/comps"
	"github.com/rohanthewiz/grmob/core"
)

// profileHeight is the hero photo's height. The photo is cropped to fill
// (ContentModeFill), so this is a framing choice, not the image's own size.
const profileHeight = 320

// profileScreen is the pushed route for one cat: hero photo, story, facts and
// an adoption-request button.
//
//	┌ AppBar  ‹  Marmalade ───────┐
//	│ [ hero photo, cropped ]      │
//	│ Marmalade                    │
//	│ Ginger tabby · 4 yrs · Male  │
//	│ (Banner once requested)      │
//	│ story paragraph…             │
//	│ Age ........... 4 yrs        │
//	│ …                            │
//	│ [ Ask to adopt Marmalade ]   │
//	│ Photo: Wikimedia Commons…    │
//	└──────────────────────────────┘
//
// The route declares no hooks of its own. The request flag lives in the App's
// requested map (see App), because this frame and its state are discarded on
// Pop, and a request should still show the next time the profile is opened.
func profileScreen(ctx *core.Context, c Cat, requested core.State[map[string]bool]) core.View {
	th := ctx.Theme()
	asked := requested.Get()[c.ID]

	ask := func() {
		// Copy, change the copy, Set it: mutating the map in place would
		// leave the slot holding the same value, and the request would not
		// re-render reliably.
		next := make(map[string]bool, len(requested.Get())+1)
		for k, v := range requested.Get() {
			next[k] = v
		}
		next[c.ID] = true
		requested.Set(next)
	}

	label := "Ask to adopt " + c.Name
	if asked {
		label = "Request sent"
	}

	return comps.Screen{
		Scroll: true,
		Gap:    float64(th.Spacing.MD),
		Style:  []core.StyleProp{core.BackgroundColor(th.Colors.Background)},
		Children: []core.View{
			// The back button is automatic: AppBar draws it whenever
			// core.CanPop is true, and it pops by default.
			comps.AppBar{Title: c.Name},
			// Placed straight in the Screen's column, which already carries
			// the theme's inset, so the photo lines up with the text below.
			core.ImageWithMode(photoURL(c.Photo), core.ContentModeFill,
				core.Width("100%"), core.Height(fmt.Sprintf("%dpx", profileHeight)),
				core.BorderRadius(th.Components.Card.BorderRadius),
				core.BackgroundColor(th.Colors.Surface), // placeholder while it loads
				core.AccessibilityLabel(c.Alt),
			),
			core.Column(
				// The theme insets every Column; this one is nested inside
				// the Screen's already padded column, so clear it or the
				// name would sit indented from the photo above.
				core.Padding(0),
				core.Gap(float64(th.Spacing.XS)),
				core.Text(c.Name, core.UseStyle(th.Typography.Title),
					core.AccessibilityRole(core.RoleHeading)),
				core.Text(c.Coat+" · "+c.subtitle(),
					core.UseStyle(th.Typography.Body),
					core.TextColor(th.Colors.TextSecondary)),
			),
			// A conditional view is fine (Banner holds no hooks); it is
			// conditional hooks that shift slots.
			core.If(asked, comps.Banner{
				Variant: comps.VariantSuccess,
				Text:    "Thanks! The shelter will contact you about " + c.Name + ".",
			}),
			core.Text(c.Story, core.UseStyle(th.Typography.Body)),
			comps.KeyValueList{Label: "About " + c.Name, Rows: c.details(), Dividers: true},
			comps.Button{
				Label:     label,
				OnTap:     ask,
				FullWidth: true,
				Disabled:  asked,
			},
			core.Text("Photo: "+photoCredit(c.Photo),
				core.UseStyle(th.Typography.Caption),
				core.TextColor(th.Colors.TextSecondary)),
		},
	}
}

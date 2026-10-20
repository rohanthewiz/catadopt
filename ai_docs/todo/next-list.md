# Next list

catadopt's follow-ups in one place. Sessions edit this file in place: they add,
correct, move and close items here, and their session doc's `## Next` section
just summarises the change (`Closed: … Raised: …`). Nothing is copied forward
from one doc to the next. Seeded on 2026-10-19 from the two session docs then in
`ai_docs/claude_sessions/` by `/next-list seed`.

## Conventions

- **IDs** (`N-001`…) are permanent and never reused, even after an item closes.
- **`raised`** is the stem of the session doc an item first appeared in.
- **Age** is computed, never stored: the number of session docs since `raised`
  (an item raised in the newest doc has age 0). `/next-list` works it out.
- **Value** is the payoff of doing the item, not the effort:
  - `high`: something is being worked around today, or a second independent
    consumer has arrived.
  - `medium`: it blocks one named thing, or it is a visible defect nobody has
    to route around yet.
  - `low`: a gap nobody has bumped into, or contingent on something that does
    not exist.
- **Sections:**
  - **Open** is the next build work: code to write, a defect to fix, a
    decision to make.
  - **Validate** is the next testing: run, look, hear or measure, or write or
    repair a test. No product change is planned unless the check finds a
    defect.
  - **Roadmap** is wanted, but later.
  - **Non-goals** is what we will likely never do.
- Nothing leaves Open, Validate or Roadmap without a line in another section.
- Open, Validate and Roadmap stay in ID order.

**Next ID:** N-011

## Open

- **N-001** · raised `2026-1019-2251-cat-adoption-gallery-app` · value low
  Real listings: replace the sample `catalogue()` with a shelter feed, or a
  bytdb store opened lazily per the GrMob persistence pattern.

- **N-002** · raised `2026-1019-2251-cat-adoption-gallery-app` · value medium
  Send adoption requests somewhere: email, a form backend, or the shelter's
  API. They are persisted on the device now (bytdb, `app/store.go`, table
  `adoption_requests` with `requested_at`), but nothing leaves the phone, and
  the browser build still keeps them in memory only (no data directory).

- **N-003** · raised `2026-1019-2251-cat-adoption-gallery-app` · value medium
  Per-photo author and licence on the profile. The profile prints the Commons
  file-page URL (`photoCredit`) as caption text. That page carries the author
  and licence, but the profile doesn't state them, and the URL isn't tappable.
  The session docs called it a "link".

- **N-004** · raised `2026-1019-2251-cat-adoption-gallery-app` · value low
  Photos are hotlinked from upload.wikimedia.org. Consider bundling or caching
  them for offline use on the natives. Hotlinking has already bitten once:
  Wikimedia answered okhttp's default agent with 403 on Android, which was
  fixed upstream as grmob N-116 in v0.7.1.

- **N-005** · raised `2026-1019-2251-cat-adoption-gallery-app` · value low
  Filters (kittens, good with kids or dogs). Gallery items are keyed by ID, so
  a filtered list keeps each card's node. `Cat.Age` and `Cat.GoodWith` are
  display strings ("4 yrs", "Kids, other cats, cat-friendly dogs"), so a
  filter first needs structured fields, for example age in months and
  per-companion flags, rather than matching on text.

## Validate

Items whose remaining work is purely testing: run, look, hear, measure, or
write or repair a test. No product change is planned unless a check finds a
defect, and if one does, the defect is raised as a new Open item.

- **N-006** · raised `2026-1019-2251-cat-adoption-gallery-app` · value medium
  Android, still unseen: the Fold6's inner (unfolded) screen; the profile
  screen; the edge tap zones (grmob's N-111 has never tapped one on the
  natives), long press and the fullscreen Lightbox. Photos, the strip, the
  middle tap and the full-height background are seen on the cover screen.

- **N-007** · raised `2026-1019-2251-cat-adoption-gallery-app` · value medium
  iOS, seen in the simulator only and nothing tapped (simctl cannot tap): the
  flip, the edge zones, Zoom's Lightbox and the profile screen with its back
  control. Then a real device.

- **N-008** · raised `2026-1019-2251-cat-adoption-gallery-app` · value low
  Screen readers: the card names ("3 of 9, Juniper"), the Info and Zoom
  controls, and the profile heading and Banner, on TalkBack and VoiceOver.

- **N-010** · raised `2026-1019-2359-persist-adoption-requests-bytdb` · value medium
  Adoption-request persistence on a device: ask to adopt, force-stop the app,
  relaunch, and expect "Request sent" on that cat's profile, on Android and
  on iOS. Covered by `TestAdoptionSurvivesRelaunch` in Go, but not yet seen
  through the shells' `SetDataDir`.

## Roadmap

Wanted, but deferred. This section was seeded empty; move items here from Open
when they are parked rather than next.

## Non-goals

- **N-009** · declined `2026-1019-2251-cat-adoption-gallery-app` — Working
  around Dark Reader beyond `<meta name="darkreader-lock">`. It is a per-user
  extension setting: turn it off for localhost. In filter mode it ignores the
  lock and its image filter hides Gallery's captions, whatever the theme.

## Closed

Earlier closures live in the session docs. Two parts of the first doc's
Validate item didn't carry into N-006 or N-007, for these reasons:

- The build and run on Android and iOS was done in
  `2026-1019-2347-warm-theme-android-and-ios`.
- The "dark status-bar strip" check was overtaken when the app moved from
  DarkTheme to AmberTheme (2179d13).

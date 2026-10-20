# Warm theme, Android on a Fold6, iOS in the simulator

Session: `ada40c4b-9595-4b2b-bcbc-e2e6ecae2cf7`
**Date:** 2026-10-19 23:47 · **Repo:** https://github.com/rohanthewiz/catadopt (79aa537 → e53cf8d, plus this doc) · **grmob:** v0.7.0 → v0.7.2

## Goal

Follows `2026-1019-2251-cat-adoption-gallery-app`. The user asked to:

1. switch to a warm theme, commit, and build for Android and install on
   their device;
2. fix upstream in grmob what the device showed, then release and update;
3. build for iOS and run in the simulator, then fix that upstream and
   update too.

## Warm theme (2179d13)

- `appTheme` in `app/app.go` holds the palette in one place; it is now
  `core.AmberTheme`, grmob's bundled warm theme.
- AmberTheme keeps its page (`Background`) and cards white and puts its
  warmth in `Surface` (amber 50, `#FFF8E1`). So both screens, and their
  AppBars, paint `Surface`. A white page with a white card did not read as
  warm. The AppBar's own white fill drew as a stray band inside the
  screen's inset.
- `wasm/index.html`: light colour scheme, `#FFF8E1` behind the first frame,
  and a sand page (`#EFE3C8`) around the phone frame. The `darkreader-lock`
  meta stays, because Dark Reader's image filter hides Gallery captions in
  any theme.
- `TestDarkMode` became `TestAppTheme`, which checks that the safe area
  paints `appTheme.Colors.Surface`.

## Android on the Galaxy Z Fold6 (cb8791d, 79a31e1)

- The first `grmob android` build vendored the shell into `android/` and
  registered gomobile/gobind as tools in `go.mod`.
- Installed with `grmob android -install`. Device SM-F956U1, viewed on the
  cover screen, 968×2376. Taking screenshots needs `screencap -d <id>`,
  because the Fold has two displays (`dumpsys SurfaceFlinger
  --display-id`).
- What showed up, and what fixed it:
  - **Blank photos.** Wikimedia answers `okhttp/4.x` with 403 ("Please set
    a user-agent…"). A local `GrMobApplication.kt` set a descriptive agent
    for Coil. Upstreamed as grmob N-116, in v0.7.1.
  - **White band below the content.** Fixed with `Screen.Fill`.
  - **Thumbnail strip over the ‹ › arrows.** Compose floored the Scroll at
    min-content. Fixed upstream as grmob N-117, in v0.7.1.
- After v0.7.1 and `grmob android -refresh`, the shell's
  `GrMobApplication.kt` is grmob's own. Photos load, the strip scrolls
  between the arrows, and a middle tap flips the card to its info.

## iOS in the simulator (e53cf8d)

- `grmob ios -run` vendored the shell into `ios/` and ran on the booted
  iPhone 17 Pro (iOS 26.5).
- What showed up:
  - **The whole screen about 693pt wide and centred.** A filled photo
    widened its frame. Fixed upstream as grmob N-118.
  - **White below the content.** `Screen{Scroll, Fill}` didn't grow the
    Scroll on SwiftUI. Fixed upstream as grmob N-119.

  Both are in v0.7.2.
- Bisected with two throwaway rebuilds: thumbnails → dots, then Fill off.
  Neither was the cause.
- On v0.7.2 (`ios -refresh`, 35 shell source files committed; build
  products are ignored) the screen fits, the description wraps, the photo
  is cropped, the strip sits between the arrows, and the Surface reaches
  the bottom.

## Notes

- grmob's sessions for the upstream half:
  `2026-1019-2246-gallery-tapzone-zindex` (v0.7.0),
  `2026-1019-2323-android-image-agent-and-scroll-strip-floor` (v0.7.1),
  `2026-1019-2343-ios-filled-image-frame-and-screen-fill-scroll` (v0.7.2).
- Upstream fixes were tested here before release through a temporary
  `replace` plus `-refresh`. catadopt was restored to its commit each
  time.
- Android emulators on this machine: `Medium_Phone_API_36.1` and
  `GrMob_Foldable`. The phone emulator stood in for the Fold6 once it was
  unplugged.

## Next

- Real listings: replace the sample `catalogue()` with a shelter feed, or a
  bytdb store opened lazily per the GrMob persistence pattern.
- Persist adoption requests (bytdb), and send them somewhere: email, a form
  backend, or the shelter's API. Today they live in memory only.
- Per-photo author and licence on the profile. It links to each Commons
  file page now, but doesn't state them.
- Photos are hotlinked from upload.wikimedia.org. Consider bundling or
  caching them for offline use on the natives.
- Filters (kittens, good with kids or dogs). Gallery items are keyed by ID,
  so a filtered list keeps each card's node.
- Non-goal: working around Dark Reader beyond `darkreader-lock`. It is a
  per-user extension setting (turn it off for localhost).

**Validate**

- Android, still unseen: the Fold6's inner (unfolded) screen; the
  profile screen; the edge tap zones, long press and the fullscreen
  Lightbox. Photos, the strip, the middle tap and the full-height
  background are seen on the cover screen.
- iOS, seen in the simulator only and nothing tapped (simctl cannot tap):
  the flip, the edge zones, Zoom's Lightbox and the profile screen with
  its back control. Then a real device.
- Screen readers: the card names ("3 of 9, Juniper"), the Info and Zoom
  controls, and the profile heading and Banner, on TalkBack and VoiceOver.

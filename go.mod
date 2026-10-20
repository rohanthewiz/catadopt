module catadopt

go 1.26.10

require github.com/rohanthewiz/grmob v0.7.2

require (
	github.com/rohanthewiz/element v0.7.0 // indirect
	github.com/rohanthewiz/rweb v0.1.28 // indirect
	github.com/rohanthewiz/serr v1.4.0 // indirect
	golang.org/x/mobile v0.0.0-20251021151156-188f512ec823 // indirect
	golang.org/x/mod v0.29.0 // indirect
	golang.org/x/sync v0.17.0 // indirect
	golang.org/x/tools v0.38.0 // indirect
)

tool (
	github.com/rohanthewiz/grmob/serve
	golang.org/x/mobile/cmd/gobind
	golang.org/x/mobile/cmd/gomobile
)

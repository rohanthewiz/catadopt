package app

import (
	"fmt"

	"github.com/rohanthewiz/grmob/comps"
)

// Cat is one adoptable cat. The fields are plain data so the same record
// feeds both the gallery card (short description, details on the back) and
// the profile screen (full story, adoption request).
type Cat struct {
	ID    string // stable key: gallery card key and adoption-request map key
	Name  string
	Age   string
	Sex   string
	Coat  string // breed or coat pattern, shown as the card's subtitle
	Short string // one or two sentences for the front of the card
	Story string // the longer write-up on the profile screen

	GoodWith   string
	Vaccinated bool
	Fixed      bool // spayed or neutered
	Fee        int  // adoption fee in USD

	// Photo is a Wikimedia Commons file path ("3/3a/Cat03.jpg"): the hash
	// directory plus the file name, which is how upload.wikimedia.org lays
	// out originals and thumbnails. Keeping the path rather than a full URL
	// lets photoURL pick the thumbnail size in one place.
	Photo string
	Alt   string
}

// photoWidth is the Wikimedia thumbnail width requested. The app column is
// at most ~430pt wide, so 960px covers a 2x display without pulling the
// multi-megabyte originals. Wikimedia only serves a fixed set of thumbnail
// widths (250, 330, 500, 960, 1280, …); arbitrary widths 404.
const photoWidth = 960

// photoURL builds the thumbnail URL for a Commons file path:
//
//	https://upload.wikimedia.org/wikipedia/commons/thumb/3/3a/Cat03.jpg/960px-Cat03.jpg
//	                                                     └ path ────────┘ └ size + name ┘
func photoURL(path string) string {
	return fmt.Sprintf("https://upload.wikimedia.org/wikipedia/commons/thumb/%s/%dpx-%s",
		path, photoWidth, fileName(path))
}

// photoCredit points at the Commons file page, which carries the author and
// licence each photo requires us to credit.
func photoCredit(path string) string {
	return "https://commons.wikimedia.org/wiki/File:" + fileName(path)
}

func fileName(path string) string {
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '/' {
			return path[i+1:]
		}
	}
	return path
}

// yesNo renders a bool for the details list.
func yesNo(b bool) string {
	if b {
		return "Yes"
	}
	return "Not yet"
}

// details are the facts on the back of the gallery card and on the profile
// screen. One function so the two views cannot disagree.
func (c Cat) details() []comps.KeyValue {
	return []comps.KeyValue{
		{Key: "Age", Value: c.Age},
		{Key: "Sex", Value: c.Sex},
		{Key: "Coat", Value: c.Coat},
		{Key: "Good with", Value: c.GoodWith},
		{Key: "Vaccinated", Value: yesNo(c.Vaccinated)},
		{Key: "Spayed / neutered", Value: yesNo(c.Fixed)},
		{Key: "Adoption fee", Value: fmt.Sprintf("$%d", c.Fee)},
	}
}

// subtitle is the line under the name on the card: "4 yrs · Male".
func (c Cat) subtitle() string { return c.Age + " · " + c.Sex }

// catalogue is sample data. The photos are real Commons images; the names,
// ages and personalities are written to match each picture but are invented.
// A real shelter would load this from its listings API or a bytdb store.
func catalogue() []Cat {
	return []Cat{
		{
			ID: "marmalade", Name: "Marmalade", Age: "4 yrs", Sex: "Male", Coat: "Ginger tabby",
			Short:    "A gentle ginger gentleman who greets every visitor at the door and lives for a sunny windowsill.",
			Story:    "Marmalade came to us when his elderly owner moved into care. He is calm, polite and a little bit nosy: expect him to inspect your shopping bags. He settles quickly in a quiet home and will happily share the sofa.",
			GoodWith: "Adults, calm kids", Vaccinated: true, Fixed: true, Fee: 85,
			Photo: "3/3a/Cat03.jpg",
			Alt:   "A ginger tabby cat with amber eyes looking up at the camera",
		},
		{
			ID: "pepper", Name: "Pepper", Age: "3 yrs", Sex: "Female", Coat: "Brown mackerel tabby",
			Short:    "Master of the afternoon stretch. Calm, independent and happiest lounging on a warm ledge.",
			Story:    "Pepper is the definition of relaxed. She enjoys company but doesn't demand it, which makes her a great fit for someone who works from home. She has a soft spot for chin scratches and cardboard boxes.",
			GoodWith: "Adults, other cats", Vaccinated: true, Fixed: true, Fee: 85,
			Photo: "1/15/Cat_August_2010-4.jpg",
			Alt:   "A brown tabby cat stretched out on a white stone step",
		},
		{
			ID: "juniper", Name: "Juniper", Age: "2 yrs", Sex: "Female", Coat: "Brown tabby",
			Short:    "Alert and curious. She'll supervise your every move from the highest perch in the room.",
			Story:    "Juniper is a watcher. Give her a tall cat tree near a window and she'll be content for hours. She warms up over a few days and then follows her favourite person from room to room.",
			GoodWith: "Adults, older kids", Vaccinated: true, Fixed: true, Fee: 95,
			Photo: "4/4d/Cat_November_2010-1a.jpg",
			Alt:   "A brown tabby cat with green eyes sitting upright on a wall",
		},
		{
			ID: "frost", Name: "Frost", Age: "1 yr", Sex: "Male", Coat: "Brown tabby",
			Short:    "An adventurous young tabby who isn't fazed by anything, not even snow. Needs daily playtime.",
			Story:    "Frost has energy to spare. He loves feather wands, ping-pong balls and racing down hallways. A home with another playful cat, or a family ready for lots of games, would suit him best.",
			GoodWith: "Kids, other cats, cat-friendly dogs", Vaccinated: true, Fixed: false, Fee: 110,
			Photo: "b/b6/Felis_catus-cat_on_snow.jpg",
			Alt:   "A young brown tabby cat walking across snow and looking back over its shoulder",
		},
		{
			ID: "saffron", Name: "Saffron", Age: "5 yrs", Sex: "Male", Coat: "Abyssinian",
			Short:    "Sleek, athletic and chatty. An Abyssinian who loves climbing and will follow you everywhere.",
			Story:    "Saffron is busy, bright and very talkative. He loves puzzle feeders and high shelves. He's happiest with people who are home a lot and enjoy an opinionated companion.",
			GoodWith: "Adults, other cats", Vaccinated: true, Fixed: true, Fee: 150,
			Photo: "9/9b/Gustav_chocolate.jpg",
			Alt:   "A slender ruddy Abyssinian cat sitting on a fleece and looking over its shoulder",
		},
		{
			ID: "biscuit", Name: "Biscuit", Age: "6 yrs", Sex: "Female", Coat: "Tabby and white",
			Short:    "Playful paws and big green eyes. Biscuit will high-five you for treats.",
			Story:    "Biscuit learned to high-five in foster care and now offers it to everyone. She is affectionate, food-motivated and still kitten-playful at six. She'd love a home with a sunny spot and a stash of treats.",
			GoodWith: "Adults, kids", Vaccinated: true, Fixed: true, Fee: 75,
			Photo: "b/bb/Kittyply_edit1.jpg",
			Alt:   "A close-up of a tabby and white cat with green eyes reaching a paw toward the camera",
		},
		{
			ID: "rusty", Name: "Rusty", Age: "3 yrs", Sex: "Male", Coat: "Orange and white",
			Short:    "A sweet orange-and-white boy with a steady gaze. Great with older kids.",
			Story:    "Rusty is easy-going and confident. He enjoys a good brush, short play sessions and a long nap afterwards. He has lived with children before and takes noise in his stride.",
			GoodWith: "Kids, other cats", Vaccinated: true, Fixed: true, Fee: 85,
			Photo: "6/68/Orange_tabby_cat_sitting_on_fallen_leaves-Hisashi-01A.jpg",
			Alt:   "An orange and white cat sitting upright on fallen autumn leaves",
		},
		{
			ID: "nugget", Name: "Nugget", Age: "10 wks", Sex: "Male", Coat: "Ginger tabby",
			Short:    "A tiny ginger kitten full of mischief. Litter-trained and ready for a patient home.",
			Story:    "Nugget is all whiskers and curiosity. He's litter-trained, eating solid food and learning his manners. Kittens do best in pairs, so ask us about adopting him with a sibling.",
			GoodWith: "Everyone (with supervision)", Vaccinated: false, Fixed: false, Fee: 125,
			Photo: "a/a5/Red_Kitten_01.jpg",
			Alt:   "A small ginger kitten lying on a wooden deck",
		},
		{
			ID: "skye", Name: "Skye", Age: "2 yrs", Sex: "Female", Coat: "Tabby and white, medium hair",
			Short:    "Striking blue eyes and a soft, fluffy coat. Shy at first, then a devoted lap cat.",
			Story:    "Skye needs a little time to trust new people, but once she does she's a loyal lap cat who purrs at the slightest touch. A quiet, adult household would help her shine.",
			GoodWith: "Adults", Vaccinated: true, Fixed: true, Fee: 95,
			Photo: "c/c7/Tabby_cat_with_blue_eyes-3336579.jpg",
			Alt:   "A fluffy tabby and white cat with bright blue eyes",
		},
	}
}

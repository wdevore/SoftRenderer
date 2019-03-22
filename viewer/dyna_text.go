package viewer

import (
	"github.com/veandco/go-sdl2/sdl"
)

type dChar struct {
	width, height int32
	txt           *Text
}

// DynaText represents dynamically changing text.
// The text is broken down into characters and each
// character is rendered using an existing character atlas.
// This is a base class for dynamic text.
type DynaText struct {
	nFont    *Font
	renderer *sdl.Renderer

	// Char textures
	txtChars []*dChar

	surface *sdl.Surface
	texture *sdl.Texture
	Color   sdl.Color

	bounds sdl.Rect

	X      int32
	Y      int32
	Width  int32
	Height int32
}

// NewDynaText creates a Text object.
func NewDynaText(font *Font, renderer *sdl.Renderer, color sdl.Color) *DynaText {
	t := new(DynaText)
	t.nFont = font
	t.renderer = renderer
	t.X = 0
	t.Y = 0
	t.Color = color
	t.initialize()

	return t
}

// Initialize sets up Text based on TextPath
func (t *DynaText) initialize() error {
	// Iterate through each character generating a texture.
	chars := "0123456789-:.,<> ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	t.txtChars = make([]*dChar, len(chars))

	i := 0
	for _, c := range chars {
		txtC := NewText(t.nFont, t.renderer)
		txtC.SetText(string(c), t.Color)
		t.txtChars[i] = new(dChar)
		t.txtChars[i].width = txtC.Bounds.W
		t.txtChars[i].height = txtC.Bounds.H
		t.txtChars[i].txt = txtC
		i++
	}

	return nil
}

// Draw renders text
func (t *DynaText) Draw(text string) {
	x := t.X
	for _, c := range text {
		for _, dc := range t.txtChars {
			if string(c) == dc.txt.text {
				dc.txt.Bounds.X = x
				dc.txt.Bounds.Y = t.Y
				dc.txt.Draw()
				x += dc.width
			}
		}
	}
}

// DrawAt renders text at the specified position
func (t *DynaText) DrawAt(x, y int32, text string) {
	for _, c := range text {
		for _, dc := range t.txtChars {
			if string(c) == dc.txt.text {
				dc.txt.Bounds.X = x
				dc.txt.Bounds.Y = y
				dc.txt.Draw()
				x += dc.width
			}
		}
	}
}

// Destroy closes the Text
func (t *DynaText) Destroy() {
	for _, dc := range t.txtChars {
		dc.txt.Destroy()
	}
}

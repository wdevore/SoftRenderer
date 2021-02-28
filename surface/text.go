package surface

import (
	"github.com/veandco/go-sdl2/sdl"
)

// Text represents a text texture for rendering.
type Text struct {
	nFont    *Font
	renderer *sdl.Renderer

	texture *sdl.Texture
	color   sdl.Color
	text    string
	Bounds  sdl.Rect
}

// NewText creates a Text object.
func NewText(font *Font, renderer *sdl.Renderer) *Text {
	t := new(Text)
	t.nFont = font
	t.renderer = renderer
	t.initialize()

	return t
}

// Initialize sets up Text based on TextPath
func (t *Text) initialize() error {
	return nil
}

// SetText builds an SDL texture. Be sure to call Destroy before
// program exit.
func (t *Text) SetText(text string, color sdl.Color) (err error) {
	t.text = text
	t.color = color

	t.Destroy()

	var surface *sdl.Surface

	// First we draw an image to a surface
	surface, err = t.nFont.font.RenderUTF8Solid(text, t.color)
	if err != nil {
		return err
	}

	// Now generate a texture for rendering, using the surface.
	t.texture, err = t.renderer.CreateTextureFromSurface(surface)
	if err != nil {
		t.Destroy()
		return err
	}

	_, _, width, height, qerr := t.texture.Query()
	if qerr != nil {
		t.Destroy()
		return qerr
	}

	t.Bounds = sdl.Rect{X: 0, Y: 0, W: width, H: height}

	// We don't need the surface any longer now that the texture
	// is created.
	surface.Free()

	return nil
}

// Draw renders text
func (t *Text) Draw() {
	t.renderer.Copy(t.texture, nil, &t.Bounds)
}

// DrawAt renders text
func (t *Text) DrawAt(x, y int32) {
	t.Bounds.X = x
	t.Bounds.Y = y
	t.renderer.Copy(t.texture, nil, &t.Bounds)
}

// Destroy closes the Text
func (t *Text) Destroy() error {
	if t.texture != nil {
		return t.texture.Destroy()
	}

	return nil
}

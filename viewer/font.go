package viewer

import "github.com/veandco/go-sdl2/ttf"

// Font wraps the SDL TTF fonts.
type Font struct {
	font     *ttf.Font
	fontPath string
	size     int
}

// NewFont creates a Font object:
// Ex NewFont("neuropol x rg.ttf", 16)
func NewFont(fontPath string, size int) (*Font, error) {
	f := new(Font)
	f.fontPath = fontPath
	f.size = size

	err := f.initialize()
	if err != nil {
		return nil, err
	}

	return f, nil
}

// Initialize sets up font based on fontPath
func (f *Font) initialize() (err error) {
	ttf.Init()

	f.font, err = ttf.OpenFont(f.fontPath, f.size)
	if err != nil {
		return err
	}

	return nil
}

// Destroy closes the font
func (f *Font) Destroy() {
	f.font.Close()
	ttf.Quit()
}

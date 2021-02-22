package api

import "github.com/veandco/go-sdl2/sdl"

// ISurface is the graph viewer
type ISurface interface {
	// Open(IHost)
	Open()
	Run()
	Close()
	Quit()
	Configure()
	SetFont(fontPath string, size int) error

	SetDrawColor(color sdl.Color)
	SetPixel(x, y int)
}

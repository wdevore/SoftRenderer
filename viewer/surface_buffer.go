package viewer

import (
	"SoftRenderer/api"
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

// SurfaceBuffer is the GUI and shows the plots and graphs.
// It receives commands for graphing and viewing various graphs.
type SurfaceBuffer struct {
	window   *sdl.Window
	surface  *sdl.Surface
	renderer *sdl.Renderer
	texture  *sdl.Texture

	// drawing buffer
	pixels *image.RGBA
	bounds image.Rectangle

	// mouse
	mx int32
	my int32

	// Debug/testing stuff
	mod int

	running bool

	opened bool

	nFont        *Font
	txtSimStatus *Text
	txtFPSLabel  *Text
	txtLoopLabel *Text
	txtMousePos  *Text
	dynaTxt      *DynaText
}

// NewSurfaceBuffer creates a new viewer and initializes it.
func NewSurfaceBuffer() api.ISurface {
	v := new(SurfaceBuffer)
	v.opened = false
	v.mod = 200
	return v
}

// Open shows the viewer and begins event polling
// (host deuron.IHost)
func (v *SurfaceBuffer) Open() {
	v.initialize()

	v.opened = true
}

// SetFont sets the font based on path and size.
func (v *SurfaceBuffer) SetFont(fontPath string, size int) error {
	var err error
	v.nFont, err = NewFont(fontPath, size)
	return err
}

// filterEvent returns false if it handled the event. Returning false
// prevents the event from being added to the queue.
func (v *SurfaceBuffer) filterEvent(e sdl.Event, userdata interface{}) bool {
	switch t := e.(type) {
	case *sdl.QuitEvent:
		v.running = false
		return false // We handled it. Don't allow it to be added to the queue.
	case *sdl.MouseMotionEvent:
		v.mx = t.X
		v.my = t.Y
		// fmt.Printf("[%d ms] MouseMotion\ttype:%d\tid:%d\tx:%d\ty:%d\txrel:%d\tyrel:%d\n",
		// 	t.Timestamp, t.Type, t.Which, t.X, t.Y, t.XRel, t.YRel)
		return false // We handled it. Don't allow it to be added to the queue.
		// case *sdl.MouseButtonEvent:
		// 	fmt.Printf("[%d ms] MouseButton\ttype:%d\tid:%d\tx:%d\ty:%d\tbutton:%d\tstate:%d\n",
		// 		t.Timestamp, t.Type, t.Which, t.X, t.Y, t.Button, t.State)
		// case *sdl.MouseWheelEvent:
		// 	fmt.Printf("[%d ms] MouseWheel\ttype:%d\tid:%d\tx:%d\ty:%d\n",
		// 		t.Timestamp, t.Type, t.Which, t.X, t.Y)
	case *sdl.KeyboardEvent:
		if t.State == sdl.PRESSED {
			switch t.Keysym.Scancode {
			case sdl.SCANCODE_ESCAPE:
				v.running = false
				// case 's':
				// 	// Start sim
				// 	// simStatus = "Starting"
				// case 'o':
				// 	// Stop sim
				// 	// simStatus = "Stopping"
				// case 'p':
				// 	// Pause sim
				// 	// simStatus = "Pausing"
				// case 'e':
				// Step sim
				// simStatus = "Stepping"
			}
		}
		// fmt.Printf("[%d ms] Keyboard\ttype:%d\tsym:%c\tmodifiers:%d\tstate:%d\trepeat:%d\n",
		// 	t.Timestamp, t.Type, t.Keysym.Sym, t.Keysym.Mod, t.State, t.Repeat)
		return false
	}

	return true
}

// Run starts the polling event loop. This must run on
// the main thread.
func (v *SurfaceBuffer) Run() {
	// log.Println("Starting viewer polling")
	v.running = true
	// var simStatus = ""
	var frameStart time.Time
	var elapsedTime float64
	var loopTime float64

	sleepDelay := 0.0

	// Get a reference to SDL's internal keyboard state. It is updated
	// during sdl.PollEvent()
	keyState := sdl.GetKeyboardState()
	c := color.RGBA{}
	c.B = 127
	c.A = 255
	// rect := sdl.Rect{X: 0, Y: 0, W: 100, H: 100}

	sdl.SetEventFilterFunc(v.filterEvent, nil)

	for v.running {
		frameStart = time.Now()

		sdl.PumpEvents()

		if keyState[sdl.SCANCODE_Z] != 0 {
			v.mod--
		}
		if keyState[sdl.SCANCODE_X] != 0 {
			v.mod++
		}

		v.clearDisplay()

		// This full loop takes about 20ms for an 800x800 dimension.
		for y := 0; y < width; y++ {
			for x := 0; x < height; x++ {
				// v.SetPixel(x, y, sdl.Color{R: uint8(x % mod), G: uint8(y % mod), B: 127, A: 255})
				c.R = uint8(x % v.mod)
				c.G = uint8(y % v.mod)
				v.pixels.SetRGBA(x, y, c)
			}
		}

		// v.texture.Update(nil, v.pixels, v.pixelPitch)
		// This takes on average 5-7ms
		v.texture.Update(nil, v.pixels.Pix, v.pixels.Stride)
		v.renderer.Copy(v.texture, nil, nil)

		// v.renderer.SetDrawColor(255, 127, 127, 255)
		// v.renderer.FillRect(&rect)

		v.txtFPSLabel.DrawAt(10, 10)
		f := fmt.Sprintf("%2.2f", 1.0/elapsedTime*1000.0)
		v.dynaTxt.DrawAt(v.txtFPSLabel.Bounds.W+10, 10, f)

		// v.mx, v.my, _ = sdl.GetMouseState()
		v.txtMousePos.DrawAt(10, 25)
		f = fmt.Sprintf("<%d, %d>", v.mx, v.my)
		v.dynaTxt.DrawAt(v.txtMousePos.Bounds.W+10, 25, f)

		v.txtLoopLabel.DrawAt(10, 40)
		f = fmt.Sprintf("%2.2f", loopTime)
		v.dynaTxt.DrawAt(v.txtLoopLabel.Bounds.W+10, 40, f)

		v.renderer.Present()

		// time.Sleep(time.Millisecond * 5)
		loopTime = float64(time.Since(frameStart).Nanoseconds() / 1000000.0)
		// elapsedTime = float64(time.Since(frameStart).Seconds())

		sleepDelay = math.Floor(framePeriod - loopTime)
		// fmt.Printf("%3.5f ,%3.5f, %3.5f \n", framePeriod, elapsedTime, sleepDelay)
		if sleepDelay > 0 {
			sdl.Delay(uint32(sleepDelay))
			elapsedTime = framePeriod
		} else {
			elapsedTime = loopTime
		}
	}
}

// Quit stops the gui from running, effectively shutting it down.
func (v *SurfaceBuffer) Quit() {
	v.running = false
}

// Close closes the viewer.
// Be sure to setup a "defer x.Close()"
func (v *SurfaceBuffer) Close() {
	if !v.opened {
		return
	}
	var err error

	if v.nFont == nil {
		return
	}
	v.nFont.Destroy()

	v.txtFPSLabel.Destroy()
	v.txtMousePos.Destroy()
	v.dynaTxt.Destroy()

	log.Println("Destroying texture")
	err = v.texture.Destroy()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Destroying renderer")
	v.renderer.Destroy()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Destroying window")
	err = v.window.Destroy()
	if err != nil {
		log.Fatal(err)
	}

	sdl.Quit()

	if err != nil {
		log.Fatal(err)
	}
}

func (v *SurfaceBuffer) initialize() {
	var err error

	err = sdl.Init(sdl.INIT_TIMER | sdl.INIT_VIDEO | sdl.INIT_EVENTS)
	if err != nil {
		panic(err)
	}

	v.window, err = sdl.CreateWindow("Soft renderer", windowPosX, windowPosY,
		width, height, sdl.WINDOW_SHOWN)

	if err != nil {
		panic(err)
	}

	// Using GetSurface requires using window.UpdateSurface() rather than renderer.Present.
	// v.surface, err = v.window.GetSurface()
	// if err != nil {
	// 	panic(err)
	// }
	// v.renderer, err = sdl.CreateSoftwareRenderer(v.surface)
	// OR create renderer manually
	v.renderer, err = sdl.CreateRenderer(v.window, -1, sdl.RENDERER_SOFTWARE)
	if err != nil {
		panic(err)
	}

	v.texture, err = v.renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, width, height)
	if err != nil {
		panic(err)
	}

	v.bounds = image.Rect(0, 0, width, height)
	v.pixels = image.NewRGBA(v.bounds)
}

// Configure view with draw objects
func (v *SurfaceBuffer) Configure() {
	// rect := sdl.Rect{X: 0, Y: 0, W: 200, H: 200}
	// v.renderer.SetDrawColor(255, 127, 0, 255)
	// v.renderer.FillRect(&rect)

	v.txtSimStatus = NewText(v.nFont, v.renderer)
	err := v.txtSimStatus.SetText("Sim Status: ", sdl.Color{R: 0, G: 0, B: 255, A: 255})
	if err != nil {
		v.Close()
		panic(err)
	}

	v.txtFPSLabel = NewText(v.nFont, v.renderer)
	err = v.txtFPSLabel.SetText("FPS: ", sdl.Color{R: 200, G: 200, B: 200, A: 255})
	if err != nil {
		v.Close()
		panic(err)
	}

	v.txtMousePos = NewText(v.nFont, v.renderer)
	err = v.txtMousePos.SetText("Mouse: ", sdl.Color{R: 255, G: 127, B: 0, A: 255})
	if err != nil {
		v.Close()
		panic(err)
	}

	v.txtLoopLabel = NewText(v.nFont, v.renderer)
	err = v.txtLoopLabel.SetText("Loop: ", sdl.Color{R: 255, G: 127, B: 0, A: 255})
	if err != nil {
		v.Close()
		panic(err)
	}

	v.dynaTxt = NewDynaText(v.nFont, v.renderer, sdl.Color{R: 255, G: 255, B: 255, A: 255})
}

func (v *SurfaceBuffer) clearDisplay() {
	v.renderer.SetDrawColor(255, 127, 127, 255)
	v.renderer.Clear()
	v.window.UpdateSurface()
}

// SetDrawColor --
func (v *SurfaceBuffer) SetDrawColor(color sdl.Color) {
	v.renderer.SetDrawColor(color.R, color.G, color.B, color.A)
}

// SetPixel --
func (v *SurfaceBuffer) SetPixel(x, y int) {

}

package surface

import (
	"SoftRenderer/api"
	graphics "SoftRenderer/graphcs"
	"SoftRenderer/renderer"
	"fmt"
	"image/color"
	"log"
	"math"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

// WindowSurface is the GUI and shows the plots and graphs.
// It receives commands for graphing and viewing various graphs.
type WindowSurface struct {
	window   *sdl.Window
	surface  *sdl.Surface
	renderer *sdl.Renderer
	texture  *sdl.Texture

	rasterBuffer api.IRasterBuffer

	// mouse
	mx int32
	my int32

	// Debug/testing stuff
	mod  int
	dir  int
	dir2 int
	dir3 int
	xx   int
	xx2  int
	xx3  int

	running bool
	animate bool
	step    bool

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
	o := new(WindowSurface)
	o.opened = false
	o.animate = true
	o.step = false
	o.mod = 200
	o.xx = 75 // -41 //75 // x2
	o.dir = 1
	o.xx2 = 0 //-29 //0 // x1
	o.dir2 = 1
	o.xx3 = 100 //28 //100 // y1
	o.dir3 = 1
	//x1  -29 y1  8 x2  -41
	return o
}

func (ws *WindowSurface) initialize() {
	var err error

	err = sdl.Init(sdl.INIT_TIMER | sdl.INIT_VIDEO | sdl.INIT_EVENTS)
	if err != nil {
		panic(err)
	}

	ws.window, err = sdl.CreateWindow("Soft renderer", windowPosX, windowPosY,
		width, height, sdl.WINDOW_SHOWN)

	if err != nil {
		panic(err)
	}

	// Using GetSurface requires using window.UpdateSurface() rather than renderer.Present.
	// ws.surface, err = ws.window.GetSurface()
	// if err != nil {
	// 	panic(err)
	// }
	// ws.renderer, err = sdl.CreateSoftwareRenderer(ws.surface)
	// OR create renderer manually
	ws.renderer, err = sdl.CreateRenderer(ws.window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}

	ws.texture, err = ws.renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, width, height)
	if err != nil {
		panic(err)
	}

	// ws.renderer.SetDrawBlendMode(sdl.BLENDMODE_BLEND)

	ws.rasterBuffer = renderer.NewRasterBuffer(width, height)
	// ws.rasterBuffer.EnableAlphaBlending(true)
}

// Configure view with draw objects
func (ws *WindowSurface) Configure() {
	ws.txtSimStatus = NewText(ws.nFont, ws.renderer)
	err := ws.txtSimStatus.SetText("Sim Status: ", sdl.Color{R: 0, G: 0, B: 255, A: 255})
	if err != nil {
		ws.Close()
		panic(err)
	}

	ws.txtFPSLabel = NewText(ws.nFont, ws.renderer)
	err = ws.txtFPSLabel.SetText("FPS: ", sdl.Color{R: 200, G: 200, B: 200, A: 255})
	if err != nil {
		ws.Close()
		panic(err)
	}

	ws.txtMousePos = NewText(ws.nFont, ws.renderer)
	err = ws.txtMousePos.SetText("Mouse: ", sdl.Color{R: 200, G: 200, B: 200, A: 255})
	if err != nil {
		ws.Close()
		panic(err)
	}

	ws.txtLoopLabel = NewText(ws.nFont, ws.renderer)
	err = ws.txtLoopLabel.SetText("Loop: ", sdl.Color{R: 200, G: 200, B: 200, A: 255})
	if err != nil {
		ws.Close()
		panic(err)
	}

	ws.dynaTxt = NewDynaText(ws.nFont, ws.renderer, sdl.Color{R: 255, G: 255, B: 255, A: 255})
}

// Open shows the viewer and begins event polling
// (host deuron.IHost)
func (ws *WindowSurface) Open() {
	ws.initialize()

	ws.opened = true
}

// SetFont sets the font based on path and size.
func (ws *WindowSurface) SetFont(fontPath string, size int) error {
	var err error
	ws.nFont, err = NewFont(fontPath, size)
	return err
}

// filterEvent returns false if it handled the event. Returning false
// prevents the event from being added to the queue.
func (ws *WindowSurface) filterEvent(e sdl.Event, userdata interface{}) bool {
	switch t := e.(type) {
	case *sdl.QuitEvent:
		ws.running = false
		return false // We handled it. Don't allow it to be added to the queue.
	case *sdl.MouseMotionEvent:
		ws.mx = t.X
		ws.my = t.Y
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
				ws.running = false
			case sdl.SCANCODE_A:
				ws.animate = !ws.animate
			case sdl.SCANCODE_S:
				ws.step = true
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
func (ws *WindowSurface) Run() {
	// log.Println("Starting viewer polling")
	ws.running = true
	// var simStatus = ""
	var frameStart time.Time
	var elapsedTime float64
	var loopTime float64

	sleepDelay := 0.0

	// Get a reference to SDL's internal keyboard state. It is updated
	// during sdl.PollEvent()
	keyState := sdl.GetKeyboardState()

	rasterizer := renderer.NewBresenHamRasterizer()

	sdl.SetEventFilterFunc(ws.filterEvent, nil)

	for ws.running {
		frameStart = time.Now()

		sdl.PumpEvents()

		if keyState[sdl.SCANCODE_Z] != 0 {
			ws.mod--
		}
		if keyState[sdl.SCANCODE_X] != 0 {
			ws.mod++
		}

		ws.clearDisplay()

		ws.render(rasterizer)

		// This takes on average 5-7ms
		// ws.texture.Update(nil, ws.pixels.Pix, ws.pixels.Stride)
		ws.texture.Update(nil, ws.rasterBuffer.Pixels().Pix, ws.rasterBuffer.Pixels().Stride)
		ws.renderer.Copy(ws.texture, nil, nil)

		ws.txtFPSLabel.DrawAt(10, 10)
		f := fmt.Sprintf("%2.2f", 1.0/elapsedTime*1000.0)
		ws.dynaTxt.DrawAt(ws.txtFPSLabel.Bounds.W+10, 10, f)

		// ws.mx, ws.my, _ = sdl.GetMouseState()
		ws.txtMousePos.DrawAt(10, 25)
		f = fmt.Sprintf("<%d, %d>", ws.mx, ws.my)
		ws.dynaTxt.DrawAt(ws.txtMousePos.Bounds.W+10, 25, f)

		ws.txtLoopLabel.DrawAt(10, 40)
		f = fmt.Sprintf("%2.2f", loopTime)
		ws.dynaTxt.DrawAt(ws.txtLoopLabel.Bounds.W+10, 40, f)

		ws.renderer.Present()

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

func (ws *WindowSurface) render(rasterizer api.IRasterizer) {
	// c := color.RGBA{R: 255, G: 127, B: 0, A: 255}
	// This full loop takes about 20ms for an 800x800 dimension.
	// for y := 0; y < height; y++ {
	// 	for x := 0; x < width; x++ {
	// 		c.R = uint8(x % ws.mod)
	// 		c.G = uint8(y % ws.mod)
	// 		ws.rasterBuffer.SetPixelColor(c)
	// 		ws.rasterBuffer.SetPixel(x, y, 0.0)
	// 	}
	// }

	x := 0
	y := 0
	left := false
	rasterizer.DrawLineAmmeraal(ws.rasterBuffer, left, x, y, x+100, y+25) // blue dx>0

	x = 0
	y = 0
	rasterizer.DrawLineAmmeraal(ws.rasterBuffer, left, x, y+25, x+100, y) // blue dx>0

	x = 50
	y = 50
	down := true
	rasterizer.DrawLineAmmeraal(ws.rasterBuffer, down, x, y+50, x+50, y+150) // red
	x = 50
	y = 50
	rasterizer.DrawLineAmmeraal(ws.rasterBuffer, down, x+50, y+50, x, y+150) // red

	// Horizontal
	x = 100
	y = 5
	left = false
	rasterizer.DrawLineAmmeraal(ws.rasterBuffer, left, x, y, x+100, y) // blue
	x = 100
	y = 10
	left = true
	rasterizer.DrawLineAmmeraal(ws.rasterBuffer, left, x, y, x+100, y) // blue

	// Vertical
	x = 100
	y = 20
	// down = false
	// rasterizer.DrawLineAmmeraal(v, down, x, y+100, x, y) // red
	// Or
	down = true
	rasterizer.DrawLineAmmeraal(ws.rasterBuffer, down, x, y, x, y+100) // red
	x = 110
	y = 20
	down = false
	rasterizer.DrawLineAmmeraal(ws.rasterBuffer, down, x, y, x, y+100) // red

	ws.rasterBuffer.SetPixelColor(color.RGBA{R: 255, G: 255, B: 255, A: 255})

	// Triangle flat-bottom ----------------------------------
	x = 200
	y = 25

	x1 := 0
	y1 := 50
	x2 := 50
	y2 := 50
	x3 := 25
	y3 := 0
	// Make sure Y's are consitent
	// rasterizer.Sort(&x1, &y1, &x2, &y2, &x3, &y3)

	tri := graphics.NewTriangle()

	down = false
	// rasterizer.DrawLineAmmeraal(ws.rasterBuffer, left, x+x1, y+y1, x+x2, y+y2) // blue horz
	// rasterizer.DrawLineAmmeraal(ws.rasterBuffer, down, x+x3, y+y3, x+x2, y+y2) // red
	// rasterizer.DrawLineAmmeraal(ws.rasterBuffer, down, x+x3, y+y3, x+x1, y+y1) // red
	tri.Set(x+x1, y+y1, x+x2, y+y2, x+x3, y+y3)
	tri.Fill(ws.rasterBuffer)

	ws.rasterBuffer.SetPixelColor(color.RGBA{R: 0, G: 255, B: 255, A: 127})
	// Triangle flat-top ----------------------------------
	x = 200
	y = 50

	x1 = 25
	y1 = 50
	x2 = 0
	y2 = 0
	x3 = 50
	y3 = 0
	// Make sure Y's are consitent
	// rasterizer.Sort(&x1, &y1, &x2, &y2, &x3, &y3)

	down = false
	// rasterizer.DrawLineAmmeraal(ws.rasterBuffer, left, x+x1, y+y1, x+x2, y+y2) // blue horz
	// rasterizer.DrawLineAmmeraal(ws.rasterBuffer, down, x+x2, y+y2, x+x3, y+y3) // red
	// rasterizer.DrawLineAmmeraal(ws.rasterBuffer, down, x+x3, y+y3, x+x1, y+y1) // red

	tri.SetWithZ(x+x1, y+y1, 2.0, x+x2, y+y2, 2.0, x+x3, y+y3, 2.0)
	tri.Fill(ws.rasterBuffer)

	ws.rasterBuffer.SetPixelColor(color.RGBA{R: 255, G: 255, B: 255, A: 255})

	// Triangle split top ----------------------------------
	x = 200
	y = 200

	x1 = 25
	y1 = 50
	x2 = 0
	y2 = -50
	x3 = 50
	y3 = 0
	tri.Set(x+x1, y+y1, x+x2, y+y2, x+x3, y+y3)
	tri.Fill(ws.rasterBuffer)

	// Triangle split bottom ----------------------------------
	x = 350
	y = 200

	if ws.animate || ws.step {
		if ws.xx2 < -50 {
			ws.dir2 = 2
		} else if ws.xx2 > 100 {
			ws.dir2 = -2
		}
		ws.xx2 += ws.dir2
	}
	x1 = ws.xx2

	//y1 = 100
	if ws.animate || ws.step {
		if ws.xx3 < 0 {
			ws.dir3 = 1
		} else if ws.xx3 > 100 {
			ws.dir3 = -1
		}
		ws.xx3 += ws.dir3
	}
	y1 = ws.xx3

	if ws.animate || ws.step {
		if ws.xx < -50 {
			ws.dir = 1
		} else if ws.xx > 100 {
			ws.dir = -1
		}
		ws.xx += ws.dir
	}
	x2 = ws.xx // 75 cause overdraw, 50 is fine
	// fmt.Println("x1 ", x1, "y1 ", y1, "x2 ", x2)
	y2 = 50
	x3 = 25
	y3 = 0
	// fmt.Println(x+x1, y+y1, x+x2, y+y2, x+x3, y+y3)
	ws.step = false

	tri.Set(x+x1, y+y1, x+x2, y+y2, x+x3, y+y3)
	tri.Fill(ws.rasterBuffer)
}

// Quit stops the gui from running, effectively shutting it down.
func (ws *WindowSurface) Quit() {
	ws.running = false
}

// Close closes the viewer.
// Be sure to setup a "defer x.Close()"
func (ws *WindowSurface) Close() {
	if !ws.opened {
		return
	}
	var err error

	if ws.nFont == nil {
		return
	}
	ws.nFont.Destroy()

	ws.txtFPSLabel.Destroy()
	ws.txtMousePos.Destroy()
	ws.dynaTxt.Destroy()

	log.Println("Destroying texture")
	err = ws.texture.Destroy()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Destroying renderer")
	ws.renderer.Destroy()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Destroying window")
	err = ws.window.Destroy()
	if err != nil {
		log.Fatal(err)
	}

	sdl.Quit()

	if err != nil {
		log.Fatal(err)
	}
}

func (ws *WindowSurface) clearDisplay() {
	ws.rasterBuffer.Clear()
	ws.window.UpdateSurface()
}

// SetDrawColor --
func (ws *WindowSurface) SetDrawColor(color sdl.Color) {
}

// SetPixel --
func (ws *WindowSurface) SetPixel(x, y int) {
}

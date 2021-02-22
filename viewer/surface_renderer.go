package viewer

import (
	"SoftRenderer/api"
	"SoftRenderer/renderer"
	"fmt"
	"image"
	"log"
	"math"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

// RendererSurface is
type RendererSurface struct {
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
	paused bool

	nFont        *Font
	txtSimStatus *Text
	txtFPSLabel  *Text
	txtLoopLabel *Text
	txtMousePos  *Text
	dynaTxt      *DynaText
}

// NewRendererSurface creates a new viewer and initializes it.
func NewRendererSurface() api.ISurface {
	v := new(RendererSurface)
	v.opened = false
	v.paused = false
	v.mod = 200
	return v
}

// Open shows the viewer and begins event polling
// (host deuron.IHost)
func (v *RendererSurface) Open() {
	v.initialize()

	v.opened = true
}

// SetFont sets the font based on path and size.
func (v *RendererSurface) SetFont(fontPath string, size int) error {
	var err error
	v.nFont, err = NewFont(fontPath, size)
	return err
}

// filterEvent returns false if it handled the event. Returning false
// prevents the event from being added to the queue.
func (v *RendererSurface) filterEvent(e sdl.Event, userdata interface{}) bool {

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
		// println(t.Keysym.Scancode)
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
			case sdl.SCANCODE_P:
				// Pause sim
				v.paused = !v.paused
				// if v.paused {
				// 	fmt.Println("Paused")
				// }
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
func (v *RendererSurface) Run() {
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

	rasterizer := renderer.NewBresenHamRasterizer()

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

		v.AmmeraalTests(rasterizer)

		v.txtFPSLabel.DrawAt(10, 10)
		f := fmt.Sprintf("%2.2f", 1.0/elapsedTime*1000.0)
		v.dynaTxt.DrawAt(v.txtFPSLabel.Bounds.W+10, 10, f)

		// v.mx, v.my, _ = sdl.GetMouseState()
		// v.txtMousePos.DrawAt(10, 25)
		// f = fmt.Sprintf("<%d, %d>", v.mx/SurfaceScale, v.my/SurfaceScale)
		// v.dynaTxt.DrawAt(v.txtMousePos.Bounds.W+10, 25, f)

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
func (v *RendererSurface) Quit() {
	v.running = false
}

// Close closes the viewer.
// Be sure to setup a "defer x.Close()"
func (v *RendererSurface) Close() {
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

func (v *RendererSurface) initialize() {
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

	v.renderer, err = sdl.CreateRenderer(v.window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}

	v.renderer.SetDrawBlendMode(sdl.BLENDMODE_BLEND)
	rw := width / 320
	rh := height / 240

	v.renderer.SetLogicalSize(int32(width/rw), int32(height/rh))
}

// Configure view with draw objects
func (v *RendererSurface) Configure() {
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

func (v *RendererSurface) clearDisplay() {
	v.renderer.SetDrawColor(127, 127, 127, 255)
	v.renderer.Clear()
	// v.window.UpdateSurface()
}

// SetDrawColor --
func (v *RendererSurface) SetDrawColor(color sdl.Color) {
	v.renderer.SetDrawColor(color.R, color.G, color.B, color.A)
}

// SetPixel --
func (v *RendererSurface) SetPixel(x, y int) {
	v.renderer.DrawPoint(int32(x), int32(y))
}

// BasicTests --
func (v *RendererSurface) BasicTests(rasterizer api.IRasterizer) {
	// v.renderer.SetDrawColor(64, 64, 64, 255)
	// v.renderer.FillRect(&rect)

	// v.renderer.SetDrawColor(255, 255, 255, 127)
	// v.renderer.FillRect(&sdl.Rect{X: 80, Y: 80, W: 25, H: 25})

	// v.renderer.SetDrawColor(255, 255, 255, 255)
	// v.renderer.DrawPoint(150, 150)

	// v.renderer.SetDrawColor(255, 0, 0, 255)
	// rasterizer.DrawLine(v, 250, 200, 100, y)
	// rasterizer.DrawLine(v, 200, 50, 250, y)

	// if !v.paused {
	// 	y += yDir

	// 	if y < 5 {
	// 		yDir = 1
	// 	}
	// 	if y > 150 {
	// 		yDir = -1
	// 	}
	// }
	// rect := sdl.Rect{X: 50, Y: 50, W: 50, H: 50}
	// yDir := -1

	// Sorting:
	// Before: 50  100  100 50  50 50
	// After:  100 50   50  50  50 100

	//  50,50
	//    3 G  2              y1
	//    .---.  100, 50
	//    |  /
	//  R | / O
	//    |/                  y2
	//    .
	//    1
	//   50,100               y3
	// Trigger flat-top              1        2       3
	// rasterizer.DrawTriangle(v, 50, 100, 100, 50, 50, 50)

	//   50,50
	//    3
	//    .
	//    |\
	//  R | \  O
	//    |  \
	//    .---.
	//   1  G  2
	// 50,100   100,100
	// Trigger flat-bottom       1         2        3
	// rasterizer.DrawTriangle(v, 50, 100, 100, 100, 50, 50)

	//       50,50
	//        3
	//        .
	//        |\
	//        | \
	//      R |  \ B
	//        |   \
	//        .-G--.
	//        1\    \
	//  50,100   \   \
	//          R  \  \ B
	//               \ \
	//                 \\
	//                   .
	//                  2
	//               100,150
	// Trigger split             1         2        3
	// rasterizer.DrawTriangle(v, 50, 100, 100, 150, 50, 50)
}

// AmmeraalTests --
func (v *RendererSurface) AmmeraalTests(rasterizer api.IRasterizer) {
	v.renderer.SetDrawColor(255, 127, 127, 255)
	// rasterizer.DrawLineAmmeraal(v, v.paused, 100, 100, 200, 125) // blue

	x := 0
	y := 0
	left := false
	rasterizer.DrawLineAmmeraal(v, left, x, y, x+100, y+25) // blue dx>0
	x = 0
	y = 0
	rasterizer.DrawLineAmmeraal(v, left, x, y+25, x+100, y) // blue dx>0

	x = 50
	y = 50
	down := true
	rasterizer.DrawLineAmmeraal(v, down, x, y+50, x+50, y+150) // red
	x = 50
	y = 50
	rasterizer.DrawLineAmmeraal(v, down, x+50, y+50, x, y+150) // red

	// Horizontal
	x = 100
	y = 5
	left = false
	rasterizer.DrawLineAmmeraal(v, left, x, y, x+100, y) // blue
	x = 100
	y = 10
	left = true
	rasterizer.DrawLineAmmeraal(v, left, x, y, x+100, y) // blue

	// Vertical
	x = 100
	y = 20
	down = true
	rasterizer.DrawLineAmmeraal(v, down, x, y, x, y+100) // red
	x = 110
	y = 20
	down = false
	rasterizer.DrawLineAmmeraal(v, down, x, y, x, y+100) // red

	// Triangle ----------------------------------
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

	down = false
	rasterizer.DrawLineAmmeraal(v, left, x+x1, y+y1, x+x2, y+y2) // blue horz
	rasterizer.DrawLineAmmeraal(v, down, x+x3, y+y3, x+x2, y+y2) // red
	rasterizer.DrawLineAmmeraal(v, down, x+x3, y+y3, x+x1, y+y1) // red

	x = 200
	y = 100

	x1 = 25
	y1 = 50
	x2 = 0
	y2 = 0
	x3 = 50
	y3 = 0
	// Make sure Y's are consitent
	// rasterizer.Sort(&x1, &y1, &x2, &y2, &x3, &y3)

	down = false
	rasterizer.DrawLineAmmeraal(v, left, x+x1, y+y1, x+x2, y+y2) // blue horz
	rasterizer.DrawLineAmmeraal(v, down, x+x2, y+y2, x+x3, y+y3) // red
	rasterizer.DrawLineAmmeraal(v, down, x+x3, y+y3, x+x1, y+y1) // red

	// rasterizer.DrawLineAmmeraalDyGtDx2(v, v.paused, 100, 200, 150, 50) // red
	// rasterizer.DrawLineAmmeraalDyGtDx2(v, false, 100, 50, 150, 200) // red
}

package main

// Note: You may need a ram drive:
// export TMPDIR="/Volumes/RAMDisk"
// diskutil erasevolume HFS+ 'RAMDisk' `hdiutil attach -nomount ram://2097152`

import (
	"github.com/wdevore/SoftRenderer/viewer"
)

var gview viewer.IViewer // The GUI

func run() int {
	gview = viewer.NewViewer()
	defer gview.Close()

	gview.Open()

	gview.SetFont("galacticstormexpand.ttf", 16)
	gview.Configure()

	gview.Run()

	return 0
}

func main() {
	run()
	// // os.Exit(..) must run AFTER sdl.Main(..) below; so keep track of exit
	// // status manually outside the closure passed into sdl.Main(..) below
	// var exitcode int
	// sdl.Main(func() {
	// 	exitcode = run()
	// })
	// // os.Exit(..) must run here! If run in sdl.Main(..) above, it will cause
	// // premature quitting of sdl.Main(..) function; resource cleaning deferred
	// // calls/closing of channels may never run
	// os.Exit(exitcode)
}

package main

import "SoftRenderer/viewer"

func main() {
	surface := viewer.NewRendererSurface()
	defer surface.Close()

	surface.Open()

	// err := gview.SetFont("../../assets/galacticstormexpand.ttf", 16)
	err := surface.SetFont("../../assets/MontserratAlternates-Light.otf", 16)

	if err != nil {
		panic(err)
	}

	surface.Configure()

	surface.Run()
}

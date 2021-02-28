package main

import "SoftRenderer/surface"

func main() {
	surface := surface.NewSurfaceBuffer()
	defer surface.Close()

	surface.Open()

	err := surface.SetFont("../../assets/MontserratAlternates-Light.otf", 16)

	if err != nil {
		panic(err)
	}

	surface.Configure()

	surface.Run()
}

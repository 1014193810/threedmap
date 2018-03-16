package main

import (
	"image"
	"image/color"
	"math"

	"os"

	_ "image/jpeg"
	_ "image/png"

	"threedmap"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

func createsmap() (*threedmap.Smap, error) {

	m := threedmap.Initmap(threedmap.Ssize{0, 0, 1000, 1000})

	for i := 0; i < 1000; i++ {
		for j := 0; j < 1000; j++ {
			//z := 0.0
			z := 100*math.Sin(math.Pi/100*float64(i)) + 100
			m.SetZ(j, i, z)
			m.SetColor(j, i, color.RGBA{uint8(250 * (math.Abs(math.Sin(math.Pi / 100 * float64(i))))), uint8(250 * math.Abs(math.Sin(math.Pi/100*float64(i)))), uint8(250 * math.Abs(math.Sin(math.Pi/100*float64(i)))), 255})
		}
	}

	return m, nil
}

func loadpic(path string) (*threedmap.Smap, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	Image, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	x0, y0, x, y := Image.Bounds().Min.X, Image.Bounds().Min.Y, Image.Bounds().Max.X, Image.Bounds().Max.Y
	m := threedmap.Initmap(threedmap.Ssize{x0, y0, x, y})
	////
	for i := x0; i < x; i++ {
		for j := y0; j < y; j++ {
			//z := 0.0
			Color := Image.At(i, j)
			_, z, _, _ := Color.RGBA()
			m.SetZ(i, j, float64(z)/256/7)
			m.SetColor(i, j, Color)
		}
	}
	////
	return m, nil
}

func main() {
	cfg := pixelgl.WindowConfig{
		Title:  "Pixel Rocks!",
		Bounds: pixel.R(0, 0, 1024, 768),
	}
	m, err := createsmap() //create a smap to show
	if err != nil {
		panic(err)
	}
	run := threedmap.GenerateRunfunc(m, cfg, colornames.White)
	pixelgl.Run(run)
}

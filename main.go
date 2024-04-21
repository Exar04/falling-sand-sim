package main

import (
	"fmt"
	"image"
	_ "image/png"
	"os"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

type Square struct {
	posx float64
	posy float64

	size float64

	color pixel.RGBA
}

func CreateSquare(posX, posY, s float64, color pixel.RGBA) Square {
	newRec := Square{}

	newRec.posx = posX
	newRec.posy = windowSizeY - posY - s
	newRec.size = s
	newRec.color = color

	return newRec
}

func (r *Square) DrawSquare() *imdraw.IMDraw {

	imd := imdraw.New(nil)

	// imd.Color = r.color
	// imd.Color = pixel.RGBA{R: 1, G: 1, B: 1, A: 0}

	imd.Push(pixel.V(r.posx, r.posy))
	imd.Push(pixel.V(r.posx+r.size, r.posy))
	imd.Push(pixel.V(r.posx+r.size, r.posy+r.size))
	imd.Push(pixel.V(r.posx, r.posy+r.size))

	imd.Polygon(0)

	return imd
}

func main() {
	pixelgl.Run(run)
}

// so to create a sand effect we first need to divide our window into a matrix
// after that when a on click even occours on mouse we will take the position of the mouse and put a rectangle there
// and then check its adjecent blocks if they are empty the rectangle will fall down
// and if we have mouse clicked and there is already a block on that matrix we will not add a new one

const windowSizeX = 800
const windowSizeY = 800

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Pixel Rocks!",
		Bounds: pixel.R(0, 0, windowSizeX, windowSizeY),
		VSync:  true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	nero := CreateSquare(10, 10, 5, pixel.RGB(1, 1, 1))
	imd := nero.DrawSquare()

	for !win.Closed() {

		imd.Draw(win)
		win.Update()
	}
}
func runo() {
	cfg := pixelgl.WindowConfig{
		Title:  "Pixel Rocks!",
		Bounds: pixel.R(0, 0, windowSizeX, windowSizeY),
		VSync:  true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	// creating a matrix of size of window
	rows := 8
	cols := 8
	squareMat := make([][]int, rows)
	for i := range squareMat {
		squareMat[i] = make([]int, cols)
	}

	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			squareMat[i][j] = j // Assigning a simple value for demonstration
		}
	}

	for !win.Closed() {

		win.Clear(pixel.RGB(0, 0, 0))

		for i := 0; i < rows; i++ {
			for j := 0; j < cols; j++ {
				if squareMat[i][j] == 2 {

				}
				// if squareMat[i][j] == 1 {
				// 	imd := imdraw.New(nil)

				// 	imd.EndShape = imdraw.RoundEndShape
				// 	transparentWhite := colornames.White
				// 	transparentWhite.A = uint8(normalizeTo_0_255(float64(j), 0, 8))
				// 	imd.Color = transparentWhite

				// 	imd.Push(pixel.V(1024/2, 780/2), pixel.V(1024/2, 780/2))
				// 	imd.Line(5)
				// 	imd.Draw(win)
				// }
			}
		}

		win.Update()
	}

}

func normalizeTo_0_255(value, min, max float64) float64 {
	return (value - min) * (255.0 / (max - min))
}

// func normalizeToWindowSizeX(value, min, max, newMin, newMax float64) float64 {
// 	return ((value - min) * (newMax - newMin) / (max - min)) + newMin
// }

func normalizeToWindowSizeX(value, min, max float64) float64 {
	// return (value - min) * (windowSizeX / (max - min))
	return (value-min)*(((windowSizeX-10)-10)/(max-min)) + 10
}
func normalizeToWindowSizeY(value, min, max float64) float64 {
	// return (value - min) * (windowSizeY / (max - min))
	return (value-min)*(((windowSizeY-10)-10)/(max-min)) + 10
}

func loadImage(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open image: %w", err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	return pixel.PictureDataFromImage(img), nil
}

func eventHandler(win *pixelgl.Window) {
	if win.Pressed(pixelgl.MouseButtonLeft) {
		fmt.Println("left clicked")
	}

	if win.Pressed(pixelgl.MouseButtonRight) {
		fmt.Println("right clicked")
	}
}

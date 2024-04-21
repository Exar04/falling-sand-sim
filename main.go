package main

import (
	"fmt"
	"image"
	_ "image/png"
	"os"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

const windowSizeX = 800
const windowSizeY = 800

const rows = 200
const cols = 200
const sizeOfBlock = float64(windowSizeX / rows)

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

	imd.Color = r.color

	imd.Push(pixel.V(r.posx, r.posy))
	imd.Push(pixel.V(r.posx+r.size, r.posy))
	imd.Push(pixel.V(r.posx+r.size, r.posy+r.size))
	imd.Push(pixel.V(r.posx, r.posy+r.size))

	imd.Polygon(0)

	return imd
}

func (r *Square) pushSqr() [4]pixel.Vec {
	arr := [4]pixel.Vec{}

	arr[0] = pixel.Vec{X: r.posx, Y: r.posy}
	arr[1] = pixel.Vec{X: r.posx + r.size, Y: r.posy}
	arr[2] = pixel.Vec{X: r.posx + r.size, Y: r.posy + r.size}
	arr[3] = pixel.Vec{X: r.posx, Y: r.posy + r.size}

	return arr
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

	imd := imdraw.New(nil)
	imd.Push(pixel.V(100, 200))
	imd.Push(pixel.V(100, 100))
	imd.Push(pixel.V(200, 100))
	imd.Push(pixel.V(200, 200))

	imd.Polygon(0)

	imd.Push(pixel.V(300, 400))
	imd.Push(pixel.V(300, 300))
	imd.Push(pixel.V(400, 300))
	imd.Push(pixel.V(400, 400))

	imd.Polygon(0)

	imd.Draw(win)
	for !win.Closed() {

		imd.Draw(win)
		win.Update()
	}
}
func main() {
	pixelgl.Run(run)
}

// so to create a sand effect we first need to divide our window into a matrix
// after that when a on click even occours on mouse we will take the position of the mouse and put a rectangle there
// and then check its adjecent blocks if they are empty the rectangle will fall down
// and if we have mouse clicked and there is already a block on that matrix we will not add a new one

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

	squareMat := make([][]bool, rows)
	preDefPositionMatrix := make([][]pixel.Vec, rows)
	normalisedMatrixPositionValues := make([][]pixel.Vec, rows)

	for i := range squareMat {
		squareMat[i] = make([]bool, cols)
		preDefPositionMatrix[i] = make([]pixel.Vec, cols)
		normalisedMatrixPositionValues[i] = make([]pixel.Vec, cols)
		// rasterizedColors[i] = make([]int, cols)
	}

	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			squareMat[i][j] = false

			pb := getPositionOfBlock(i, j, (windowSizeX / rows))
			preDefPositionMatrix[i][j] = pixel.Vec{X: pb.X, Y: pb.Y}

			npx := normalizeToWindowSizeX(float64(i), 0, float64(rows))
			npy := normalizeToWindowSizeY(float64(j), 0, float64(cols))

			normalisedMatrixPositionValues[i][j] = pixel.Vec{X: npx, Y: npy}
		}
	}

	rasterizedColors := make([]int, rows*cols)
	for !win.Closed() {

		win.Clear(pixel.RGB(0, 0, 0))
		rasterizedMatrix := make([][4]pixel.Vec, 0)

		imd := imdraw.New(nil)

		for i := 0; i < rows; i++ {
			for j := 0; j < cols; j++ {

				if win.Pressed(pixelgl.MouseButtonLeft) {
					if isMouseInsideBlock(win, preDefPositionMatrix[i][j].X, preDefPositionMatrix[i][j].Y, sizeOfBlock) {
						rasterizedColors[i*cols+j] = 1
					}
				}

				if rasterizedColors[i*cols+j] == 1 && j < cols-1 {
					if rasterizedColors[i*cols+(j+1)] == 0 {
						rasterizedColors[i*cols+(j+1)] = 2
						rasterizedColors[i*cols+j] = 0
					}
				}
				if rasterizedColors[i*cols+j] == 2 && j < cols-1 {
					rasterizedColors[i*cols+j] = 1
				}

				nero := CreateSquare(
					normalisedMatrixPositionValues[i][j].X,
					normalisedMatrixPositionValues[i][j].Y,
					float64(windowSizeX/rows), pixel.RGB(1, 1, 1),
				)
				rasterizedMatrix = append(rasterizedMatrix, nero.pushSqr())
			}
		}

		// fmt.Println(rasterizedColors)
		for i := range rasterizedMatrix {
			var colo float64
			if rasterizedColors[i] == 1 || rasterizedColors[i] == 2 {
				colo = 1
			} else {
				colo = 0
			}
			// colo := float64(rasterizedColors[i])
			imd.Color = pixel.RGB(colo, colo, colo)
			for j := range rasterizedMatrix[i] {
				imd.Push(rasterizedMatrix[i][j])
			}
			imd.Polygon(0)
		}
		imd.Draw(win)

		win.Update()
		time.Sleep(time.Millisecond * 10)
	}
	// fmt.Println(rasterizedMatrix)

}

// in an array there are blocks we check if our mouse is inside a given block
// we will calculate it by ckecking if our mouse is inside the range of your block
func isMouseInsideBlock(win *pixelgl.Window, blockPosX, blockPosY, size float64) bool {
	mp := win.MousePosition()
	if mp.X >= blockPosX && mp.X <= blockPosX+size && mp.Y >= blockPosY && mp.Y <= blockPosY+size {
		return true
	} else {
		return false
	}
}

func getPositionOfBlock(x, y, s int) pixel.Vec {
	px := normalizeToWindowSizeX(float64(x), 0, float64(rows))
	py := windowSizeY - normalizeToWindowSizeY(float64(y), 0, float64(cols)) - float64(s)
	p := pixel.Vec{X: px, Y: py}
	return p
}

func normalizeTo_0_255(value, min, max float64) float64 {
	return (value - min) * (255.0 / (max - min))
}

func normalizeTo_0_1(value, min, max float64) float64 {
	return (value - min) * (1.0 / (max - min))
}

// min represents the mininum range of value we are giving that is 5 from range (0 to 10)
// this function converts our value from our given range to window size range
func normalizeToWindowSizeX(value, min, max float64) float64 {
	return (value - min) * (windowSizeX / (max - min))
}
func normalizeToWindowSizeY(value, min, max float64) float64 {
	return (value - min) * (windowSizeY / (max - min))
}

func normalizeToWindowSizeXWithPadding(value, min, max, padding float64) float64 {
	return (value-min)*((windowSizeX-padding*2)/(max-min)) + padding
}
func normalizeToWindowSizeYWithPadding(value, min, max, padding float64) float64 {
	return (value-min)*((windowSizeY-padding*2)/(max-min)) + padding
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
	}

	if win.Pressed(pixelgl.MouseButtonRight) {
		fmt.Println("right clicked")
	}
}

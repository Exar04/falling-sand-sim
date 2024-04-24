package main

import (
	"fmt"
	"image"
	_ "image/png"
	"os"
	"sync"

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

		for i := 0; i < rows; i++ {
			for j := 0; j < cols; j++ {
				npx := normalizeToWindowSizeX(float64(i), 0, float64(rows))
				npy := normalizeToWindowSizeY(float64(j), 0, float64(cols))

				sq1 := CreateSquare(npx, npy, windowSizeX/rows, pixel.RGB(1, 1, 1))
				drawsq1 := sq1.DrawSquare()

				sq2 := CreateSquare(npx+3, npy+3, windowSizeX/rows-6, pixel.RGB(0, 0, 0))
				drawsq2 := sq2.DrawSquare()

				drawsq1.Draw(win)
				drawsq2.Draw(win)
			}
		}
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
	// rasterizedColors := make([]int, rows*cols)
	rasterizedColors := make([][]int, rows)

	sandColor := [3]float64{0.96, 0.84, 0.65}
	waterColor := [3]float64{0.65, 0.88, 0.93}

	for i := range squareMat {
		squareMat[i] = make([]bool, cols)
		preDefPositionMatrix[i] = make([]pixel.Vec, cols)
		normalisedMatrixPositionValues[i] = make([]pixel.Vec, cols)
		rasterizedColors[i] = make([]int, cols)
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

	// throw mouse listener in different thread
	// go mouseLeftClickListener(win, colorSwitch)

	wg := &sync.WaitGroup{}
	for !win.Closed() {

		win.Clear(pixel.RGB(0.1, .17, .12))
		rasterizedMatrix := make([][4]pixel.Vec, 0)

		imd := imdraw.New(nil)

		for i := 0; i < rows; i++ {
			// this function traverses row by row
			func(i int, wg *sync.WaitGroup) {
				for j := 0; j < cols; j++ {

					// to traverse colums instead of rows

					if win.Pressed(pixelgl.MouseButtonLeft) {
						if isMouseInsideBlock(win, preDefPositionMatrix[i][j].X, preDefPositionMatrix[i][j].Y, sizeOfBlock) {
							rasterizedColors[i][j] = 1 // 1 represents sand
						}
					}
					if win.Pressed(pixelgl.MouseButtonRight) {
						if isMouseInsideBlock(win, preDefPositionMatrix[i][j].X, preDefPositionMatrix[i][j].Y, sizeOfBlock) {
							rasterizedColors[i][j] = 3 // 3 represents water
						}
					}
					// sim for sand
					if rasterizedColors[i][j] == 1 {
						if rasterizedColors[i][j+1] == 0 {
							rasterizedColors[i][j+1] = 2
							rasterizedColors[i][j] = 0
						} else if i < rows-1 && rasterizedColors[i+1][j+1] == 0 {
							rasterizedColors[i+1][j+1] = 2
							rasterizedColors[i][j] = 0
						} else if i > 0 && rasterizedColors[i-1][j+1] == 0 {
							rasterizedColors[i-1][j+1] = 2
							rasterizedColors[i][j] = 0
						} else if i > 0 && (rasterizedColors[i][j+1] == 3 || rasterizedColors[i][j+1] == 4) {
							rasterizedColors[i][j+1] = 2
							rasterizedColors[i][j] = 3
						}
					}
					// 2 represents that sand block had been moved from above row to below row
					if rasterizedColors[i][j] == 2 && j < cols-1 {
						rasterizedColors[i][j] = 1
					}

					// sim for water
					if rasterizedColors[i][j] == 3 {
						if rasterizedColors[i][j+1] == 0 {
							rasterizedColors[i][j+1] = 4
							rasterizedColors[i][j] = 0
						} else if i < rows-1 && rasterizedColors[i+1][j+1] == 0 {
							rasterizedColors[i+1][j+1] = 4
							rasterizedColors[i][j] = 0
						} else if i > 0 && rasterizedColors[i-1][j+1] == 0 {
							rasterizedColors[i-1][j+1] = 4
							rasterizedColors[i][j] = 0
						} else if i > 0 && i < rows-1 {
							if rasterizedColors[i+1][j] == 0 && i < rows-2 {
								if rasterizedColors[i+2][j] == 0 {
									rasterizedColors[i+1][j] = 4
									rasterizedColors[i][j] = 0
								}
							} else if rasterizedColors[i-1][j] == 0 {
								rasterizedColors[i-1][j] = 4
								rasterizedColors[i][j] = 0
							}

						}
					}

					// 4 represents that water block had been moved from above row to below row
					if rasterizedColors[i][j] == 4 && j < cols-1 {
						rasterizedColors[i][j] = 3
					}

					nero := CreateSquare(
						normalisedMatrixPositionValues[i][j].X,
						normalisedMatrixPositionValues[i][j].Y,
						float64(windowSizeX/rows), pixel.RGB(1, 1, 1),
					)
					rasterizedMatrix = append(rasterizedMatrix, nero.pushSqr())
				}
			}(i, wg)
		}

		//	// this funciton traverses column by column
		// for i := 0; i < rows; i++ {
		// 	if i%2 == 0 {
		// 		wg.Add(1)
		// 		go func(i int, wg *sync.WaitGroup) {
		// 			defer wg.Done()
		// 			for j := 0; j < cols; j++ {
		// 				// to traverse colums instead of rows

		// 				if win.Pressed(pixelgl.MouseButtonLeft) {
		// 					if isMouseInsideBlock(win, preDefPositionMatrix[j][i].X, preDefPositionMatrix[j][i].Y, sizeOfBlock) {
		// 						if colorSwitch {
		// 							rasterizedColors[j][i] = 1 // 1 represents sand
		// 						} else {
		// 							rasterizedColors[j][i] = 3 // 3 represents water
		// 						}
		// 					}
		// 				}
		// 				// sim for sand
		// 				if i < rows-1 && rasterizedColors[j][i] == 1 {
		// 					if rasterizedColors[j][i+1] == 0 {
		// 						rasterizedColors[j][i+1] = 2
		// 						rasterizedColors[j][i] = 0
		// 					} else if j < rows-2 && rasterizedColors[j+1][i+1] == 0 {
		// 						rasterizedColors[j+1][i+1] = 2
		// 						rasterizedColors[j][i] = 0
		// 					} else if j > 0 && rasterizedColors[j-1][i+1] == 0 {
		// 						rasterizedColors[j-1][i+1] = 2
		// 						rasterizedColors[j][i] = 0
		// 					} else if j > 0 && (rasterizedColors[j][i+1] == 3 || rasterizedColors[j][i+1] == 4) {
		// 						rasterizedColors[j][i+1] = 2
		// 						rasterizedColors[j][i] = 3
		// 					}
		// 				}
		// 				// 2 represents that sand block had been moved from above row to below row
		// 				if rasterizedColors[j][i] == 2 && j < cols-1 {
		// 					rasterizedColors[j][i] = 1
		// 				}

		// 				// sim for water
		// 				if rasterizedColors[j][i] == 3 {
		// 					if rasterizedColors[j][i+1] == 0 {
		// 						rasterizedColors[j][i+1] = 4
		// 						rasterizedColors[j][i] = 0
		// 					} else if j < rows-1 && rasterizedColors[j+1][i+1] == 0 {
		// 						rasterizedColors[j+1][i+1] = 4
		// 						rasterizedColors[j][i] = 0
		// 					} else if j > 0 && rasterizedColors[j-1][i+1] == 0 {
		// 						rasterizedColors[j-1][i+1] = 4
		// 						rasterizedColors[j][i] = 0
		// 					} else if j > 0 && j < rows-1 {
		// 						if rasterizedColors[j+1][i] == 0 && j < rows-2 {
		// 							if rasterizedColors[j+2][i] == 0 {
		// 								rasterizedColors[j+1][i] = 4
		// 								rasterizedColors[j][i] = 0
		// 							}
		// 						} else if rasterizedColors[j-1][i] == 0 {
		// 							rasterizedColors[j-1][i] = 4
		// 							rasterizedColors[j][i] = 0
		// 						}
		// 					}
		// 				}

		// 				// 4 represents that water block had been moved from above row to below row
		// 				if rasterizedColors[j][i] == 4 && i < cols-1 {
		// 					rasterizedColors[j][i] = 3
		// 				}

		// 				nero := CreateSquare(
		// 					normalisedMatrixPositionValues[i][j].X,
		// 					normalisedMatrixPositionValues[i][j].Y,
		// 					float64(windowSizeX/rows), pixel.RGB(1, 1, 1),
		// 				)
		// 				mut.Lock()
		// 				rasterizedMatrix = append(rasterizedMatrix, nero.pushSqr())
		// 				mut.Unlock()
		// 			}
		// 		}(i, wg)
		// 	}
		// }
		// wg.Wait()

		for i, row := range rasterizedColors {
			for j := range row {
				if rasterizedColors[i][j] == 1 || rasterizedColors[i][j] == 2 {
					imd.Color = pixel.RGB(sandColor[0], sandColor[1], sandColor[2])

				} else if rasterizedColors[i][j] == 3 || rasterizedColors[i][j] == 4 {
					imd.Color = pixel.Alpha(0.5).Mul(pixel.RGB(waterColor[0], waterColor[1], waterColor[2]))
				} else {
					imd.Color = pixel.RGB(0.1, 0.17, 0.12)
				}

				for k := range rasterizedMatrix[i*cols+j] {
					// fmt.Printf("k : %d, i : %d, j : %d, cols: %d, i*cols+j : %d , rclen: %d \n", k, i, j, cols, i*cols+j, len(rasterizedColors))
					imd.Push(rasterizedMatrix[i*cols+j][k])
				}
				imd.Polygon(0)
			}
		}
		imd.Draw(win)

		win.Update()
		// time.Sleep(time.Millisecond * 10)
	}
}

// func wt(i int, wg *sync.WaitGroup, win *pixelgl.Window, preDefPositionMatrix [][]pixel.Vec, rasterizedColors [][]int, normalisedMatrixPositionValues [][]pixel.Vec, rasterizedMatrix [][4]pixel.Vec, colorSwitch bool) {
// }

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

func mouseLeftClickListener(win *pixelgl.Window, colorSwitch bool) {
	for true {
		if win.Pressed(pixelgl.MouseButtonLeft) {
			// if isMouseInsideBlock(win, preDefPositionMatrix[i][j].X, preDefPositionMatrix[i][j].Y, sizeOfBlock) {
			// 	if colorSwitch {
			// 		rasterizedColors[i][j] = 1 // 1 represents sand
			// 	} else {
			// 		rasterizedColors[i][j] = 3 // 3 represents water
			// 	}
			// }
		}
	}
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

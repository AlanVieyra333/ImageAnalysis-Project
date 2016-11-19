package editImage

import (
	"container/heap"
	"image"
	"image/color"
	"math"
	"tools"
)

/*	----------------	Low pass filter	----------------	*/
func ConvolutionMask(inImg image.Image, umb uint16, mask [3][3]int) *image.RGBA {
	n := inImg.Bounds().Max.Y
	m := inImg.Bounds().Max.X
	outImg := image.NewRGBA(image.Rect(0, 0, m, n))

	window := len(mask)

	for y := 0; y < n; y++ {
		for x := 0; x < m; x++ {
			var r, g, b uint16
			_, _, _, a := inImg.At(x, y).RGBA()

			var rGx, rGy, gGx, gGy, bGx, bGy int
			var i, j int
			for yy := y - (window / 2); yy <= y+(window/2); yy++ {
				i = 0
				for xx := x - (window / 2); xx <= x+(window/2); xx++ {
					if IsValid(xx, yy, n, m) {
						rr, gg, bb, _ := inImg.At(xx, yy).RGBA()

						rGx += mask[j][i] * int(rr)
						rGy += mask[window-1-i][j] * int(rr)
						gGx += mask[j][i] * int(gg)
						gGy += mask[window-1-i][j] * int(gg)
						bGx += mask[j][i] * int(bb)
						bGy += mask[window-1-i][j] * int(bb)
					}
					i++
				}
				j++
			}
			rG := uint16(math.Abs(float64(rGx)) + math.Abs(float64(rGy)))
			gG := uint16(math.Abs(float64(gGx)) + math.Abs(float64(gGy)))
			bG := uint16(math.Abs(float64(bGx)) + math.Abs(float64(bGy)))

			if rG > umb || gG > umb || bG > umb {
				r = uint16(0xffff)
				g = uint16(0xffff)
				b = uint16(0xffff)
			} else {
				r = uint16(0)
				g = uint16(0)
				b = uint16(0)
			}

			newColor := color.RGBA64{r, g, b, uint16(a)}
			outImg.Set(x, y, newColor)
		}
	}
	return outImg
}

func FilterLaplaciano(inImg image.Image, umb uint16) *image.RGBA {
	mask := [3][3]int{{-1, 0, 0},
		{0, 1, 0},
		{0, 0, 0}}
	return ConvolutionMask(inImg, umb, mask)
}

func FilterRobert(inImg image.Image, umb uint16) *image.RGBA {
	mask := [3][3]int{{-1, 0, 0},
		{0, 1, 0},
		{0, 0, 0}}
	return ConvolutionMask(inImg, umb, mask)
}

func FilterPrewitt(inImg image.Image, umb uint16) *image.RGBA {
	mask := [3][3]int{{1, 0, -1},
		{1, 0, -1},
		{1, 0, -1}}
	return ConvolutionMask(inImg, umb, mask)
}

func FilterSobel(inImg image.Image, umb uint16) *image.RGBA {
	mask := [3][3]int{{-1, 0, 1},
		{-2, 0, 2},
		{-1, 0, 1}}
	return ConvolutionMask(inImg, umb, mask)
}

/*	----------------	Low pass filter	----------------	*/
func FilterProm(inImg image.Image, window int8) *image.RGBA {
	n := inImg.Bounds().Max.Y
	m := inImg.Bounds().Max.X
	outImg := image.NewRGBA(image.Rect(0, 0, m, n))

	for y := 0; y < n; y++ {
		for x := 0; x < m; x++ {
			/*	--------------------	Process	--------------------	*/
			var r, g, b int64
			_, _, _, a := inImg.At(x, y).RGBA()

			var items uint32
			for yy := y - (int(window) / 2); yy <= y+(int(window)/2); yy++ {
				for xx := x - (int(window) / 2); xx <= x+(int(window)/2); xx++ {
					if IsValid(xx, yy, n, m) {
						rr, gg, bb, _ := inImg.At(xx, yy).RGBA()
						r += int64(rr)
						g += int64(gg)
						b += int64(bb)
						items++
					}
				}
			}
			r /= int64(items)
			g /= int64(items)
			b /= int64(items)

			newColor := color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)}
			outImg.Set(x, y, newColor)
		}
	}
	return outImg
}

func FilterPromH(inImg image.Image, window int8, N uint32) *image.RGBA {
	n := inImg.Bounds().Max.Y
	m := inImg.Bounds().Max.X
	outImg := image.NewRGBA(image.Rect(0, 0, m, n))

	for y := 0; y < n; y++ {
		for x := 0; x < m; x++ {
			/*	--------------------	Process	--------------------	*/
			var r, g, b int64
			_, _, _, a := inImg.At(x, y).RGBA()

			var items uint32
			for yy := y - (int(window) / 2); yy <= y+(int(window)/2); yy++ {
				for xx := x - (int(window) / 2); xx <= x+(int(window)/2); xx++ {
					if IsValid(xx, yy, n, m) {
						rr, gg, bb, _ := inImg.At(xx, yy).RGBA()

						if xx == x && yy == y {
							r += int64(N * rr)
							g += int64(N * gg)
							b += int64(N * bb)
						} else {
							r += int64(rr)
							g += int64(gg)
							b += int64(bb)
							items++
						}
					}
				}
			}
			r /= int64(items + N)
			g /= int64(items + N)
			b /= int64(items + N)

			newColor := color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)}
			outImg.Set(x, y, newColor)
		}
	}
	return outImg
}

/*	----------------	Other filter	----------------	*/
func FilterMediana(inImg image.Image, window int8) *image.RGBA {
	n := inImg.Bounds().Max.Y
	m := inImg.Bounds().Max.X
	outImg := image.NewRGBA(image.Rect(0, 0, m, n))

	for y := 0; y < n; y++ {
		for x := 0; x < m; x++ {
			/*	--------------------	Process	--------------------	*/
			var rr, gg, bb uint16
			_, _, _, a := inImg.At(x, y).RGBA()

			rPQ := tools.NewPriorityQueue(0)
			gPQ := tools.NewPriorityQueue(0)
			bPQ := tools.NewPriorityQueue(0)
			heap.Init(&rPQ)
			heap.Init(&gPQ)
			heap.Init(&bPQ)
			item := &tools.Item{}

			for yy := y - (int(window) / 2); yy <= y+(int(window)/2); yy++ {
				for xx := x - (int(window) / 2); xx <= x+(int(window)/2); xx++ {
					if IsValid(xx, yy, n, m) {
						rx, gx, bx, _ := inImg.At(xx, yy).RGBA()

						// Insert a new item and then modify its priority.
						item = &tools.Item{
							Priority: int(rx),
						}
						heap.Push(&rPQ, item)

						item = &tools.Item{
							Priority: int(gx),
						}
						heap.Push(&gPQ, item)

						item = &tools.Item{
							Priority: int(bx),
						}
						heap.Push(&bPQ, item)
					}
				}
			}

			// Take the items out; they arrive in decreasing priority order.
			for rPQ.Len() > 0 {
				itemR := heap.Pop(&rPQ).(*tools.Item)
				itemG := heap.Pop(&gPQ).(*tools.Item)
				itemB := heap.Pop(&bPQ).(*tools.Item)
				if rPQ.Len() == (int(window*window))/2 {
					rr = uint16(itemR.Priority)
					gg = uint16(itemG.Priority)
					bb = uint16(itemB.Priority)
				}
			}

			newColor := color.RGBA64{rr, gg, bb, uint16(a)}
			outImg.Set(x, y, newColor)
		}
	}
	return outImg
}

func FilterModa(inImg image.Image, window int8) *image.RGBA {
	n := inImg.Bounds().Max.Y
	m := inImg.Bounds().Max.X
	outImg := image.NewRGBA(image.Rect(0, 0, m, n))

	for y := 0; y < n; y++ {
		for x := 0; x < m; x++ {
			var r, g, b uint16
			_, _, _, a := inImg.At(x, y).RGBA()

			itemsR := map[int]int{}
			itemsG := map[int]int{}
			itemsB := map[int]int{}
			for yy := y - (int(window) / 2); yy <= y+(int(window)/2); yy++ {
				for xx := x - (int(window) / 2); xx <= x+(int(window)/2); xx++ {
					if IsValid(xx, yy, n, m) {
						rr, gg, bb, _ := inImg.At(xx, yy).RGBA()

						itemsR[int(rr)]++
						itemsG[int(gg)]++
						itemsB[int(bb)]++
					}
				}
			}

			rMax := 0
			for Value, Rep := range itemsR {
				if Rep > rMax {
					r = uint16(Value)
					rMax = Rep
				}
			}
			gMax := 0
			for Value, Rep := range itemsG {
				if Rep > gMax {
					g = uint16(Value)
					gMax = Rep
				}
			}
			bMax := 0
			for Value, Rep := range itemsB {
				if Rep > bMax {
					b = uint16(Value)
					bMax = Rep
				}
			}

			newColor := color.RGBA64{r, g, b, uint16(a)}
			outImg.Set(x, y, newColor)
		}
	}
	return outImg
}

func FilterMax(inImg image.Image, window int8) *image.RGBA {
	n := inImg.Bounds().Max.Y
	m := inImg.Bounds().Max.X
	outImg := image.NewRGBA(image.Rect(0, 0, m, n))

	for y := 0; y < n; y++ {
		for x := 0; x < m; x++ {
			var r, g, b uint16
			_, _, _, a := inImg.At(x, y).RGBA()

			for yy := y - (int(window) / 2); yy <= y+(int(window)/2); yy++ {
				for xx := x - (int(window) / 2); xx <= x+(int(window)/2); xx++ {
					if IsValid(xx, yy, n, m) {
						rr, gg, bb, _ := inImg.At(xx, yy).RGBA()

						if uint16(rr) > r {
							r = uint16(rr)
						}
						if uint16(gg) > g {
							g = uint16(gg)
						}
						if uint16(bb) > b {
							b = uint16(bb)
						}
					}
				}
			}

			newColor := color.RGBA64{r, g, b, uint16(a)}
			outImg.Set(x, y, newColor)
		}
	}
	return outImg
}

func FilterMin(inImg image.Image, window int8) *image.RGBA {
	n := inImg.Bounds().Max.Y
	m := inImg.Bounds().Max.X
	outImg := image.NewRGBA(image.Rect(0, 0, m, n))

	for y := 0; y < n; y++ {
		for x := 0; x < m; x++ {
			var r, g, b uint16
			_, _, _, a := inImg.At(x, y).RGBA()

			r = uint16((1 << 16) - 1)
			g = uint16((1 << 16) - 1)
			b = uint16((1 << 16) - 1)
			for yy := y - (int(window) / 2); yy <= y+(int(window)/2); yy++ {
				for xx := x - (int(window) / 2); xx <= x+(int(window)/2); xx++ {
					if IsValid(xx, yy, n, m) {
						rr, gg, bb, _ := inImg.At(xx, yy).RGBA()

						if uint16(rr) < r {
							r = uint16(rr)
						}
						if uint16(gg) < g {
							g = uint16(gg)
						}
						if uint16(bb) < b {
							b = uint16(bb)
						}
					}
				}
			}

			newColor := color.RGBA64{r, g, b, uint16(a)}
			outImg.Set(x, y, newColor)
		}
	}
	return outImg
}

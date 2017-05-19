package editImage

import (
	//"fmt"
	"image"
	"image/color"
)

func Dilatation(inImg image.Image, mask image.Image, origin [2]int) image.Image {
	n := inImg.Bounds().Max.Y
	m := inImg.Bounds().Max.X
	nM := mask.Bounds().Max.Y
	mM := mask.Bounds().Max.X
	outImg := image.NewRGBA(image.Rect(0, 0, m, n))

	inImg = NG(inImg, uint8(0))
	mask = NG(mask, uint8(0))

	//	Recorrer toda la imagen.
	for y := 0; y < n; y++ {
		for x := 0; x < m; x++ {
			_, _, _, a := inImg.At(x, y).RGBA()

			max := [3]int{0, 0, 0}
			//	Recorrer EE (mask) sobre la imagen.
			for yM := 0; yM < nM; yM++ {
				for xM := 0; xM < mM; xM++ {
					xTmp := (x - origin[0]) + xM
					yTmp := (y - origin[1]) + yM

					//	Comp. que pertenece a un pixel de la img.
					if IsValid(xTmp, yTmp, n, m) {
						r, g, b, _ := inImg.At(xTmp, yTmp).RGBA()
						rM, gM, bM, aM := mask.At(xM, yM).RGBA()

						if aM > 0 { //	Pixel dentro del objeto de MASK.
							sumR := int(r) + int(rM)
							sumG := int(g) + int(gM)
							sumB := int(b) + int(bM)

							//	Maximos.
							if sumR > max[0] {
								max[0] = sumR
							}
							if sumG > max[1] {
								max[1] = sumG
							}
							if sumB > max[2] {
								max[2] = sumB
							}
						}
					}
				}
			}
			//fmt.Printf("Max: %d\n", max[0])
			newColor := color.RGBA64{uint16(max[0]), uint16(max[1]), uint16(max[2]), uint16(a)}
			outImg.Set(x, y, newColor)
		}
	}

	return outImg
}

func Erosion(inImg image.Image, mask image.Image, origin [2]int) image.Image {
	n := inImg.Bounds().Max.Y
	m := inImg.Bounds().Max.X
	nM := mask.Bounds().Max.Y
	mM := mask.Bounds().Max.X
	outImg := image.NewRGBA(image.Rect(0, 0, m, n))

	inImg = NG(inImg, uint8(0))
	mask = NG(mask, uint8(0))

	//	Recorrer toda la imagen.
	for y := 0; y < n; y++ {
		for x := 0; x < m; x++ {
			_, _, _, a := inImg.At(x, y).RGBA()

			min := [3]int{0xffff, 0xffff, 0xffff}
			//	Recorrer EE (mask) sobre la imagen.
			for yM := 0; yM < nM; yM++ {
				for xM := 0; xM < mM; xM++ {
					xTmp := (x - origin[0]) + xM
					yTmp := (y - origin[1]) + yM

					//	Comp. que pertenece a un pixel de la img.
					if IsValid(xTmp, yTmp, n, m) {
						r, g, b, _ := inImg.At(xTmp, yTmp).RGBA()
						rM, gM, bM, aM := mask.At(xM, yM).RGBA()

						if aM > 0 { //	Pixel dentro del objeto de MASK.
							sumR := int(r) + int(rM)
							sumG := int(g) + int(gM)
							sumB := int(b) + int(bM)

							//	Maximos.
							if sumR < min[0] {
								min[0] = sumR
							}
							if sumG < min[1] {
								min[1] = sumG
							}
							if sumB < min[2] {
								min[2] = sumB
							}
						}
					}
				}
			}
			//fmt.Printf("Max: %d\n", min[0])
			newColor := color.RGBA64{uint16(min[0]), uint16(min[1]), uint16(min[2]), uint16(a)}
			outImg.Set(x, y, newColor)
		}
	}

	return outImg
}

/*	----------	Binario	---------	*/

func Opening(inImg image.Image, mask image.Image, origin [2]int) image.Image {
	return Dilatation(Erosion(inImg, mask, origin), mask, origin)
	//return Dilatation(inImg, mask, origin)
}

func Closing(inImg image.Image, mask image.Image, origin [2]int) image.Image {
	return Erosion(Dilatation(inImg, mask, origin), mask, origin)
	//return Erosion(inImg, mask, origin)
}

func HMT(inImg image.Image, mask image.Image, origin [2]int) image.Image {
	n := inImg.Bounds().Max.Y
	m := inImg.Bounds().Max.X
	nM := mask.Bounds().Max.Y
	mM := mask.Bounds().Max.X
	outImg := image.NewRGBA(image.Rect(0, 0, m, n))

	inImg = NG(inImg, uint8(0))
	mask = NG(mask, uint8(0))
	isBackground := false

	//	Recorrer toda la imagen.
	for y := 0; y < n; y++ {
		for x := 0; x < m; x++ {
			r, g, b, a := inImg.At(x, y).RGBA()
			isBackground = false
			//	Recorrer EE (mask) sobre la imagen.
		LoopEE:
			for yM := 0; yM < nM; yM++ {
				for xM := 0; xM < mM; xM++ {
					xTmp := (x - origin[0]) + xM
					yTmp := (y - origin[1]) + yM

					//	Comp. que pertenece a un pixel de la img.
					if IsValid(xTmp, yTmp, n, m) {
						r, g, b, _ = inImg.At(xTmp, yTmp).RGBA()
						rM, gM, bM, aM := mask.At(xM, yM).RGBA()

						if aM > 0 { //	Pixel dentro del objeto de MASK.
							//	Si no coinciden todos, es fondo.
							if r != rM || g != gM || b != bM {
								isBackground = true
								break LoopEE
							}
						}
					}
				}
			}
			if isBackground {
				r = 0
				g = 0
				b = 0
			} else {
				r = 0xffff
				g = 0xffff
				b = 0xffff
			}
			//fmt.Printf("Max: %d\n", max[0])
			newColor := color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)}
			outImg.Set(x, y, newColor)
		}
	}

	return outImg
}

/*	----------	Latices	---------	*/
func MorphSmoothing(inImg image.Image, mask image.Image, origin [2]int) image.Image {
	return Closing(Opening(inImg, mask, origin), mask, origin)
}

func difference(aImg image.Image, bImg image.Image) image.Image {
	n := aImg.Bounds().Max.Y
	m := aImg.Bounds().Max.X
	outImg := image.NewRGBA(image.Rect(0, 0, m, n))

	//	Recorrer toda la imagen.
	for y := 0; y < n; y++ {
		for x := 0; x < m; x++ {
			r, g, b, a := aImg.At(x, y).RGBA()
			rB, gB, bB, _ := bImg.At(x, y).RGBA()

			//	Diferencia entre f y aux(f)
			r -= rB
			g -= gB
			b -= bB

			newColor := color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)}
			outImg.Set(x, y, newColor)
		}
	}
	return outImg
}

func ErosionGradient(inImg image.Image, mask image.Image, origin [2]int) image.Image {
	auxImg := Erosion(inImg, mask, origin)

	return difference(inImg, auxImg)
}

func DilatationGradient(inImg image.Image, mask image.Image, origin [2]int) image.Image {
	auxImg := Dilatation(inImg, mask, origin)

	return difference(auxImg, inImg)
}

func TopHat(inImg image.Image, mask image.Image, origin [2]int) image.Image {
	auxImg := Opening(inImg, mask, origin)

	return difference(inImg, auxImg)
}

func BotHat(inImg image.Image, mask image.Image, origin [2]int) image.Image {
	auxImg := Closing(inImg, mask, origin)

	return difference(auxImg, inImg)
}

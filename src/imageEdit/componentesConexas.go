package imageEdit

import (
	//	"fmt"
	"image"
	"image/color"
)

func CompConex(inImg image.Image, conect int64) image.Image {
	//	fmt.Printf("conect: %d\n", conect)
	binImg := Binarizing(inImg, uint16(0x7FFF))
	outImg := image.NewRGBA64(image.Rect(0, 0, inImg.Bounds().Max.X, inImg.Bounds().Max.Y))

	xm := binImg.Bounds().Max.X
	ym := binImg.Bounds().Max.Y

	nEtiq := uint32(1)
	var tiketEq [10000]uint32
	var newColors [10000]color.Color

	// Primera ronda
	for y := 0; y < ym; y++ {
		for x := 0; x < xm; x++ {
			/*	--------------------	Process	--------------------	*/
			px, _, _, a := binImg.At(x, y).RGBA()

			if px != 0 {
				if conect == 4 {
					ok1, a1 := IsOne(x-1, y, ym, xm, binImg) // Left
					ok2, a2 := IsOne(x, y-1, ym, xm, binImg) // Up

					if ok1 {
						//	Left.
						val1, _, _, _ := outImg.At(x-1, y).RGBA()
						outImg.Set(x, y, color.RGBA64{uint16(val1), uint16(val1), uint16(val1), uint16(a1)})

						//	Left-Up.
						if ok2 {
							// Igualar etiqueras.
							val2, _, _, _ := outImg.At(x, y-1).RGBA()

							i := val1
							for ; tiketEq[i] != 0; i = tiketEq[i] {
							}
							j := val2
							for ; tiketEq[j] != 0; j = tiketEq[j] {
							}

							if i < j {
								tiketEq[j] = i
							} else if i > j {
								tiketEq[i] = j
							}
						}
					} else if ok2 {
						//	Up.
						val2, _, _, _ := outImg.At(x, y-1).RGBA()
						outImg.Set(x, y, color.RGBA64{uint16(val2), uint16(val2), uint16(val2), uint16(a2)})
					} else {
						//	Nothing.
						outImg.Set(x, y, color.RGBA64{uint16(nEtiq), uint16(nEtiq), uint16(nEtiq), uint16(a)})
						nEtiq++
					}
				} else {
					ok1, a1 := IsOne(x-1, y-1, ym, xm, binImg) // Left-Up
					ok2, a2 := IsOne(x+1, y-1, ym, xm, binImg) // Right-Up

					if ok1 {
						//	Left-Up.
						val1, _, _, _ := outImg.At(x-1, y-1).RGBA()
						outImg.Set(x, y, color.RGBA64{uint16(val1), uint16(val1), uint16(val1), uint16(a1)})

						//	Left-Up/Right-Up.
						if ok2 {
							// Igualar etiqueras.
							val2, _, _, _ := outImg.At(x+1, y-1).RGBA()

							i := val1
							for ; tiketEq[i] != 0; i = tiketEq[i] {
							}
							j := val2
							for ; tiketEq[j] != 0; j = tiketEq[j] {
							}

							if i < j {
								tiketEq[j] = i
							} else if i > j {
								tiketEq[i] = j
							}
						}
					} else if ok2 {
						//	Right-Up.
						val2, _, _, _ := outImg.At(x+1, y-1).RGBA()
						outImg.Set(x, y, color.RGBA64{uint16(val2), uint16(val2), uint16(val2), uint16(a2)})
					} else {
						//	Nothing.
						outImg.Set(x, y, color.RGBA64{uint16(nEtiq), uint16(nEtiq), uint16(nEtiq), uint16(a)})
						nEtiq++
					}
				}
			} else {
				outImg.Set(x, y, color.RGBA64{uint16(px), uint16(px), uint16(px), uint16(a)})
			}
		}
	}

	// Generate nEtiq numbers of colors
	//	fmt.Printf("tikets=%d\n", nEtiq)
	for i := uint32(0); i < nEtiq; i++ {
		newColors[i] = genColorRGBA64()
	}

	//***Ronda 2
	for y := 0; y < ym; y++ {
		for x := 0; x < xm; x++ {
			px, _, _, a := outImg.At(x, y).RGBA()
			//			fmt.Printf("%d ", px)
			if px != 0 {
				i := px
				for tiketEq[i] != 0 {
					i = tiketEq[i]
				}

				r, g, b, _ := newColors[i-1].RGBA()
				outImg.Set(x, y, color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)})
			}
		}
		//		fmt.Println("")
	}

	//color.GrayModel.Convert(oldColor)

	return outImg
	//return binImg
}

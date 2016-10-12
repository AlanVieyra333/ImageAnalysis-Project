package imageEdit

import (
	"image"
	"image/color"
)

func Sum(inImg image.Image, mask image.Image) image.Image {
	xm := inImg.Bounds().Max.X
	ym := inImg.Bounds().Max.Y
	outImg := image.NewRGBA64(image.Rect(0, 0, xm, ym))
	xm2 := mask.Bounds().Max.X
	ym2 := mask.Bounds().Max.Y

	pxMax := uint32((1 << 16) - 1)

	for x := 0; x < xm; x++ {
		for y := 0; y < ym; y++ {
			/*	--------------------	Process	--------------------	*/
			r, g, b, a := inImg.At(x, y).RGBA()

			if IsValid(x, y, ym2, xm2) {
				r2, g2, b2, _ := mask.At(x, y).RGBA()
				if r+r2 > pxMax {
					r = pxMax
				} else {
					r = r + r2
				}

				if g+g2 > pxMax {
					g = pxMax
				} else {
					g = g + g2
				}

				if b+b2 > pxMax {
					b = pxMax
				} else {
					b = b + b2
				}
			}

			newColor := color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)}
			outImg.Set(x, y, newColor)
		}
	}

	return outImg
}

func Sub(inImg image.Image, mask image.Image) image.Image {
	xm := inImg.Bounds().Max.X
	ym := inImg.Bounds().Max.Y
	outImg := image.NewRGBA64(image.Rect(0, 0, xm, ym))
	xm2 := mask.Bounds().Max.X
	ym2 := mask.Bounds().Max.Y

	pxMin := uint32(0)

	for x := 0; x < xm; x++ {
		for y := 0; y < ym; y++ {
			/*	--------------------	Process	--------------------	*/
			r, g, b, a := inImg.At(x, y).RGBA()

			if IsValid(x, y, ym2, xm2) {
				r2, g2, b2, _ := mask.At(x, y).RGBA()
				if r-r2 < pxMin {
					r = pxMin
				} else {
					r = r - r2
				}

				if g-g2 < pxMin {
					g = pxMin
				} else {
					g = g - g2
				}

				if b-b2 < pxMin {
					b = pxMin
				} else {
					b = b - b2
				}
			}

			newColor := color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)}
			outImg.Set(x, y, newColor)
		}
	}

	return outImg
}

func And(inImg image.Image, mask image.Image) image.Image {
	xm := inImg.Bounds().Max.X
	ym := inImg.Bounds().Max.Y
	outImg := image.NewRGBA64(image.Rect(0, 0, xm, ym))
	xm2 := mask.Bounds().Max.X
	ym2 := mask.Bounds().Max.Y

	inImgNG := NG(inImg, uint8(0))
	maskNG := NG(mask, uint8(0))

	for x := 0; x < xm; x++ {
		for y := 0; y < ym; y++ {
			/*	--------------------	Process	--------------------	*/
			r, g, b, a := inImgNG.At(x, y).RGBA()

			if IsValid(x, y, ym2, xm2) {
				r2, _, _, _ := maskNG.At(x, y).RGBA()
				var rNew uint32
				rNew = 0
				for i := uint(15); i >= 0; i-- {
					rNew = (rNew << 1) | (((r >> i) & (r2 >> i)) & 1)
				}
				r = rNew
			}
			g = r
			b = r

			newColor := color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)}
			outImg.Set(x, y, newColor)
		}
	}

	return outImg
}

func Or(inImg image.Image, mask image.Image) image.Image {
	xm := inImg.Bounds().Max.X
	ym := inImg.Bounds().Max.Y
	outImg := image.NewRGBA64(image.Rect(0, 0, xm, ym))
	xm2 := mask.Bounds().Max.X
	ym2 := mask.Bounds().Max.Y

	inImgNG := NG(inImg, uint8(0))
	maskNG := NG(mask, uint8(0))

	for x := 0; x < xm; x++ {
		for y := 0; y < ym; y++ {
			/*	--------------------	Process	--------------------	*/
			r, g, b, a := inImgNG.At(x, y).RGBA()

			if IsValid(x, y, ym2, xm2) {
				r2, _, _, _ := maskNG.At(x, y).RGBA()
				var rNew uint32
				rNew = 0
				for i := uint(15); i >= 0; i-- {
					rNew = (rNew << 1) | (((r >> i) | (r2 >> i)) & 1)
				}
				r = rNew
			}
			g = r
			b = r

			newColor := color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)}
			outImg.Set(x, y, newColor)
		}
	}

	return outImg
}

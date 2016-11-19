package editImage

import (
	"image"
	"image/color"
)

func Negative(inImg image.Image) image.Image {
	outImg := image.NewRGBA64(image.Rect(0, 0, inImg.Bounds().Max.X, inImg.Bounds().Max.Y))

	for x := 0; x < inImg.Bounds().Max.X; x++ {
		for y := 0; y < inImg.Bounds().Max.Y; y++ {
			/*	--------------------	Process	--------------------	*/
			r, g, b, a := inImg.At(x, y).RGBA()

			rNew := uint16((0xFFFF) - r)
			gNew := uint16((0xFFFF) - g)
			bNew := uint16((0xFFFF) - b)
			aNew := uint16(a)

			newColor := color.RGBA64{rNew, gNew, bNew, aNew}
			outImg.Set(x, y, newColor)
		}
	}

	return outImg
}

package imageEdit

import (
	"image"
	"image/color"
)

func NG(inImg image.Image) image.Image {
	outImg := image.NewRGBA64(image.Rect(0, 0, inImg.Bounds().Max.X, inImg.Bounds().Max.Y))

	for x := 0; x < inImg.Bounds().Max.X; x++ {
		for y := 0; y < inImg.Bounds().Max.Y; y++ {
			/*	--------------------	Process	--------------------	*/
			r, g, b, a := inImg.At(x, y).RGBA()

			grayColor := uint16(GetGrayWeighted(r, g, b))

			newColor := color.RGBA64{grayColor, grayColor, grayColor, uint16(a)}
			outImg.Set(x, y, newColor)
		}
	}

	return outImg
}

func GetGrayWeighted(R uint32, G uint32, B uint32) uint32 {
	var result uint32
	result = 0
	result += uint32(float64(R) * 0.3)
	result += uint32(float64(R) * 0.59)
	result += uint32(float64(R) * 0.11)

	return result
}

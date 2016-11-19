package editImage

import (
	"image"
	"image/color"
)

func Binarizing(inImg image.Image, umb uint16, channel uint8) *image.RGBA {
	//	fmt.Printf("umb: %d\n", umb)
	outImg := image.NewRGBA(image.Rect(0, 0, inImg.Bounds().Max.X, inImg.Bounds().Max.Y))

	for x := 0; x < inImg.Bounds().Max.X; x++ {
		for y := 0; y < inImg.Bounds().Max.Y; y++ {
			/*	--------------------	Process	--------------------	*/
			r, g, b, a := inImg.At(x, y).RGBA()

			grayColor := uint16(GetGrayWeighted(r, g, b))

			if grayColor > umb {
				grayColor = uint16(0xffff)
			} else {
				grayColor = 0
			}

			var newColor color.Color
			switch channel {
			case 1:
				newColor = color.RGBA64{grayColor, uint16(0), uint16(0), uint16(a)}
			case 2:
				newColor = color.RGBA64{uint16(0), grayColor, uint16(0), uint16(a)}
			case 3:
				newColor = color.RGBA64{uint16(0), uint16(0), grayColor, uint16(a)}
			default:
				newColor = color.RGBA64{grayColor, grayColor, grayColor, uint16(a)}
			}

			outImg.Set(x, y, newColor)
		}
	}

	return outImg
}

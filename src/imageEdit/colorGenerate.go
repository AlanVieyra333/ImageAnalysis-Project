package imageEdit

import (
	//	"fmt"
	"image/color"
	"math/rand"
)

func genColorRGBA64() color.Color {
	r := uint16(rand.Intn((1 << 16) - 1))
	g := uint16(rand.Intn((1 << 16) - 1))
	b := uint16(rand.Intn((1 << 16) - 1))
	a := uint16(rand.Intn((1 << 16) - 1))
	newColor := color.RGBA64{r, g, b, a}

	return newColor
}

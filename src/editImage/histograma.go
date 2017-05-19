/*
http://sabia.tic.udc.es/gc/Contenidos%20adicionales/trabajos/Imagenyvideo/Procesado%20Digital%20de%20la%20Imagen/pagina_superior5.htm
*/
package editImage

import (
	"image"
	"image/color"
	"math"
)

func Histograma(inImg image.Image) ([0xff]float64, float64, float64, float64, float64, float64) {
	n := inImg.Bounds().Max.Y
	m := inImg.Bounds().Max.X
	inImg = NG(inImg, 0)

	M := m * n
	var grayLevel [0xff]int

	for y := 0; y < n; y++ {
		for x := 0; x < m; x++ {
			r, _, _, _ := inImg.At(x, y).RGBA()
			grayLevel[int((float32(r)/float32(0xffff))*float32(255))]++
		}
	}

	var Pr [0xff]float64
	for k := 0; k < len(Pr); k++ {
		Pr[k] = float64(grayLevel[k]) * 100.0 / float64(M)
	}

	// Propiedades.
	media := HistMedia(Pr)
	varianza := HistVarianza(Pr, media)
	asimetria := HistAsimetria(Pr, media)
	energia := HistEnergia(Pr)
	entropia := HistEntropia(Pr)

	//return Pr, media, varianza, asimetria, energia, entropia
	return Pr, media, varianza, asimetria, energia, entropia
	//return Pr, 1.0, 2.0, 3.0, 4.0, 5.0
}

func HistMedia(P [0xff]float64) float64 {
	var res float64

	for g := 0; g < len(P); g++ {
		res += float64(g) * P[g]
	}

	return res
}

func HistVarianza(P [0xff]float64, media float64) float64 {
	var res float64

	for g := 0; g < len(P); g++ {
		res += (float64(g) - media) * (float64(g) - media) * P[g]
	}

	return res
}

func HistAsimetria(P [0xff]float64, media float64) float64 {
	var res float64

	for g := 0; g < len(P); g++ {
		res += (float64(g) - media) * (float64(g) - media) * (float64(g) - media) * P[g]
	}

	return res
}

func HistEnergia(P [0xff]float64) float64 {
	var res float64

	for g := 0; g < len(P); g++ {
		res += P[g] * P[g]
	}

	return res
}

func HistEntropia(P [0xff]float64) float64 {
	var res float64
	res = 0.0
	for g := 0; g < len(P); g++ {
		if P[g] != 0.0 {
			res -= P[g] * math.Log2(P[g])
		}
	}

	return res
}

func Ecualization(inImg image.Image) image.Image {
	n := inImg.Bounds().Max.Y
	m := inImg.Bounds().Max.X
	outImg := image.NewRGBA64(image.Rect(0, 0, m, n))
	inImg = NG(inImg, 0)

	Pr, _, _, _, _, _ := Histograma(inImg)

	for y := 0; y < n; y++ {
		for x := 0; x < m; x++ {
			k, _, _, a := inImg.At(x, y).RGBA()

			Sk := 0.0
			for j := 0; j <= int(k*255/0xffff); j++ {
				Sk += Pr[j]
			}

			px := uint16(Sk / 100.0 * 0xffff)
			if px > 0xffff {
				px = 0xffff
			}

			newColor := color.RGBA64{px, px, px, uint16(a)}
			outImg.Set(x, y, newColor)
		}
	}

	return outImg
}

func Contraction(inImg image.Image, valMax uint32, valMin uint32) image.Image {
	n := inImg.Bounds().Max.Y
	m := inImg.Bounds().Max.X
	outImg := image.NewRGBA64(image.Rect(0, 0, m, n))
	inImg = NG(inImg, 0)

	fMax := uint32(0)
	fMin := uint32(0xffff)

	for y := 0; y < n; y++ {
		for x := 0; x < m; x++ {
			r, _, _, _ := inImg.At(x, y).RGBA()
			if r > fMax {
				fMax = r
			}
			if r < fMin {
				fMin = r
			}
		}
	}

	for y := 0; y < n; y++ {
		for x := 0; x < m; x++ {
			f, _, _, a := inImg.At(x, y).RGBA()

			g := uint16(float64(float64(valMax-valMin)/float64(fMax-fMin))*float64(f-fMin) + float64(valMin))

			newColor := color.RGBA64{g, g, g, uint16(a)}
			outImg.Set(x, y, newColor)
		}
	}

	return outImg
}

func Expansion(inImg image.Image, valMax uint32, valMin uint32) image.Image {
	n := inImg.Bounds().Max.Y
	m := inImg.Bounds().Max.X
	outImg := image.NewRGBA64(image.Rect(0, 0, m, n))
	inImg = NG(inImg, 0)

	fMax := uint32(0)
	fMin := uint32(0xffff)

	for y := 0; y < n; y++ {
		for x := 0; x < m; x++ {
			r, _, _, _ := inImg.At(x, y).RGBA()
			if r > fMax {
				fMax = r
			}
			if r < fMin {
				fMin = r
			}
		}
	}

	for y := 0; y < n; y++ {
		for x := 0; x < m; x++ {
			f, _, _, a := inImg.At(x, y).RGBA()

			g := uint16((float64(f-fMin)/float64(fMax-fMin))*float64(valMax-valMin) + float64(valMin))

			newColor := color.RGBA64{g, g, g, uint16(a)}
			outImg.Set(x, y, newColor)
		}
	}

	return outImg
}

func Displacement(inImg image.Image, valDes int) image.Image {
	n := inImg.Bounds().Max.Y
	m := inImg.Bounds().Max.X
	outImg := image.NewRGBA64(image.Rect(0, 0, m, n))
	inImg = NG(inImg, 0)

	for y := 0; y < n; y++ {
		for x := 0; x < m; x++ {
			f, _, _, a := inImg.At(x, y).RGBA()

			valNew := int(f) + valDes

			/*if f == 3598 {
				fmt.Printf("valNew: %d, %d, %d\n", valNew, f, valDes)
			}*/

			if valNew > 0xffff {
				valNew = 0xffff
			}
			if valNew < 0 {
				valNew = 0
			}

			g := uint16(valNew)

			newColor := color.RGBA64{g, g, g, uint16(a)}
			outImg.Set(x, y, newColor)
		}
	}

	return outImg
}

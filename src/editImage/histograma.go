/*
http://sabia.tic.udc.es/gc/Contenidos%20adicionales/trabajos/Imagenyvideo/Procesado%20Digital%20de%20la%20Imagen/pagina_superior5.htm
*/
package editImage

import (
	"image"
	"math"
)

func Histograma(inImg image.Image) ([0xff]float64, float64, float64, float64, float64, float64) {
	n := inImg.Bounds().Max.Y
	m := inImg.Bounds().Max.X

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
		Pr[k] = float64(grayLevel[k]) / float64(M)
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

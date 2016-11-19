package editImage

import (
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"strconv"
	"strings"
)

type InfoJSON struct {
	Operation    int
	FileName     string
	FileNameEdit string
	Args         string
}

type DataOutJSON struct {
	Data1 [0xff]float64
	Data2 [5]float64
}

func Edit(data InfoJSON, newFile string, folder string) (DataOutJSON, error) {
	/*	--------------------	Args	--------------------	*/
	args := strings.Split(data.Args, ";")
	var dataJSON DataOutJSON
	//	len(args)
	//fmt.Println("holaaaaaaaa: --------" + data.Args)
	/*	--------------------	Open image	--------------------	*/
	var outImg image.Image
	inFile, err := os.Open(folder + data.FileNameEdit)
	if err != nil {
		return dataJSON, err
	}
	//	fmt.Println("open")
	defer inFile.Close()

	inImg, imageType, err := image.Decode(inFile)
	//fmt.Println("imageType: " + imageType)
	if err != nil {
		return dataJSON, err
	}
	fmt.Printf("operation: %d\n", data.Operation)

	/*	--------------------	Operations	--------------------	*/
	switch data.Operation {
	case 1:
		aux, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return dataJSON, err
		}
		outImg = CompConex(inImg, aux)
	case 2:
		aux, err := strconv.ParseFloat(args[0], 64)
		if err != nil {
			return dataJSON, err
		}
		outImg = Binarizing(inImg, uint16(aux*((1<<16)-1)), uint8(0))
	case 3:
		outImg = NG(inImg, uint8(0))
	case 4:
		outImg = Negative(inImg)
	case 5, 6, 7, 8:
		tmpFile, err := os.Open(folder + args[0])
		if err != nil {
			return dataJSON, err
		}
		defer tmpFile.Close()

		tmpImg, _, err := image.Decode(tmpFile)
		if err != nil {
			return dataJSON, err
		}

		if data.Operation == 5 {
			outImg = Sum(inImg, tmpImg)
		} else if data.Operation == 6 {
			outImg = Sub(inImg, tmpImg)
		} else if data.Operation == 7 {
			outImg = And(inImg, tmpImg)
		} else if data.Operation == 8 {
			outImg = Or(inImg, tmpImg)
		}
	case 9: // Umbralizacion by channel.
		umb, err := strconv.ParseFloat(args[0], 64)
		if err != nil {
			return dataJSON, err
		}
		channel, err := strconv.ParseInt(args[1], 10, 64)
		if err != nil {
			return dataJSON, err
		}

		outImg = Binarizing(inImg, uint16(umb*((1<<16)-1)), uint8(channel))
	case 10: // Gray scale by channel.
		channel, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return dataJSON, err
		}
		outImg = NG(inImg, uint8(channel))
	case 12, 13: // Histograma & sus Propiedades
		outImg = NG(inImg, uint8(0))
		Pr, media, varianza, asimetria, energia, entropia := Histograma(outImg)
		outImg = inImg
		dataJSON.Data1 = Pr
		dataJSON.Data2 = [5]float64{media, varianza, asimetria, energia, entropia}
	/*case 14: // Desplazamiento
	case 15: // Ensanchamiento
	case 16: // Estiramiento
	case 17: // Ecualización*/

	case 18: // Filtro Laplaciano
		umb, err := strconv.ParseFloat(args[0], 64)
		if err != nil {
			return dataJSON, err
		}
		outImg = FilterLaplaciano(inImg, uint16(umb*((1<<16)-1)))
	case 19: // Filtro Robert
		umb, err := strconv.ParseFloat(args[0], 64)
		if err != nil {
			return dataJSON, err
		}
		outImg = FilterRobert(inImg, uint16(umb*((1<<16)-1)))
	case 20: // Filtro Prewitt
		umb, err := strconv.ParseFloat(args[0], 64)
		if err != nil {
			return dataJSON, err
		}
		outImg = FilterPrewitt(inImg, uint16(umb*((1<<16)-1)))
	case 21: // Filtro Sobel
		umb, err := strconv.ParseFloat(args[0], 64)
		if err != nil {
			return dataJSON, err
		}
		outImg = FilterSobel(inImg, uint16(umb*((1<<16)-1)))
	case 22: // Filtro promedio
		window, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return dataJSON, err
		}
		outImg = FilterProm(inImg, int8(window))
	case 23: // Filtro promedio pesado
		window, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return dataJSON, err
		}
		n, err := strconv.ParseInt(args[1], 10, 64)
		if err != nil {
			return dataJSON, err
		}
		outImg = FilterPromH(inImg, int8(window), uint32(n))
	case 24: // Filtro Mediana
		window, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return dataJSON, err
		}
		outImg = FilterMediana(inImg, int8(window))
	case 25: // Filtro Moda
		window, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return dataJSON, err
		}
		outImg = FilterModa(inImg, int8(window))
	case 26: // Filtro Max
		window, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return dataJSON, err
		}
		outImg = FilterMax(inImg, int8(window))
	case 27: // Filtro Min
		window, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return dataJSON, err
		}
		outImg = FilterMin(inImg, int8(window))
	default:
		outImg = NG(inImg, uint8(0))
	}

	/*	--------------------	Save new image	--------------------	*/
	if outImg != nil {
		outFile, err := os.Create(folder + newFile)
		if err != nil {
			return dataJSON, err
		}
		defer outFile.Close()

		switch imageType {
		case "jpeg":
			jpeg.Encode(outFile, image.Image(outImg), &jpeg.Options{jpeg.DefaultQuality})
		case "png":
			png.Encode(outFile, outImg)
		case "gif":
			gif.Encode(outFile, outImg, nil)
		default:
			fmt.Println("Type image: " + imageType)
			png.Encode(outFile, outImg)
		}
	}
	fmt.Println("Image edited succesful!: " + newFile)

	return dataJSON, nil
}

func IsOne(x int, y int, n int, m int, binImg *image.RGBA) (bool, uint32) {
	if IsValid(x, y, n, m) {
		r, _, _, a := binImg.At(x, y).RGBA()
		if r != 0 {
			return true, a
		} else {
			return false, 0
		}
	} else {
		return false, 0
	}
}

func IsValid(x int, y int, n int, m int) bool {
	if (x >= 0) && (y >= 0) && (x < m) && (y < n) {
		return true
	}
	return false
}

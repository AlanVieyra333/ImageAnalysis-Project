package imageEdit

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

func Edit(data InfoJSON, newFile string, folder string) error {
	/*	--------------------	Args	--------------------	*/
	args := strings.Split(data.Args, ";")
	//	len(args)
	//fmt.Println("holaaaaaaaa: --------" + data.Args)
	/*	--------------------	Open image	--------------------	*/
	var outImg image.Image
	inFile, err := os.Open(folder + data.FileNameEdit)
	if err != nil {
		return err
	}
	//	fmt.Println("open")
	defer inFile.Close()

	inImg, imageType, err := image.Decode(inFile)
	//fmt.Println("imageType: " + imageType)
	if err != nil {
		return err
	}
	fmt.Printf("operation: %d\n", data.Operation)

	/*	--------------------	Operations	--------------------	*/
	switch data.Operation {
	case 1:
		aux, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return err
		}
		outImg = CompConex(inImg, aux)
	case 2:
		aux, err := strconv.ParseFloat(args[0], 64)
		if err != nil {
			return err
		}
		outImg = Binarizing(inImg, uint16(aux*((1<<16)-1)), uint8(0))
	case 3:
		outImg = NG(inImg, uint8(0))
	case 4:
		outImg = Negative(inImg)
	case 5, 6, 7, 8:
		tmpFile, err := os.Open(folder + args[0])
		if err != nil {
			return err
		}
		defer tmpFile.Close()

		tmpImg, _, err := image.Decode(tmpFile)
		if err != nil {
			return err
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
			return err
		}
		channel, err := strconv.ParseInt(args[1], 10, 64)
		if err != nil {
			return err
		}

		outImg = Binarizing(inImg, uint16(umb*((1<<16)-1)), uint8(channel))
	case 10: // Gray scale by channel.
		channel, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return err
		}
		outImg = NG(inImg, uint8(channel))
	case 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43:
		outImg = NG(inImg, uint8(0))
	}

	/*	--------------------	Save new image	--------------------	*/
	if outImg != nil {
		outFile, err := os.Create(folder + newFile)
		if err != nil {
			return err
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
	fmt.Println("Image edited succesful!")

	return nil
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

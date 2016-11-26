package editImage

import (
	"container/heap"
	"fmt"
	"image"
	"image/color"
	"tools"
)

const MASK = -2 /* valor inicial del nivel del Umbral */
const WSHED = 0 /* valor de los pixeles que pertenecen a los watersheds */
const INIT = -1 /* valor inicial de im0 */

func GetClockNeighbor(x, y, i int) (int, int) {
	switch i {
	case 1:
		return x, y - 1
	case 2:
		return x + 1, y - 1
	case 3:
		return x + 1, y
	case 4:
		return x + 1, y + 1
	case 5:
		return x, y + 1
	case 6:
		return x - 1, y + 1
	case 7:
		return x - 1, y
	case 8:
		return x - 1, y - 1
	default:
		return x, y
	}
}

func Watershed(inImg image.Image) *image.RGBA {
	n := inImg.Bounds().Max.Y
	m := inImg.Bounds().Max.X
	var x, y int
	imgGray := NG(inImg, uint8(0))
	outImg := image.NewRGBA(image.Rect(0, 0, m, n))
	var imI = make([][]int, m)
	var imO = make([][]int, m)
	var imD = make([][]int, m)
	//	Queue.
	pixels := tools.NewPriorityQueue(0)
	heap.Init(&pixels)
	item := &tools.Item{}
	hMin := 0xffff
	hMax := 0
	//	Queue 2.
	fifo := tools.NewQueue(0)
	//	Etiquetas.
	etiqueta_actual := 0
	var actual_dist int

	/*	Inicializacion.	*/
	for x = 0; x < m; x++ {
		imI[x] = make([]int, n)
		imO[x] = make([]int, n)
		imD[x] = make([]int, n)
		for y = 0; y < n; y++ {
			val, _, _, _ := imgGray.At(x, y).RGBA()
			imI[x][y] = int(val)
			imO[x][y] = INIT

			if imI[x][y] < hMin {
				hMin = imI[x][y]
			}
			if imI[x][y] > hMax {
				hMax = imI[x][y]
			}
		}
	}
	/*	<----	*/

	/*	Ordenar pixeles.	*/
	for y = 0; y < n; y++ {
		for x = 0; x < m; x++ {
			item = &tools.Item{
				Priority: imI[x][y],
				X:        x,
				Y:        y,
			}
			//fmt.Printf("*** PUSH: %d (%d,%d)\n", item.Priority, item.X, item.Y)
			heap.Push(&pixels, item)
		}
	}
	/*	<----	*/

	/*for pixels.Len() > 0 {
		pTmp := heap.Pop(&pixels).(*tools.Item)
		fmt.Printf("Pixels: %d (%d,%d)\n", pTmp.Priority, pTmp.X, pTmp.Y)
	}//*/

	//fmt.Printf("hMin:%d, hMax:%d\n", hMin, hMax)
	//fmt.Printf("First: pixel:%d, len:%d\n", pixels.Top().Priority, pixels.Len())
	/*	Proceso de inundamiento.	*/
	for h := hMin; h <= hMax; h++ {
		//fmt.Printf("Pixel:%d\n", h)
		/*	Localizar pixeles p tal que imi(p) = h	*/
		pH := make([]*tools.Item, 0) //	Pixeles con un N.G = h
		for pixels.Len() > 0 && pixels.Top().Priority == h {
			pH = append(pH, heap.Pop(&pixels).(*tools.Item))
		}
		/*	<----	*/
		if pixels.Len() > 0 {
			h = pixels.Top().Priority - 1
		}
		//fmt.Printf("-%d\n", len(pH))

		/*	Para cada pixel p tal que imI(p) = h	*/
		for _, p := range pH {
			imO[p.X][p.Y] = MASK
			/*	Al menos un vecino ya tenga etiqueta o sea WSHED.	*/
			for i := 1; i <= 8; i++ {
				x, y = GetClockNeighbor(p.X, p.Y, i)
				if IsValid(x, y, n, m) { //	Solo vecinos.

					if imO[x][y] > 0 || imO[x][y] == WSHED {
						imD[p.X][p.Y] = 1
						fifo.Push(p)
					}
				}
			}
		}
		/*	<----	*/
		actual_dist = 1
		item = &tools.Item{ //	Pixel ficticio.
			Priority: -1,
		}
		fifo.Push(item)

		for {
			p := fifo.Pop()
			if p.Priority == -1 { //	Es ficticio.
				if fifo.Len() == 0 {
					break //	El p no tiene vecinos etiquetados o WSHED.
				} else {
					item = &tools.Item{ //	Pixel ficticio.
						Priority: -1,
					}
					fifo.Push(item)
					actual_dist++
					p = fifo.Pop()
				}
			}

			/*	Para todos los vecinos de p.	*/
			for i := 1; i <= 8; i++ {
				x, y = GetClockNeighbor(p.X, p.Y, i)
				if IsValid(x, y, n, m) { //	Solo vecinos.

					if imD[x][y] < actual_dist && (imO[x][y] > 0 || imO[x][y] == WSHED) {
						/* Este vecino ya pertenece a un basin o al watersheds. */
						if imO[x][y] > 0 {
							if imO[p.X][p.Y] == MASK || imO[p.X][p.Y] == WSHED {
								imO[p.X][p.Y] = imO[x][y]
							} else if imO[p.X][p.Y] != imO[x][y] {
								//	Es WSHED porque hay encuentro de 2 etiquetas.
								imO[p.X][p.Y] = WSHED
							}
						} else if imO[p.X][p.Y] == MASK {
							imO[p.X][p.Y] = WSHED
						}
					} else if imO[x][y] == MASK && imD[x][y] == 0 {
						//	Este vecino estaba en espera (MASK).
						imD[x][y] = actual_dist + 1
						q := &tools.Item{ //	Pixel vecino.
							Priority: imI[x][y],
							X:        x,
							Y:        y,
						}
						fifo.Push(q)
					}
				}
			}
			/*	<----	*/
		}

		/*fmt.Println("-------------	FST imO	-------------------")
		printArray(imO, m, n)
		fmt.Println("-------------	imD	-------------------")
		printArray(imD, m, n) //*/

		/*	Checar si nuevo minimo ha sido descubierto.	*/
		/*	Para cada pixel p tal que imi(p) = h	*/
		for _, p := range pH {
			//fmt.Printf("2Pixel: %d\n", p.Priority)
			imD[p.X][p.Y] = 0 //	Inicializar distancia.
			if imO[p.X][p.Y] == MASK {
				//fmt.Printf("NewMASK: (%d,%d)\n", p.X, p.Y)
				etiqueta_actual++
				fifo.Push(p)
				imO[p.X][p.Y] = etiqueta_actual
				for fifo.Len() > 0 {
					q := fifo.Pop()
					//fmt.Printf("Px: %d,%d\n", q.X, q.Y)
					/*	Añadir a Queue vecinos MASK y darles etiqueta.	*/
					for i := 1; i <= 8; i++ {
						x, y = GetClockNeighbor(q.X, q.Y, i)
						if IsValid(x, y, n, m) { //	Solo vecinos.
							//fmt.Printf("Vecino: %d,%d\n", x, y)
							if imO[x][y] == MASK {
								u := &tools.Item{ //	Pixel vecino.
									Priority: imI[x][y],
									X:        x,
									Y:        y,
								}
								fifo.Push(u)
								imO[x][y] = etiqueta_actual
							}
						}
					}
					/*	<----	*/
				}
			}
		}
		/*	<----	*/
		/*	<----	*/

		//fmt.Println("-------------	SND imO	-------------------")
		//printArray(imO, m, n)
		//fmt.Println("-------------	imD	-------------------")
		//printArray(imD, m, n)

	}
	/*	<----	*/

	//	Colores.
	var newColors = make([]color.Color, etiqueta_actual)
	/*	Generar  colores.	*/
	//fmt.Printf("tikets=%d\n", etiqueta_actual)
	for i := 0; i < etiqueta_actual; i++ {
		newColors[i] = genColorRGBA64()
	}

	for y = 0; y < n; y++ {
		for x = 0; x < m; x++ {
			_, _, _, a := imgGray.At(x, y).RGBA()
			var r, g, b uint32

			if imO[x][y] > 0 {
				r, g, b, _ = newColors[imO[x][y]-1].RGBA()
				/*}else if imO[x][y] < 0 {
				fmt.Printf("%d (%d,%d)\n", imO[x][y], x, y)
				//*/
				/*r = 0xffff
				g = 0xffff
				b = 0xffff //*/
			} else {
				if imO[x][y] < 0 {
					fmt.Printf("%d (%d,%d)\n", imO[x][y], x, y)
				}
				r = 0
				g = 0
				b = 0
			}

			outImg.Set(x, y, color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)})
		}
	}

	return outImg
}

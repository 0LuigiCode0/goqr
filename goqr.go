package goqr

import (
	"errors"
	"image"
	"image/color"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"strings"
)

var maxDataM = []int{
	128, 224, 352, 512, 688, 864, 992, 1232, 1456, 1728,
	2032, 2320, 2672, 2920, 3320, 3624, 4056, 4504, 5016, 5352,
	5712, 6256, 6880, 7312, 8000, 8496, 9024, 9544, 10136, 10984,
	11640, 12328, 13048, 13800, 14496, 15312, 15936, 16816, 17728, 18672,
}
var maxDataH = []int{
	72, 128, 208, 288, 368, 480, 528, 688, 800, 976,
	1120, 1264, 1440, 1576, 1784, 2024, 2264, 2504, 2728, 3080,
	3248, 3536, 3712, 4112, 4304, 4768, 5024, 5288, 5608, 5960,
	6344, 6760, 7208, 7688, 7888, 8432, 8768, 9136, 9776, 10208,
}

var blocksM = []int{
	1, 1, 1, 2, 2, 4, 4, 4, 5, 5,
	5, 8, 9, 9, 10, 10, 11, 13, 14, 16,
	17, 17, 18, 20, 21, 23, 25, 26, 28, 29,
	31, 33, 35, 37, 38, 40, 43, 45, 47, 49,
}
var blocksH = []int{
	1, 1, 2, 4, 4, 4, 5, 6, 8, 8,
	11, 11, 16, 16, 18, 16, 19, 21, 25, 25,
	25, 34, 30, 32, 35, 37, 40, 42, 45, 48,
	51, 54, 57, 60, 63, 66, 70, 74, 77, 81,
}

var byteCorectM = []int{
	10, 16, 26, 18, 24, 16, 18, 22, 22, 26,
	30, 22, 22, 24, 24, 28, 28, 26, 26, 26,
	26, 28, 28, 28, 28, 28, 28, 28, 28, 28,
	28, 28, 28, 28, 28, 28, 28, 28, 28, 28,
}
var byteCorectH = []int{
	17, 28, 22, 16, 22, 28, 26, 26, 24, 28,
	24, 28, 22, 24, 24, 30, 28, 28, 26, 28,
	30, 24, 30, 30, 30, 30, 30, 30, 30, 30,
	30, 30, 30, 30, 30, 30, 30, 30, 30, 30,
}

const (
	levelCorrectM = 0x5e7c
	levelCorrectH = 0x1ce7
)

var polinom = map[int][]int{
	7:  {87, 229, 146, 149, 238, 102, 21},
	10: {251, 67, 46, 61, 118, 70, 64, 94, 32, 45},
	13: {74, 152, 176, 100, 86, 100, 106, 104, 130, 218, 206, 140, 78},
	15: {8, 183, 61, 91, 202, 37, 51, 58, 58, 237, 140, 124, 5, 99, 105},
	16: {120, 104, 107, 109, 102, 161, 76, 3, 91, 191, 147, 169, 182, 194, 225, 120},
	17: {43, 139, 206, 78, 43, 239, 123, 206, 214, 147, 24, 99, 150, 39, 243, 163, 136},
	18: {215, 234, 158, 94, 184, 97, 118, 170, 79, 187, 152, 148, 252, 179, 5, 98, 96, 153},
	20: {17, 60, 79, 50, 61, 163, 26, 187, 202, 180, 221, 225, 83, 239, 156, 164, 212, 212, 188, 190},
	22: {210, 171, 247, 242, 93, 230, 14, 109, 221, 53, 200, 74, 8, 172, 98, 80, 219, 134, 160, 105, 165, 231},
	24: {229, 121, 135, 48, 211, 117, 251, 126, 159, 180, 169, 152, 192, 226, 228, 218, 111, 0, 117, 232, 87, 96, 227, 21},
	26: {173, 125, 158, 2, 103, 182, 118, 17, 145, 201, 111, 28, 165, 53, 161, 21, 245, 142, 13, 102, 48, 227, 153, 145, 218, 70},
	28: {168, 223, 200, 104, 224, 234, 108, 180, 110, 190, 195, 147, 205, 27, 232, 201, 21, 43, 245, 87, 42, 195, 212, 119, 242, 37, 9, 123},
	30: {41, 173, 145, 152, 216, 31, 179, 182, 50, 48, 110, 86, 239, 96, 222, 125, 42, 173, 226, 193, 224, 130, 156, 37, 251, 216, 238, 40, 192, 180},
}

var fieldGalua = []int{
	1, 2, 4, 8, 16, 32, 64, 128, 29, 58, 116, 232, 205, 135, 19, 38,
	76, 152, 45, 90, 180, 117, 234, 201, 143, 3, 6, 12, 24, 48, 96, 192,
	157, 39, 78, 156, 37, 74, 148, 53, 106, 212, 181, 119, 238, 193, 159, 35,
	70, 140, 5, 10, 20, 40, 80, 160, 93, 186, 105, 210, 185, 111, 222, 161,
	95, 190, 97, 194, 153, 47, 94, 188, 101, 202, 137, 15, 30, 60, 120, 240,
	253, 231, 211, 187, 107, 214, 177, 127, 254, 225, 223, 163, 91, 182, 113, 226,
	217, 175, 67, 134, 17, 34, 68, 136, 13, 26, 52, 104, 208, 189, 103, 206,
	129, 31, 62, 124, 248, 237, 199, 147, 59, 118, 236, 197, 151, 51, 102, 204,
	133, 23, 46, 92, 184, 109, 218, 169, 79, 158, 33, 66, 132, 21, 42, 84,
	168, 77, 154, 41, 82, 164, 85, 170, 73, 146, 57, 114, 228, 213, 183, 115,
	230, 209, 191, 99, 198, 145, 63, 126, 252, 229, 215, 179, 123, 246, 241, 255,
	227, 219, 171, 75, 150, 49, 98, 196, 149, 55, 110, 220, 165, 87, 174, 65,
	130, 25, 50, 100, 200, 141, 7, 14, 28, 56, 112, 224, 221, 167, 83, 166,
	81, 162, 89, 178, 121, 242, 249, 239, 195, 155, 43, 86, 172, 69, 138, 9,
	18, 36, 72, 144, 61, 122, 244, 245, 247, 243, 251, 235, 203, 139, 11, 22,
	44, 88, 176, 125, 250, 233, 207, 131, 27, 54, 108, 216, 173, 71, 142, 1,
}

var reversefieldGalua = []int{
	0, 0, 1, 25, 2, 50, 26, 198, 3, 223, 51, 238, 27, 104, 199, 75,
	4, 100, 224, 14, 52, 141, 239, 129, 28, 193, 105, 248, 200, 8, 76, 113,
	5, 138, 101, 47, 225, 36, 15, 33, 53, 147, 142, 218, 240, 18, 130, 69,
	29, 181, 194, 125, 106, 39, 249, 185, 201, 154, 9, 120, 77, 228, 114, 166,
	6, 191, 139, 98, 102, 221, 48, 253, 226, 152, 37, 179, 16, 145, 34, 136,
	54, 208, 148, 206, 143, 150, 219, 189, 241, 210, 19, 92, 131, 56, 70, 64,
	30, 66, 182, 163, 195, 72, 126, 110, 107, 58, 40, 84, 250, 133, 186, 61,
	202, 94, 155, 159, 10, 21, 121, 43, 78, 212, 229, 172, 115, 243, 167, 87,
	7, 112, 192, 247, 140, 128, 99, 13, 103, 74, 222, 237, 49, 197, 254, 24,
	227, 165, 153, 119, 38, 184, 180, 124, 17, 68, 146, 217, 35, 32, 137, 46,
	55, 63, 209, 91, 149, 188, 207, 205, 144, 135, 151, 178, 220, 252, 190, 97,
	242, 86, 211, 171, 20, 42, 93, 158, 132, 60, 57, 83, 71, 109, 65, 162,
	31, 45, 67, 216, 183, 123, 164, 118, 196, 23, 73, 236, 127, 12, 111, 246,
	108, 161, 59, 82, 41, 157, 85, 170, 251, 96, 134, 177, 187, 204, 62, 90,
	203, 89, 95, 176, 156, 169, 160, 81, 11, 245, 22, 235, 122, 117, 44, 215,
	79, 174, 213, 233, 230, 231, 173, 232, 116, 214, 244, 234, 168, 80, 88, 175,
}

var qrBlocks = []int{
	21, 25, 29, 33, 37, 41, 45, 49,
	53, 57, 61, 65, 69, 73, 77, 81,
	85, 89, 93, 97, 101, 105, 109, 113,
	117, 121, 125, 129, 133, 137, 141, 145,
	149, 153, 157, 161, 165, 169, 173, 177,
}

var codeVersion = [][]int{
	{}, {}, {}, {}, {}, {},
	{0x2, 0x1e, 0x26},
	{0x11, 0x1c, 0x38},
	{0x37, 0x18, 0x4},
	{0x29, 0x3e, 0x0},
	{0xf, 0x3a, 0x3c},
	{0xd, 0x24, 0x1a},
	// {0x,0x,0x},	101011 100000 100110
	// {0x,0x,0x},	110101 000110 100010
	// {0x,0x,0x},	010011 000010 011110
	// {0x,0x,0x},	011100 010001 011100
	// {0x,0x,0x},	111010 010101 100000
	// {0x,0x,0x},	100100 110011 100100
	// {0x,0x,0x},	000010 110111 011000
	// {0x,0x,0x},	000000 101001 111110
	// {0x,0x,0x},	100110 101101 000010
	// {0x,0x,0x},	111000 001011 000110
	// {0x,0x,0x},	011110 001111 111010
	// {0x,0x,0x},	001101 001101 100100
	// {0x,0x,0x},	101011 001001 011000
	// {0x,0x,0x},	110101 101111 011100
	// {0x,0x,0x},	010011 101011 100000
	// {0x,0x,0x},	010001 110101 000110
	// {0x,0x,0x},	110111 110001 111010
	// {0x,0x,0x},	101001 010111 111110
	// {0x,0x,0x},	001111 010011 000010
	// {0x,0x,0x},	101000 011000 101101
	// {0x,0x,0x},	001110 011100 010001
	// {0x,0x,0x},	010000 111010 010101
	// {0x,0x,0x},	110110 111110 101001
	// {0x,0x,0x},	110100 100000 001111
	// {0x,0x,0x},	010010 100100 110011
	// {0x,0x,0x},	001100 000010 110111
	// {0x,0x,0x},	101010 000110 001011
	// {0x,0x,0x},	111001 000100 010101
}

var coordAnchor = [][][]int{
	{},
	{{18, 18}},
	{{22, 22}},
	{{26, 26}},
	{{30, 30}},
	{{34, 34}},
	{{6, 22}, {22, 6}, {22, 22}, {22, 38}, {38, 22}, {38, 38}},
	{{6, 24}, {24, 6}, {24, 24}, {24, 42}, {42, 24}, {42, 42}},
	{{6, 26}, {26, 6}, {26, 26}, {26, 46}, {46, 26}, {46, 46}},
	{{6, 28}, {28, 6}, {28, 28}, {28, 50}, {50, 28}, {50, 50}},
	{{6, 30}, {30, 6}, {30, 30}, {30, 54}, {54, 30}, {54, 54}},
	{{6, 32}, {32, 6}, {32, 32}, {32, 58}, {58, 32}, {58, 58}},
}

const (
	search0  byte = 2
	search1  byte = 3
	sync0    byte = 4
	sync1    byte = 5
	mask0    byte = 6
	mask1    byte = 7
	version0 byte = 8
	version1 byte = 9
	anchor0  byte = 10
	anchor1  byte = 11
)

//QRGenerate генерирует qr
func QRGenerate(content, imagePath, qrPath string, sizeImg float64) error {
	if qrPath == "" {
		return errors.New("qrPath is nil")
	}

	maxData := &maxDataM
	blocks := &blocksM
	byteCorect := &byteCorectM
	levelCorrect := levelCorrectM
	var gachi interface{}
	var maxSizeGachi int

	if imagePath != "" {
		file, err := os.OpenFile(imagePath, os.O_RDONLY, 0777)
		if err != nil {
			return err
		}
		defer file.Close()
		buf, _ := ioutil.ReadFile(imagePath)
		switch http.DetectContentType(buf) {
		case "image/jpeg":
			gachi, err = jpeg.Decode(file)
			if err != nil {
				return err
			}
			if !strings.HasSuffix(qrPath, ".jpg") {
				return errors.New("QR not jpg")
			}
		case "image/png":
			gachi, err = png.Decode(file)
			if err != nil {
				return err
			}
			if !strings.HasSuffix(qrPath, ".png") {
				return errors.New("QR not png")
			}
		case "image/gif":
			gachi, err = gif.DecodeAll(file)
			if err != nil {
				return err
			}
			if !strings.HasSuffix(qrPath, ".gif") {
				return errors.New("QR not gif")
			}
		default:
			return errors.New("image wrong type")
		}

		maxData = &maxDataH
		blocks = &blocksH
		byteCorect = &byteCorectH
		levelCorrect = levelCorrectH
	}

	//Перевод строки в двоичную последовательность
	length, data := utfToBit(content)
	//Выбор версии QR кода и длины системных данных
	version, lenSystemData, err := howToVersion(length, maxData)
	if err != nil {
		return err
	}
	maxSizeGachi = int(math.Sqrt(float64((qrBlocks[version]*qrBlocks[version])-240-(len(coordAnchor[version])*25)-(qrBlocks[version]*2)) * sizeImg))
	if maxSizeGachi%2 == 0 {
		maxSizeGachi--
	}

	//Запись системных данных в начало массива
	data = addServicesData(content, version, lenSystemData, maxData, data)
	//Дозаполнение пустышками до необходимой длины
	addVoidData(lenSystemData, length, version, maxData, &data)
	//Пстроение блоков
	block, byteBlock, size := buildBlock(version, maxData, blocks, &data)
	//Создание байт коррекции
	countByteCorect, corectBlock, length := buildCorectBlock(version, block, length, byteCorect, &byteBlock)
	//Групирование блоков данных
	data = groupData(length, size, countByteCorect, &byteBlock, &corectBlock)

	//Рисование
	size = qrBlocks[version]
	dataImg := make([][]byte, size)
	for i := range dataImg {
		dataImg[i] = make([]byte, size)
	}
	searchPoint(&dataImg)
	syncLine(&dataImg)
	maskInfo(&dataImg, levelCorrect)
	codeVer(&dataImg, version)
	anchor(&dataImg, version)
	write(&dataImg, &data)

	//Вывод изображения
	file1, err := os.OpenFile(qrPath, os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		return err
	}

	if img2, ok := gachi.(image.Image); ok {
		if err := png.Encode(file1, paintImage(size, maxSizeGachi, &dataImg, img2)); err != nil {
			return err
		}
	} else if img2, ok := gachi.(*gif.GIF); ok {
		if err := gif.EncodeAll(file1, paintGIF(size, maxSizeGachi, &dataImg, img2)); err != nil {
			return err
		}
	} else {
		if err := png.Encode(file1, paintImage(size, maxSizeGachi, &dataImg, nil)); err != nil {
			return err
		}
	}
	return nil
}

//Перевод строки в двоичную последовательность
func utfToBit(content string) (length int, dataBit []int) {
	count := 0
	length = len(content) * 8
	dataBit = make([]int, length)
	for l, i := len(content), 0; i < l; i++ {
		mask := 0x80
		a := int(content[i])
		for mask > 0 {
			if mask&a != 0 {
				dataBit[count] = 1
			}
			mask >>= 1
			count++
		}
	}
	return
}

//Выбор версии QR кода и длины системных данных
func howToVersion(length int, maxData *[]int) (version int, lenSystemData int, err error) {
	for i := 0; i < 40; i++ {
		max := (*maxData)[i]
		if length > max {
			continue
		}
		switch {
		case i < 9:
			if length+92 > max {
				i++
				if i == 9 {
					version = i
					lenSystemData = 20
					break
				}
			}
			version = i
			lenSystemData = 12
		case i >= 9 && i < 26:
			if length+100 > max {
				i++
				if i == 26 {
					version = i
					lenSystemData = 20
					break
				}
			}
			version = i
			lenSystemData = 20
		case i >= 26 && i < 40:
			if length+100 > max {
				i++
				if i == 40 {
					err = errors.New("data's oversize")
					break
				}
			}
			version = i
			lenSystemData = 20
		}
		break
	}
	return
}

//Запись системных данных в начало массива
func addServicesData(content string, version int, lenSystemData int, maxData *[]int, dataBit []int) []int {
	newData := make([]int, (*maxData)[version])
	newData[1] = 1
	mask := 1 << (lenSystemData - 4 - 1)
	countSymbol := len(content)
	for i := 4; i < lenSystemData; i++ {
		if countSymbol&mask != 0 {
			newData[i] = 1
		}
		mask >>= 1
	}
	copy(newData[lenSystemData:], dataBit)
	return newData
}

//Дозаполнение пустышками до необходимой длины
func addVoidData(lenSystemData int, length int, version int, maxData, newData *[]int) {
	minMultiplyLenData := lenSystemData + length
	for minMultiplyLenData%8 != 0 {
		minMultiplyLenData++
	}
	f := true
	for i := minMultiplyLenData; i < (*maxData)[version]; {
		switch f {
		case true:
			(*newData)[i] = 1
			(*newData)[i+1] = 1
			(*newData)[i+2] = 1
			(*newData)[i+4] = 1
			(*newData)[i+5] = 1
			i += 8
			f = false
		case false:
			(*newData)[i+3] = 1
			(*newData)[i+7] = 1
			i += 8
			f = true
		}
	}
}

//Пстроение блоков
func buildBlock(version int, maxData, blocks, newData *[]int) (block int, byteBlock [][]int, size int) {
	length := 0
	count := 0
	block = (*blocks)[version]
	byteBlock = make([][]int, block)
	maxByte := (*maxData)[version] / 8
	size, resid := maxByte/block, maxByte%block
	for i := 0; i < block; i++ {
		if resid >= block-i {
			byteBlock[i] = make([]int, size+1)
			length += size + 1
		} else {
			byteBlock[i] = make([]int, size)
			length += size
		}
		for j := 0; j < len(byteBlock[i]); j++ {
			x := (*newData)[count] * (2 << 6)
			x += (*newData)[count+1] * (2 << 5)
			x += (*newData)[count+2] * (2 << 4)
			x += (*newData)[count+3] * (2 << 3)
			x += (*newData)[count+4] * (2 << 2)
			x += (*newData)[count+5] * (2 << 1)
			x += (*newData)[count+6] * (2 << 0)
			x += (*newData)[count+7] * (1)
			byteBlock[i][j] = x
			count += 8
		}
	}
	return
}

//Создание байт коррекции
func buildCorectBlock(version int, block int, length int, byteCorect *[]int, byteBlock *[][]int) (int, [][]int, int) {
	countByteCorect := (*byteCorect)[version]
	polinomCorect := polinom[countByteCorect]
	corectBlock := make([][]int, block)
	for i := range corectBlock {
		if len((*byteBlock)[i]) > countByteCorect {
			corectBlock[i] = make([]int, len((*byteBlock)[i]))
			length += len((*byteBlock)[i])
		} else {
			corectBlock[i] = make([]int, countByteCorect)
			length += countByteCorect
		}
		copy(corectBlock[i], (*byteBlock)[i])
		for range (*byteBlock)[i] {
			x := corectBlock[i][0]
			copy(corectBlock[i], corectBlock[i][1:])
			corectBlock[i][len(corectBlock[i])-1] = 0
			if x == 0 {
				continue
			}
			x = reversefieldGalua[x]
			for j := 0; j < countByteCorect; j++ {
				y := polinomCorect[j] + x
				if y > 254 {
					y %= 255
				}
				corectBlock[i][j] ^= fieldGalua[y]
			}
		}
	}
	return countByteCorect, corectBlock, length
}

//Групирование блоков данных
func groupData(length, size, countByteCorect int, byteBlock, corectBlock *[][]int) (data []int) {
	count := 0
	data = make([]int, length)
	for j := 0; j < size+1; j++ {
		for _, v := range *byteBlock {
			if len(v) > j {
				data[count] = v[j]
				count++
			}
		}
	}
	for j := 0; j < countByteCorect; j++ {
		for _, v := range *corectBlock {
			if v[j] != 0 {
				data[count] = v[j]
				count++
			}
		}
	}
	return
}

//Поисковые маячки
func searchPoint(img *[][]byte) {
	posx := []int{0, 0, len(*img) - 7}
	posy := []int{0, len(*img) - 7, 0}
	for k := 0; k < 3; k++ {
		x := search1
		for i := 0; i < 3; i++ {
			for j := i; j < 7-i; j++ {
				(*img)[i+posy[k]][j+posx[k]] = x
				(*img)[j+posy[k]][i+posx[k]] = x
				(*img)[6-i+posy[k]][j+posx[k]] = x
				(*img)[j+posy[k]][6-i+posx[k]] = x
			}
			(*img)[6-i+posy[k]][6-i+posx[k]] = x
			if x == search0 {
				x = search1
			} else {
				x = search0
			}
		}
		(*img)[posy[k]+3][posx[k]+3] = search1
		for i := -1; i < 8; i++ {
			if posx[k]+i > -1 && posx[k]+i < len(*img) {
				if posy[k]+7 < len(*img) {
					(*img)[posy[k]+7][posx[k]+i] = x
				} else if posy[k]-1 > -1 {
					(*img)[posy[k]-1][posx[k]+i] = x
				}
			}
			if posy[k]+i > -1 && posy[k]+i < len(*img) {
				if posx[k]+7 < len(*img) {
					(*img)[posy[k]+i][posx[k]+7] = x
				} else if posx[k]-1 > -1 {
					(*img)[posy[k]+i][posx[k]-1] = x
				}
			}
		}
	}
}

//Полосы синхранизации
func syncLine(img *[][]byte) {
	f := true
	for i := 8; i < len(*img)-8; i++ {
		if f {
			(*img)[6][i] = sync1
			(*img)[i][6] = sync1
			f = false
		} else {
			(*img)[6][i] = sync0
			(*img)[i][6] = sync0
			f = true
		}
	}
}

//Информация о версии и маске
func maskInfo(img *[][]byte, code int) {
	var a, b int
	mask := 0x4000
	for i := 0; i < 15; i++ {
		if i > 6 {
			a = 8
			b = len(*img) - 15 + i

		} else {
			a = len(*img) - 1 - i
			b = 8
		}
		if mask&code != 0 {
			(*img)[a][b] = mask1
		} else {
			(*img)[a][b] = mask0
		}
		mask >>= 1
	}
	(*img)[len(*img)-8][8] = mask1

	mask = 0x4000
	for i := 0; i < 17; i++ {
		if i > 8 {
			a = 16 - i
			b = 8
		} else {
			a = 8
			b = i
		}
		if (*img)[a][b] != sync1 {
			if mask&code != 0 {
				(*img)[a][b] = mask1
			} else {
				(*img)[a][b] = mask0
			}
		} else {
			continue
		}
		mask >>= 1
	}
}

func codeVer(img *[][]byte, version int) {
	size := len(*img) - 1
	vers := codeVersion[version]
	for i := range vers {
		mask := 0x20
		for j := 0; mask > 0; j++ {
			if mask&vers[i] == 0 {
				(*img)[size-10+i][j] = version0
				(*img)[j][size-10+i] = version0
			} else {
				(*img)[size-10+i][j] = version1
				(*img)[j][size-10+i] = version1
			}
			mask >>= 1
		}
	}
}

func anchor(img *[][]byte, version int) {
	coordLisn := coordAnchor[version]
	for i := range coordLisn {
		y, x := coordLisn[i][0]-2, coordLisn[i][1]-2
		k := anchor1
		for d := 0; d < 2; d++ {
			for j := d; j < 5-d; j++ {
				(*img)[d+y][j+x] = k
				(*img)[j+y][d+x] = k
				(*img)[4-d+y][j+x] = k
				(*img)[j+y][4-d+x] = k
			}
			if k == anchor0 {
				k = anchor1
			} else {
				k = anchor0
			}
		}
		(*img)[y+2][x+2] = anchor1
	}
}

func write(img *[][]byte, data *[]int) {
	var i int
	var direct bool
	bufBit := make([]byte, len(*data)*8)
	for _, v := range *data {
		mask := 0x80
		for mask > 0 {
			if mask&v != 0 {
				bufBit[i] = 1
			}
			i++
			mask >>= 1
		}
	}
	i = 0

	for x := len(*img) - 1; x > -1; {
		if x != 6 {
			if direct {
				for y := 0; y < len(*img); y++ {
					for k := 0; k < 2; k++ {
						if (*img)[y][x-k] == 0 {
							var a byte
							if len(bufBit) > i {
								a = bufBit[i]
								i++
							}
							if (x-k)%3 == 0 {
								if a == 0 {
									(*img)[y][x-k] = 1
								} else {
									(*img)[y][x-k] = 0
								}
							} else {
								(*img)[y][x-k] = a
							}

						}
					}
				}
				direct = false
			} else {
				for y := len(*img) - 1; y > -1; y-- {
					for k := 0; k < 2; k++ {
						if (*img)[y][x-k] == 0 {
							var a byte
							if len(bufBit) > i {
								a = bufBit[i]
								i++
							}
							if (x-k)%3 == 0 {
								if a == 0 {
									(*img)[y][x-k] = 1
								} else {
									(*img)[y][x-k] = 0
								}
							} else {
								(*img)[y][x-k] = a
							}
						}
					}
				}
				direct = true
			}
			x -= 2
		} else {
			x--
		}
	}
}

func paintImage(size, maxSizeImg int, dataImg *[][]byte, image2 image.Image) *image.CMYK {
	var sizeImg, shift int
	coeff := 1
	if image2 != nil {
		coeff = (image2.Bounds().Dx() / maxSizeImg) + 1
		size = coeff * (size + 8)
		sizeImg = image2.Bounds().Dx() - 1
		maxSizeImg = maxSizeImg * coeff
		shift = ((maxSizeImg - sizeImg) / 2)
	} else {
		size = size + 8
	}

	rect := image.Rect(0, 0, size, size)
	image1 := image.NewCMYK(rect)

	x0 := (size / 2) - (maxSizeImg / 2)
	y0 := (size / 2) - (maxSizeImg / 2)

	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			if y > coeff*4-1 && y < size-(coeff*4) && x > coeff*4-1 && x < size-(coeff*4) {
				if (*dataImg)[(y-(coeff*4))/coeff][(x-(coeff*4))/coeff]%2 == 1 {
					image1.Set(x, y, color.Black)
				}
			}
		}
	}
	if image2 != nil {
		for y := 0; y < maxSizeImg; y++ {
			for x := 0; x < maxSizeImg; x++ {
				if x > maxSizeImg-shift-1 || x < shift || y > maxSizeImg-shift-1 || y < shift {
					image1.Set(x0+x, y0+y, color.White)
				} else {
					image1.Set(x0+x, y0+y, image2.At(x-shift, y-shift))
				}
			}
		}
	}
	return image1
}

func paintGIF(size, maxSizeImg int, dataImg *[][]byte, image2 *gif.GIF) *gif.GIF {
	image1 := &gif.GIF{
		Delay:     image2.Delay,
		LoopCount: image2.LoopCount,
	}
	for _, img := range image2.Image {
		frame := paintImage(size, maxSizeImg, dataImg, img)
		pall := image.NewPaletted(frame.Rect, palette.Plan9)
		draw.FloydSteinberg.Draw(pall, frame.Rect, frame, image.Point{})

		image1.Image = append(image1.Image, pall)
	}

	return image1
}

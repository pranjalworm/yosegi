package imager

import (
	"image"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/pranjalworm/stitch/models"
)

func ReadTargetImage(imageName string, imageDivisionFactor int, targetImageChannel chan<- *models.AveragedImageData) {

	filePath := "./images/target/" + imageName + ".jpg"
	img := readFilePathAsImage(filePath)
	analysedImage := analyseAndAverageOut(img, imageName, imageDivisionFactor)

	targetImageChannel <- &analysedImage

	close(targetImageChannel)
}

func ReadCollection(collectionPath string, collectionChannel chan<- []models.AveragedImageData) {

	files, err := os.ReadDir(collectionPath)

	if err != nil {
		panic(err)
	}

	imageArr := make([]models.AveragedImageData, 0, len(files))

	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, file := range files {
		ext := strings.ToLower(filepath.Ext(file.Name()))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
			continue
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			imagePath := collectionPath + file.Name()
			image := readFilePathAsImage(imagePath)
			analysedImage := analyseAndAverageOut(image, file.Name(), 1)

			mu.Lock()
			imageArr = append(imageArr, analysedImage)
			mu.Unlock()
		}()
	}

	wg.Wait()

	collectionChannel <- imageArr
}

func analyseAndAverageOut(img image.Image, imageName string, imageDivisionFactor int) models.AveragedImageData {

	const BitFactor = 257

	bounds := img.Bounds()

	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y

	var orientation string
	if width > height {
		orientation = "landscape"
	} else {
		orientation = "portrait"
	}

	imageData := models.AveragedImageData{ImageName: imageName, Orientation: orientation, Width: width, Height: height, Pixel: make([]models.Pixel, imageDivisionFactor*imageDivisionFactor)}

	xBlockLength := width / imageDivisionFactor
	yBlockLength := height / imageDivisionFactor

	maxCol := imageDivisionFactor - 1
	maxRow := imageDivisionFactor - 1

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {

			r, g, b, a := img.At(x, y).RGBA()

			red := uint32(r / BitFactor)
			green := uint32(g / BitFactor)
			blue := uint32(b / BitFactor)
			alpha := uint32(a / BitFactor)

			row := y / yBlockLength
			col := x / xBlockLength

			if col > maxCol {
				col = maxCol
			}

			if row > maxRow {
				row = maxRow
			}

			index := imageDivisionFactor*row + col

			imageData.Pixel[index].Red += red
			imageData.Pixel[index].Green += green
			imageData.Pixel[index].Blue += blue
			imageData.Pixel[index].Alpha += alpha
		}
	}

	// calculate the average color values

	gridBlockPixelCount := xBlockLength * yBlockLength

	for i, pixel := range imageData.Pixel {
		imageData.Pixel[i].Red = pixel.Red / uint32(gridBlockPixelCount)
		imageData.Pixel[i].Green = pixel.Green / uint32(gridBlockPixelCount)
		imageData.Pixel[i].Blue = pixel.Blue / uint32(gridBlockPixelCount)
		imageData.Pixel[i].Alpha = pixel.Alpha / uint32(gridBlockPixelCount)
	}

	return imageData
}

package utils

import (
	"image"
	"os"

	"github.com/pranjalworm/stitch/models"
)

const BitFactor = 257

func AnalyseImage(imagePath string, imageDivisionFactor int) models.AveragedImageData {
	file, err := os.Open(imagePath)

	if err != nil {
		panic(err)
	}

	defer file.Close()

	img, format, err := image.Decode(file)

	if err != nil {
		panic(err)
	}

	bounds := img.Bounds()

	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y

	imageData := models.AveragedImageData{Width: width, Height: height, Format: format, GridBlockPixel: make([]models.GridBlockPixel, imageDivisionFactor*imageDivisionFactor)}

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

			imageData.GridBlockPixel[index].Pixel.Red += red
			imageData.GridBlockPixel[index].Pixel.Green += green
			imageData.GridBlockPixel[index].Pixel.Blue += blue
			imageData.GridBlockPixel[index].Pixel.Alpha += alpha
		}
	}

	// calculate the average color values
	gridBlockPixelCount := xBlockLength * yBlockLength

	for i, pixel := range imageData.GridBlockPixel {
		imageData.GridBlockPixel[i].Pixel.Red = uint32(pixel.Pixel.Red / uint32(gridBlockPixelCount))
		imageData.GridBlockPixel[i].Pixel.Green = uint32(pixel.Pixel.Green / uint32(gridBlockPixelCount))
		imageData.GridBlockPixel[i].Pixel.Blue = uint32(pixel.Pixel.Blue / uint32(gridBlockPixelCount))
		imageData.GridBlockPixel[i].Pixel.Alpha = uint32(pixel.Pixel.Alpha / uint32(gridBlockPixelCount))
	}

	return imageData
}

package utils

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"time"

	"github.com/pranjalworm/stitch/models"
)

func CreateImage(imageData models.AveragedImageData, imageDivisionFactor int) {

	width := imageData.Width
	height := imageData.Height

	xBlockLength := width / imageDivisionFactor
	yBlockLength := height / imageDivisionFactor

	maxCol := imageDivisionFactor - 1
	maxRow := imageDivisionFactor - 1

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			col := x / xBlockLength
			row := y / yBlockLength

			if col > maxCol {
				col = maxCol
			}

			if row > maxRow {
				row = maxRow
			}

			index := imageDivisionFactor*row + col

			pixelColor := imageData.GridBlockPixel[index].Pixel

			img.Set(x, y, color.RGBA{R: uint8(pixelColor.Red), G: uint8(pixelColor.Green), B: uint8(pixelColor.Blue), A: uint8(pixelColor.Alpha)})
		}
	}

	file, err := os.Create("output/output" + time.Now().String() + ".png")

	if err != nil {
		panic(err)
	}

	defer file.Close()

	err = png.Encode(file, img)

	if err != nil {
		panic(err)
	}

	fmt.Println("Image created!")

}

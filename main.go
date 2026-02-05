package main

import (
	_ "image/jpeg"
	_ "image/png"
	"os"
	"strconv"

	"github.com/pranjalworm/stitch/utils"
)

func main() {

	imageName := os.Args[1]
	divisionFactor := os.Args[2]

	num, err := strconv.Atoi(divisionFactor)
	if err != nil {
		panic(err)

	}

	imageDivisionFactor := num
	imagePath := "./images/" + imageName + ".jpg"

	// read and analyse image
	imageData := utils.AnalyseImage(imagePath, imageDivisionFactor)
	utils.CreateImage(imageData, imageDivisionFactor)
}

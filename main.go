package main

import (
	_ "image/jpeg"
	_ "image/png"
	"os"
	"strconv"

	"github.com/pranjalworm/stitch/imager"
	"github.com/pranjalworm/stitch/models"
	"github.com/pranjalworm/stitch/utils"
)

func main() {

	defer utils.TrackTime("main")()
	imageName := os.Args[1]
	divisionFactor := os.Args[2]
	collectionName := os.Args[3]

	imageDivisionFactor, err := strconv.Atoi(divisionFactor)
	if err != nil {
		panic(err)
	}

	collectionChannel := make(chan []models.AveragedImageData)

	collectionPath := "./images/collections/" + collectionName + "/"
	go imager.ReadCollection(collectionPath, collectionChannel)

	targetImageChannel := make(chan *models.AveragedImageData)

	go imager.ReadTargetImage(imageName, imageDivisionFactor, targetImageChannel)

	targetImage := <-targetImageChannel

	resultImageCollection := imager.MapCollectionToTarget(targetImage, <-collectionChannel, imageDivisionFactor)

	imager.CreateCollage(resultImageCollection, collectionPath, targetImage, imageDivisionFactor)

}

package imager

import (
	"github.com/pranjalworm/stitch/models"
)

func MapCollectionToTarget(targetImage *models.AveragedImageData, collection []models.AveragedImageData, imageDivisionFactor int) []string {

	resultImageData := make([]string, imageDivisionFactor*imageDivisionFactor)

	for targetIndex := range targetImage.Pixel {

		targetPixel := targetImage.Pixel[targetIndex]

		resultImageData[targetIndex] = findMinDiffCollectionImage(targetPixel, collection, targetImage.Orientation)

	}
	return resultImageData

}

func findMinDiffCollectionImage(targetPixel models.Pixel, collection []models.AveragedImageData, targetOrientation string) string {

	var minDiff uint32 = ^uint32(0)
	var resultCollectionImage string
	var collectionImage models.AveragedImageData

	for collectionIndex := range collection {

		collectionImage = collection[collectionIndex]
		collectionPixel := collectionImage.Pixel[0]

		if targetOrientation != collectionImage.Orientation {
			continue
		}

		redDiff := abs(targetPixel.Red - collectionPixel.Red)
		blueDiff := abs(targetPixel.Blue - collectionPixel.Blue)
		greenDiff := abs(targetPixel.Green - collectionPixel.Green)
		alphaDiff := abs(targetPixel.Alpha - collectionPixel.Alpha)
		totaldiff := redDiff + blueDiff + greenDiff + alphaDiff

		if totaldiff < minDiff {
			minDiff = totaldiff
			resultCollectionImage = collectionImage.ImageName
		}
	}

	return resultCollectionImage
}

func abs(x uint32) uint32 {
	if x < 0 {
		return -x
	}

	return x
}

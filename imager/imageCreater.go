package imager

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
	"time"

	"github.com/pranjalworm/stitch/models"
	"github.com/pranjalworm/stitch/utils"
)

func CreateCollage(resultImageCollection []string, collectionPath string, targetImage *models.AveragedImageData, imageDivisionFactor int) {

	defer utils.TrackTime("CreateImage")()

	// load collection images
	loadCollectionImageChannel := make(chan []image.Image)
	go loadImages(resultImageCollection, collectionPath, loadCollectionImageChannel)

	collectionImages := <-loadCollectionImageChannel

	// create collage
	collageWidth := targetImage.Width - (targetImage.Width % imageDivisionFactor)
	collageHeight := targetImage.Height - (targetImage.Height % imageDivisionFactor)

	cols := imageDivisionFactor
	cellW, cellH := collectionImages[0].Bounds().Dx(), collectionImages[0].Bounds().Dy()
	collage := image.NewRGBA(image.Rect(0, 0, collageWidth, collageHeight))

	for i, img := range collectionImages {
		col := i % cols
		row := i / cols
		offset := image.Pt(col*cellW, row*cellH)
		draw.Draw(collage, img.Bounds().Add(offset), img, image.Point{}, draw.Src)
	}

	file, err := os.Create("images/output/output" + time.Now().String() + ".png")

	if err != nil {
		panic(err)
	}

	defer file.Close()

	err = png.Encode(file, collage)

	if err != nil {
		panic(err)
	}

	fmt.Println("Image created!")
}

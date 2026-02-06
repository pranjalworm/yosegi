package imager

import (
	"fmt"
	"image"
	"os"
	"sync"
)

func readFilePathAsImage(filePath string) image.Image {
	file, err := os.Open(filePath)

	if err != nil {
		panic(err)
	}

	defer file.Close()

	img, _, err := image.Decode(file)

	if err != nil {
		panic(err)
	}

	return img
}

func loadImages(imageNames []string, dirPath string, ch chan<- []image.Image) {

	fmt.Println("len", len(imageNames))

	images := make([]image.Image, 0, len(imageNames))

	fmt.Println(len(images), images)
	var wg sync.WaitGroup
	var mu sync.Mutex

	for i := range imageNames {
		wg.Add(1)

		index := i

		go func() {
			defer wg.Done()
			filePath := dirPath + imageNames[index]
			image := readFilePathAsImage(filePath)
			mu.Lock()
			images[index] = image
			mu.Unlock()
		}()
	}

	wg.Wait()

	ch <- images
}

// func resizeCollectionImages(images []image.Image) {

// 	resizedImages := make([]image.Image, len(images))

// 	var wg sync.WaitGroup
// 	var mu sync.Mutex

// 	for index := range images {

// 		go func() {
// 			resizedImages
// 		}

// 	}

// }

// func resizeImage(src image.Image, width, height int) *image.RGBA {
// 	dst := image.NewRGBA(image.Rect(0, 0, width, height))
// 	draw.CatmullRom.Scale(dst, dst.Bounds(), src, src.Bounds(), draw.Over, nil)
// 	return dst
// }

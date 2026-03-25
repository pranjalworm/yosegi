package main

import (
	"image"
	"image/draw"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/pranjalworm/stitch/models"
)

func readImageFile(filePath string) (image.Image, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func analyseImage(img image.Image, name string, divisionFactor int) models.AveragedImageData {
	const BitFactor = 257

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	orientation := "portrait"
	if width > height {
		orientation = "landscape"
	}

	data := models.AveragedImageData{
		ImageName:   name,
		Orientation: orientation,
		Width:       width,
		Height:      height,
		Pixel:       make([]models.Pixel, divisionFactor*divisionFactor),
	}

	xBlock := width / divisionFactor
	yBlock := height / divisionFactor
	maxCol := divisionFactor - 1
	maxRow := divisionFactor - 1

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			row := min(y/yBlock, maxRow)
			col := min(x/xBlock, maxCol)
			idx := divisionFactor*row + col

			data.Pixel[idx].Red += uint32(r / BitFactor)
			data.Pixel[idx].Green += uint32(g / BitFactor)
			data.Pixel[idx].Blue += uint32(b / BitFactor)
			data.Pixel[idx].Alpha += uint32(a / BitFactor)
		}
	}

	blockPixelCount := uint32(xBlock * yBlock)
	for i := range data.Pixel {
		data.Pixel[i].Red /= blockPixelCount
		data.Pixel[i].Green /= blockPixelCount
		data.Pixel[i].Blue /= blockPixelCount
		data.Pixel[i].Alpha /= blockPixelCount
	}

	return data
}

func readCollection(dirPath string) ([]models.AveragedImageData, error) {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	var result []models.AveragedImageData
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, file := range files {
		ext := strings.ToLower(filepath.Ext(file.Name()))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
			continue
		}

		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			imgPath := filepath.Join(dirPath, name)
			img, err := readImageFile(imgPath)
			if err != nil {
				return
			}
			analysed := analyseImage(img, name, 1)
			mu.Lock()
			result = append(result, analysed)
			mu.Unlock()
		}(file.Name())
	}

	wg.Wait()
	return result, nil
}

func mapCollectionToTarget(target *models.AveragedImageData, collection []models.AveragedImageData, divisionFactor int) []string {
	result := make([]string, divisionFactor*divisionFactor)

	for i, pixel := range target.Pixel {
		var minDiff uint32 = ^uint32(0)
		var bestName string

		for _, cImg := range collection {
			if target.Orientation != cImg.Orientation {
				continue
			}
			cp := cImg.Pixel[0]
			diff := absDiff(pixel.Red, cp.Red) + absDiff(pixel.Green, cp.Green) +
				absDiff(pixel.Blue, cp.Blue) + absDiff(pixel.Alpha, cp.Alpha)

			if diff < minDiff {
				minDiff = diff
				bestName = cImg.ImageName
			}
		}
		result[i] = bestName
	}
	return result
}

func resizeImage(src image.Image, dstW, dstH int) *image.RGBA {
	bounds := src.Bounds()
	srcW := bounds.Dx()
	srcH := bounds.Dy()
	dst := image.NewRGBA(image.Rect(0, 0, dstW, dstH))
	for y := 0; y < dstH; y++ {
		srcY := bounds.Min.Y + y*srcH/dstH
		for x := 0; x < dstW; x++ {
			srcX := bounds.Min.X + x*srcW/dstW
			dst.Set(x, y, src.At(srcX, srcY))
		}
	}
	return dst
}

func absDiff(a, b uint32) uint32 {
	if a > b {
		return a - b
	}
	return b - a
}

func buildCollage(imageNames []string, dirPath string, target *models.AveragedImageData, divisionFactor int) (*image.RGBA, error) {
	images := make([]image.Image, len(imageNames))
	var wg sync.WaitGroup
	var loadErr error
	var mu sync.Mutex

	for i, name := range imageNames {
		wg.Add(1)
		go func(idx int, fileName string) {
			defer wg.Done()
			imgPath := filepath.Join(dirPath, fileName)
			img, err := readImageFile(imgPath)
			if err != nil {
				mu.Lock()
				loadErr = err
				mu.Unlock()
				return
			}
			mu.Lock()
			images[idx] = img
			mu.Unlock()
		}(i, name)
	}

	wg.Wait()
	if loadErr != nil {
		return nil, loadErr
	}

	collageWidth := target.Width - (target.Width % divisionFactor)
	collageHeight := target.Height - (target.Height % divisionFactor)

	cellW := collageWidth / divisionFactor
	cellH := collageHeight / divisionFactor
	collage := image.NewRGBA(image.Rect(0, 0, collageWidth, collageHeight))

	for i, img := range images {
		col := i % divisionFactor
		row := i / divisionFactor
		resized := resizeImage(img, cellW, cellH)
		offset := image.Pt(col*cellW, row*cellH)
		draw.Draw(collage, resized.Bounds().Add(offset), resized, image.Point{}, draw.Src)
	}

	return collage, nil
}

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
			r, g, b, _ := img.At(x, y).RGBA()
			row := min(y/yBlock, maxRow)
			col := min(x/xBlock, maxCol)
			idx := divisionFactor*row + col

			data.Pixel[idx].Red += uint32(r / BitFactor)
			data.Pixel[idx].Green += uint32(g / BitFactor)
			data.Pixel[idx].Blue += uint32(b / BitFactor)
		}
	}

	blockPixelCount := uint32(xBlock * yBlock)
	for i := range data.Pixel {
		data.Pixel[i].Red /= blockPixelCount
		data.Pixel[i].Green /= blockPixelCount
		data.Pixel[i].Blue /= blockPixelCount
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

func filterByOrientation(collection []models.AveragedImageData, orientation string) []models.AveragedImageData {
	filtered := make([]models.AveragedImageData, 0, len(collection))
	for _, c := range collection {
		if c.Orientation == orientation {
			filtered = append(filtered, c)
		}
	}
	return filtered
}

func mapCollectionToTarget(target *models.AveragedImageData, collection []models.AveragedImageData, divisionFactor int) []string {
	filtered := filterByOrientation(collection, target.Orientation)
	result := make([]string, divisionFactor*divisionFactor)

	var wg sync.WaitGroup
	for i, pixel := range target.Pixel {
		wg.Add(1)
		go func(idx int, px models.Pixel) {
			defer wg.Done()
			var minDiff uint32 = ^uint32(0)
			var bestName string

			for _, cImg := range filtered {
				cp := cImg.Pixel[0]
				d := absDiff(px.Red, cp.Red)
				if d >= minDiff {
					continue
				}
				d += absDiff(px.Green, cp.Green)
				if d >= minDiff {
					continue
				}
				d += absDiff(px.Blue, cp.Blue)
				if d < minDiff {
					minDiff = d
					bestName = cImg.ImageName
				}
			}
			result[idx] = bestName
		}(i, pixel)
	}

	wg.Wait()
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

func buildCollage(imageNames []string, dirPath string, target *models.AveragedImageData, divisionFactor int, tileSize int) (*image.RGBA, error) {
	cellW := (target.Width - (target.Width % divisionFactor)) / divisionFactor
	cellH := (target.Height - (target.Height % divisionFactor)) / divisionFactor

	// Scale cells so the long side matches the requested tileSize.
	if cellW >= cellH {
		if cellW != tileSize {
			cellH = cellH * tileSize / cellW
			cellW = tileSize
		}
	} else {
		if cellH != tileSize {
			cellW = cellW * tileSize / cellH
			cellH = tileSize
		}
	}

	collageWidth := cellW * divisionFactor
	collageHeight := cellH * divisionFactor

	// Collect unique image names.
	uniqueNames := make(map[string]struct{})
	for _, name := range imageNames {
		uniqueNames[name] = struct{}{}
	}

	// Read and resize each unique image once, in parallel.
	resizedCache := make(map[string]*image.RGBA, len(uniqueNames))
	var wg sync.WaitGroup
	var mu sync.Mutex
	var loadErr error

	for name := range uniqueNames {
		wg.Add(1)
		go func(fileName string) {
			defer wg.Done()
			img, err := readImageFile(filepath.Join(dirPath, fileName))
			if err != nil {
				mu.Lock()
				loadErr = err
				mu.Unlock()
				return
			}
			resized := resizeImage(img, cellW, cellH)
			mu.Lock()
			resizedCache[fileName] = resized
			mu.Unlock()
		}(name)
	}

	wg.Wait()
	if loadErr != nil {
		return nil, loadErr
	}

	// Compose the collage from cached resized tiles.
	collage := image.NewRGBA(image.Rect(0, 0, collageWidth, collageHeight))
	for i, name := range imageNames {
		col := i % divisionFactor
		row := i / divisionFactor
		resized := resizedCache[name]
		offset := image.Pt(col*cellW, row*cellH)
		draw.Draw(collage, resized.Bounds().Add(offset), resized, image.Point{}, draw.Src)
	}

	return collage, nil
}

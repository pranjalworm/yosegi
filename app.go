package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx   context.Context
	isDev bool
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.isDev = wailsRuntime.Environment(ctx).BuildType == "dev"
}

type GenerateResult struct {
	OutputPath string `json:"outputPath"`
	TimeTaken  string `json:"timeTaken"`
}

// GetDefaultOutputPath returns ~/Desktop as the default output location.
func (a *App) GetDefaultOutputPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, "Desktop")
}

// SelectTargetImage opens a file dialog for the user to pick a target image.
func (a *App) SelectTargetImage() (string, error) {
	path, err := wailsRuntime.OpenFileDialog(a.ctx, wailsRuntime.OpenDialogOptions{
		Title: "Select Target Image",
		Filters: []wailsRuntime.FileFilter{
			{DisplayName: "Images", Pattern: "*.jpg;*.jpeg;*.png"},
		},
	})
	if err != nil {
		return "", err
	}
	return path, nil
}

// SelectCollectionDirectory opens a directory dialog for the user to pick the collection folder.
func (a *App) SelectCollectionDirectory() (string, error) {
	path, err := wailsRuntime.OpenDirectoryDialog(a.ctx, wailsRuntime.OpenDialogOptions{
		Title: "Select Image Collection Folder",
	})
	if err != nil {
		return "", err
	}
	return path, nil
}

// SelectOutputDirectory opens a directory dialog for the output location.
func (a *App) SelectOutputDirectory() (string, error) {
	path, err := wailsRuntime.OpenDirectoryDialog(a.ctx, wailsRuntime.OpenDialogOptions{
		Title: "Select Output Folder",
	})
	if err != nil {
		return "", err
	}
	return path, nil
}

// OpenFile opens a file using the OS default application.
func (a *App) OpenFile(filePath string) error {
	switch runtime.GOOS {
	case "darwin":
		return exec.Command("open", filePath).Start()
	case "windows":
		return exec.Command("cmd", "/c", "start", "", filePath).Start()
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
}

// CountCollectionImages returns how many valid image files are in a directory.
func (a *App) CountCollectionImages(dirPath string) int {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return 0
	}
	count := 0
	for _, f := range files {
		ext := strings.ToLower(filepath.Ext(f.Name()))
		if ext == ".jpg" || ext == ".jpeg" || ext == ".png" {
			count++
		}
	}
	return count
}

// GetImagePreview returns a base64-encoded thumbnail of an image for the frontend.
func (a *App) GetImagePreview(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return "", err
	}

	bounds := img.Bounds()
	srcW, srcH := bounds.Dx(), bounds.Dy()
	dstW, dstH := srcW, srcH

	maxDim := 800
	if dstW > maxDim || dstH > maxDim {
		if dstW >= dstH {
			dstH = dstH * maxDim / dstW
			dstW = maxDim
		} else {
			dstW = dstW * maxDim / dstH
			dstH = maxDim
		}
	}

	preview := image.NewRGBA(image.Rect(0, 0, dstW, dstH))
	for y := 0; y < dstH; y++ {
		srcY := bounds.Min.Y + y*srcH/dstH
		for x := 0; x < dstW; x++ {
			srcX := bounds.Min.X + x*srcW/dstW
			preview.Set(x, y, img.At(srcX, srcY))
		}
	}

	var buf strings.Builder
	buf.WriteString("data:image/jpeg;base64,")
	encoder := base64.NewEncoder(base64.StdEncoding, &buf)
	if err := jpeg.Encode(encoder, preview, &jpeg.Options{Quality: 90}); err != nil {
		return "", err
	}
	encoder.Close()

	return buf.String(), nil
}

// GenerateCollage runs the full mosaic generation pipeline.
func (a *App) GenerateCollage(targetPath string, collectionPath string, outputDir string, divisionFactor int, tileSize int) (*GenerateResult, error) {
	var logger *log.Logger
	if a.isDev {
		logFile, err := os.OpenFile(filepath.Join(outputDir, "yosegi_perf.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return nil, fmt.Errorf("failed to create log file: %w", err)
		}
		defer logFile.Close()
		logger = log.New(logFile, "", 0)
	} else {
		logger = log.New(io.Discard, "", 0)
	}

	start := time.Now()
	logger.Printf("=== Run started at %s ===", start.Format(time.RFC3339))
	logger.Printf("divisionFactor=%d, collection=%s", divisionFactor, collectionPath)

	wailsRuntime.EventsEmit(a.ctx, "progress", "Reading reference image...")

	stepStart := time.Now()
	targetImg, err := readImageFile(targetPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read target image: %w", err)
	}
	logger.Printf("  Read target image:      %s", time.Since(stepStart))

	stepStart = time.Now()
	targetData := analyseImage(targetImg, "target", divisionFactor)
	logger.Printf("  Analyse target image:   %s", time.Since(stepStart))

	wailsRuntime.EventsEmit(a.ctx, "progress", "Analyzing collection images...")

	stepStart = time.Now()
	collectionData, err := readCollection(collectionPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read collection: %w", err)
	}
	logger.Printf("  Read collection (%d images): %s", len(collectionData), time.Since(stepStart))

	if len(collectionData) == 0 {
		return nil, fmt.Errorf("no valid images found in collection directory")
	}

	wailsRuntime.EventsEmit(a.ctx, "progress", "Matching tiles to reference...")

	stepStart = time.Now()
	resultNames := mapCollectionToTarget(&targetData, collectionData, divisionFactor)
	logger.Printf("  Map collection to target: %s", time.Since(stepStart))

	wailsRuntime.EventsEmit(a.ctx, "progress", "Building collage...")

	stepStart = time.Now()
	collage, err := buildCollage(resultNames, collectionPath, &targetData, divisionFactor, tileSize)
	if err != nil {
		return nil, fmt.Errorf("failed to build collage: %w", err)
	}
	logger.Printf("  Build collage:          %s", time.Since(stepStart))

	timestamp := time.Now().Format("2006-01-02_15-04-05")
	outputFileName := fmt.Sprintf("yosegi_%s.jpg", timestamp)
	outputPath := filepath.Join(outputDir, outputFileName)

	outFile, err := os.Create(outputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	wailsRuntime.EventsEmit(a.ctx, "progress", "Encoding output...")

	stepStart = time.Now()
	if err := jpeg.Encode(outFile, collage, &jpeg.Options{Quality: 90}); err != nil {
		return nil, fmt.Errorf("failed to encode JPEG: %w", err)
	}
	logger.Printf("  Encode JPEG:            %s", time.Since(stepStart))

	elapsed := time.Since(start)
	logger.Printf("  TOTAL:                  %s", elapsed)
	logger.Println()

	wailsRuntime.EventsEmit(a.ctx, "progress", "Preparing for the grand reveal...")

	return &GenerateResult{
		OutputPath: outputPath,
		TimeTaken:  elapsed.Round(time.Millisecond).String(),
	}, nil
}

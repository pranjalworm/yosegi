# Yosegi

A desktop application that creates photo mosaics — collages that resemble a target image, built from your own image collection. Powered by Go + Wails with a Vue.js frontend.

## How It Works

1. Select a target image and a folder of collection images
2. The app divides the target into an N×N grid and calculates the average color of each cell
3. Each grid cell is matched to the collection image with the closest average color
4. The matched images are stitched together into a final mosaic

## Getting Started

### Prerequisites

- Go 1.25+
- Node.js 20+
- [Wails CLI](https://wails.io/docs/gettingstarted/installation) v2

### Development

```bash
wails dev
```

### Build

```bash
wails build
```

The built app is output to `build/bin/Yosegi.app` on macOS.

## Project Structure

```
main.go           # Wails app entry point
app.go            # Backend logic (image processing, file dialogs, collage generation)
models/           # Shared data structures (Pixel, AveragedImageData)
frontend/         # Vue.js + TypeScript + Vite frontend
  src/
    App.vue       # Main application view
    components/   # UI components (TitleBar, FileInput, ImagePreview)
    style.css     # Global styles
build/            # Build configuration and app icons
wails.json        # Wails project configuration
```

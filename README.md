# Stitch

A Go application that creates a photo collage resembling a given target image by stitching together images from a collection based on average color matching.

## How It Works

1. Divides the target image into an N×N grid and calculates the average color of each cell
2. Analyzes all images in a collection directory for their average color
3. Maps each grid cell to the closest-matching collection image
4. Stitches the matched images together into a final collage

## Usage

```bash
go run main.go <imageName> <divisionFactor> <collectionName>
```

- `imageName` — name of the target image (without `.jpg` extension), placed in `./images/target/`
- `divisionFactor` — grid size (e.g., `10` creates a 10×10 grid)
- `collectionName` — directory name under `./images/collections/` containing the source images

Output is saved as a timestamped PNG in `./images/output/`.

## Project Structure

```
imager/       # Core image processing (reading, matching, collage creation)
models/       # Data structures (Pixel, AveragedImageData)
utils/        # Performance tracking helpers
images/
  target/       # Place target images here
  collections/  # Place image collections here
  output/       # Generated collages
```

## Requirements

- Go 1.25+
- No external dependencies — uses only the standard library

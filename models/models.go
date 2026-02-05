package models

type Pixel struct {
	Red   uint32
	Green uint32
	Blue  uint32
	Alpha uint32
}

type GridBlockPixel struct {
	Pixel Pixel
}

type AveragedImageData struct {
	Width          int
	Height         int
	Format         string
	GridBlockPixel []GridBlockPixel
}

package models

type Pixel struct {
	Red   uint32
	Green uint32
	Blue  uint32
	Alpha uint32
}

type AveragedImageData struct {
	ImageName   string
	Orientation string // "portrait" | "landscape"
	Width       int
	Height      int
	Pixel       []Pixel
}

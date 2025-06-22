package bmp

type BMP struct {
	Header Header
	Image  Image
}

type Header struct {
	FileType       string
	FileSize       uint32
	HeaderSize     uint32
	DibHeaderSize  uint32
	WidthInPixels  int32
	HeightInPixels int32
	PixelSize      uint16
	ImageSize      uint32
}

type Image struct {
	Width  int
	Height int
	Pixels []Pixel
}

type Pixel struct {
	B, G, R uint8
}

type DIBHeader struct {
	HeaderSize      uint32
	Width           int32
	Height          int32
	Planes          uint16
	BitsPerPixel    uint16
	Compression     uint32
	ImageSize       uint32
	XPixelsPerMeter int32
	YPixelsPerMeter int32
	ColorsUsed      uint32
	ImportantColors uint32
}

const (
	bitmapFileHeaderSize = 14
	bitmapInfoHeaderSize = 40
)

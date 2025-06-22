package bmp

import (
	"io"
	"log"
	"os"
)

func ReadBMP(path string) *BMP {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	header := ReadHeader(path)
	image := ReadImage(file, *header)

	bmp := &BMP{
		Header: *header,
		Image:  image,
	}

	return bmp
}

func ReadImage(file *os.File, header Header) Image {
	width := int(header.WidthInPixels)
	height := int(header.HeightInPixels)

	_, err := file.Seek(int64(header.HeaderSize+header.DibHeaderSize), io.SeekCurrent)
	if err != nil {
		log.Fatalf("Failed to seek image data %v", err)
	}

	rowSize := ((int(header.PixelSize)*width + 31) / 32) * 4
	pixelData := make([]byte, rowSize)
	image := Image{
		Width:  width,
		Height: height,
		Pixels: make([]Pixel, width*height),
	}

	for y := height - 1; y >= 0; y-- {
		_, err := file.Read(pixelData)
		if err != nil {
			log.Fatalf("Failed to read pixel data: %v", err)
		}
		for x := 0; x < width; x++ {
			idx := x * 3
			pixelIdx := y*width + x
			image.Pixels[pixelIdx] = Pixel{
				B: pixelData[idx],
				G: pixelData[idx+1],
				R: pixelData[idx+2],
			}
		}
	}

	return image
}

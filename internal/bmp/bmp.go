package bmp

import (
	"encoding/binary"
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

func WriteBMP(path string, bmp *BMP) {
	file, err := os.Create(path)
	if err != nil {
		log.Fatalf("Failed to create file: %v", err)
	}
	defer file.Close()

	prepareHeader(&bmp.Header)

	writeHeader(file, &bmp.Header)
	writeImage(file, &bmp.Image, bmp.Header)
}

func prepareHeader(header *Header) {
	header.FileType = "BMP"
	header.HeaderSize = bitmapFileHeaderSize
	header.DibHeaderSize = bitmapInfoHeaderSize
	header.PixelSize = 24

	width := int(header.WidthInPixels)
	height := int(header.HeightInPixels)
	rowSize := ((width*3 + 3) / 4) * 4
	header.ImageSize = uint32(rowSize * height)
	header.FileSize = header.HeaderSize + header.DibHeaderSize + header.ImageSize
}

func writeHeader(file *os.File, header *Header) {
	writeLE(file, []byte(header.FileType), "file type")
	writeLE(file, header.FileSize, "file size")
	writeLE(file, uint32(0), "reserved")
	writeLE(file, header.HeaderSize+header.DibHeaderSize, "offset")
	writeLE(file, header.DibHeaderSize, "DIB header size")
	writeLE(file, header.WidthInPixels, "width")
	writeLE(file, header.HeightInPixels, "height")
	writeLE(file, uint16(1), "planes")
	writeLE(file, header.PixelSize, "pixel size")
	writeLE(file, uint32(0), "compression")
	writeLE(file, header.ImageSize, "image size")
	writeLE(file, uint32(0), "X pixels per meter")
	writeLE(file, uint32(0), "Y pixels per meter")
	writeLE(file, uint32(0), "colors used")
	writeLE(file, uint32(0), "important colors")
}

func writeImage(file *os.File, image *Image, header Header) {
	width := int(header.WidthInPixels)
	height := int(header.HeightInPixels)
	rowSize := ((width*3 + 3) / 4) * 4
	pixelData := make([]byte, rowSize)

	for y := height - 1; y >= 0; y-- {
		for x := 0; x < width; x++ {
			idx := x * 3
			pixelIdx := y*width + x
			pixelData[idx] = image.Pixels[pixelIdx].B
			pixelData[idx+1] = image.Pixels[pixelIdx].G
			pixelData[idx+2] = image.Pixels[pixelIdx].R
		}
		_, err := file.Write(pixelData)
		if err != nil {
			log.Fatalf("Failed to write pixel data: %v", err)
		}
	}

	writtenSize := uint32(rowSize * height)
	if header.ImageSize > writtenSize {
		padding := make([]byte, header.ImageSize-writtenSize)
		_, err := file.Write(padding)
		if err != nil {
			log.Fatalf("Failed to write padding: %v", err)
		}
	}
}

func writeLE(w io.Writer, data interface{}, field string) {
	err := binary.Write(w, binary.LittleEndian, data)
	if err != nil {
		log.Fatalf("Failed to write %s: %v", field, err)
	}
}

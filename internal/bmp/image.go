package bmp

import (
	"encoding/binary"
	"log"
	"os"
)

func SaveImage(img *Image, path string) {
	file, err := os.Create(path)
	if err != nil {
		log.Fatalf("Failed to create file: %v", err)
	}
	defer file.Close()

	header := buildHeader(img)

	binary.Write(file, binary.LittleEndian, []byte(header.FileType))
	binary.Write(file, binary.LittleEndian, header.FileSize)
	binary.Write(file, binary.LittleEndian, uint32(0))                                         // Зарезервировано
	binary.Write(file, binary.LittleEndian, uint32(bitmapFileHeaderSize+bitmapInfoHeaderSize)) // Offset
	binary.Write(file, binary.LittleEndian, header.DibHeaderSize)
	binary.Write(file, binary.LittleEndian, header.WidthInPixels)
	binary.Write(file, binary.LittleEndian, header.HeightInPixels)
	binary.Write(file, binary.LittleEndian, uint16(1)) // Planes
	binary.Write(file, binary.LittleEndian, header.PixelSize)
	binary.Write(file, binary.LittleEndian, uint32(0)) // Compression
	binary.Write(file, binary.LittleEndian, header.ImageSize)
	binary.Write(file, binary.LittleEndian, uint32(0)) // X Pixels per meter
	binary.Write(file, binary.LittleEndian, uint32(0)) // Y Pixels per meter
	binary.Write(file, binary.LittleEndian, uint32(0)) // Colors used
	binary.Write(file, binary.LittleEndian, uint32(0)) // Important colors

	width := img.Width
	height := img.Height

	for y := height - 1; y >= 0; y-- {
		start := y * width
		end := start + width
		for _, pixel := range img.Pixels[start:end] {
			binary.Write(file, binary.LittleEndian, pixel)
		}
	}
}

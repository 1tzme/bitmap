package transform

import (
	"bitmap/internal/bmp"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	u "bitmap/internal/utils"
)

func HandleRotateCommand() {
	var rotates []string
	rotateFlag := flag.NewFlagSet("apply", flag.ExitOnError)
	rotateFlag.Func("rotate", "Rotate directIion: right, 90, 180, 270, left, -90, -180, -270", func(s string) error {
		s = strings.ToLower(s)
		if !isValidRotation(s) {
			return fmt.Errorf("invalid rotate option: %s", s)
		}
		rotates = append(rotates, s)
		return nil
	})

	if err := rotateFlag.Parse(os.Args[2:]); err != nil {
		log.Fatalf("Error parsing flags: %v", err)
	}

	if len(rotateFlag.Args()) != 2 {
		u.PrintApplyUsage()
		log.Fatal("Error: expected source and output file")
	}

	sourceFile, outputFile := rotateFlag.Args()[0], rotateFlag.Args()[1]
	header := bmp.ReadHeader(sourceFile)

	file, err := os.Open(sourceFile)
	if err != nil {
		log.Fatalln("Failer to open file: ", err)
	}
	defer file.Close()

	img := bmp.ReadImage(file, *header)
	rotatedImg := applyRotations(img, rotates)
	bmp.SaveImage(rotatedImg, outputFile)
}

func isValidRotation(option string) bool {
	validOptions := map[string]bool{
		"right": true, "90": true, "180": true, "270": true,
		"left": true, "-90": true, "-180": true, "-270": true,
	}
	return validOptions[strings.ToLower(option)]
}

func applyRotations(img bmp.Image, rotations []string) *bmp.Image {
	result := &img
	for _, rot := range rotations {
		angle := normalizeRotation(rot)
		result = rotateImage(result, angle)
	}
	return result
}

func normalizeRotation(rot string) int {
	rot = strings.ToLower(rot)
	switch rot {
	case "right", "90":
		return 90
	case "left", "-90":
		return -90
	case "180", "-180":
		return 180
	case "270":
		return 270
	case "-270":
		return -270
	default:
		log.Fatalf("Invalid rotation: %s", rot)
		return 0
	}
}

func rotateImage(img *bmp.Image, angle int) *bmp.Image {
	angle = (angle%360 + 360) % 360 // Нормализация угла
	width, height := img.Width, img.Height
	var newWidth, newHeight int
	var newPixels []bmp.Pixel

	switch angle {
	case 90, 270:
		newWidth, newHeight = height, width
	default:
		newWidth, newHeight = width, height
	}

	newPixels = make([]bmp.Pixel, newWidth*newHeight)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			idx := y*width + x
			var newX, newY int
			switch angle {
			case 90:
				newX = height - 1 - y
				newY = x
			case 180:
				newX = width - 1 - x
				newY = height - 1 - y
			case 270:
				newX = y
				newY = width - 1 - x
			default: // 0 градусов
				newX = x
				newY = y
			}
			newIdx := newY*newWidth + newX
			newPixels[newIdx] = img.Pixels[idx]
		}
	}

	return &bmp.Image{
		Width:  newWidth,
		Height: newHeight,
		Pixels: newPixels,
	}
}

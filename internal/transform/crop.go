package transform

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	b "bitmap/internal/bmp"
	u "bitmap/internal/utils"
)

type cropParams struct {
	OffsetX, OffsetY, Width, Height int
}

func HandleCropCommand() {
	var crops []cropParams
	cropFlag := flag.NewFlagSet("apply", flag.ExitOnError)
	cropFlag.Func("crop", "Crop parameters: OffsetX-OffsetY[-Width-Height]", func(s string) error {
		params, err := parseCropParams(s)
		if err != nil {
			return fmt.Errorf("invalid crop parameters: %w", err)
		}
		crops = append(crops, params)
		return nil
	})

	if err := cropFlag.Parse(os.Args[2:]); err != nil {
		log.Fatalf("Error parsing flags: %v", err)
	}

	if len(cropFlag.Args()) != 2 {
		u.PrintApplyUsage()
		log.Fatal("Error: expected source and output file")
	}

	sourceFile, outputFile := cropFlag.Args()[0], cropFlag.Args()[1]

	header := b.ReadHeader(sourceFile)

	file, err := os.Open(sourceFile)
	if err != nil {
		log.Fatalln("Failer to open file: ", err)
	}
	defer file.Close()

	img := b.ReadImage(file, *header)
	croppedImg, err := applyCrops(&img, crops)
	if err != nil {
		log.Fatalf("Error applying crop: %v", err)
	}
	b.SaveImage(croppedImg, outputFile)
}

func parseCropParams(s string) (cropParams, error) {
	parts := strings.Split(s, "-")
	if len(parts) != 2 && len(parts) != 4 {
		return cropParams{}, fmt.Errorf("expected 2 or 4 values, got %d", len(parts))
	}

	offsetX, err := strconv.Atoi(parts[0])
	if err != nil || offsetX < 0 {
		return cropParams{}, fmt.Errorf("invalid OffsetX: %s", parts[0])
	}
	offsetY, err := strconv.Atoi(parts[1])
	if err != nil || offsetY < 0 {
		return cropParams{}, fmt.Errorf("invalid OffsetY: %s", parts[1])
	}

	var width, height int
	if len(parts) == 4 {
		width, err = strconv.Atoi(parts[2])
		if err != nil || width <= 0 {
			return cropParams{}, fmt.Errorf("invalid Width: %s", parts[2])
		}
		height, err = strconv.Atoi(parts[3])
		if err != nil || height <= 0 {
			return cropParams{}, fmt.Errorf("invalid Height: %s", parts[3])
		}
	} else {
		width = -1
		height = -1
	}

	return cropParams{OffsetX: offsetX, OffsetY: offsetY, Width: width, Height: height}, nil
}

func applyCrops(img *b.Image, crops []cropParams) (*b.Image, error) {
	result := img
	for _, crop := range crops {
		var err error
		result, err = cropImage(result, crop)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func cropImage(img *b.Image, params cropParams) (*b.Image, error) {
	width, height := img.Width, img.Height
	cropWidth := params.Width
	cropHeight := params.Height

	if cropWidth == -1 {
		cropWidth = width - params.OffsetX
	}
	if cropHeight == -1 {
		cropHeight = height - params.OffsetY
	}

	if params.OffsetX < 0 || params.OffsetX >= width || params.OffsetY < 0 || params.OffsetY >= height {
		return nil, fmt.Errorf("crop offset (%d,%d) out of image bounds (%d,%d)", params.OffsetX, params.OffsetY, width, height)
	}
	if cropWidth <= 0 || cropHeight <= 0 || params.OffsetX+cropWidth > width || params.OffsetY+cropHeight > height {
		return nil, fmt.Errorf("crop dimensions (%d,%d) at offset (%d,%d) exceed image bounds (%d,%d)", cropWidth, cropHeight, params.OffsetX, params.OffsetY, width, height)
	}

	newPixels := make([]b.Pixel, cropWidth*cropHeight)
	for y := 0; y < cropHeight; y++ {
		for x := 0; x < cropWidth; x++ {
			srcIdx := (y+params.OffsetY)*width + (x + params.OffsetX)
			dstIdx := y*cropWidth + x
			newPixels[dstIdx] = img.Pixels[srcIdx]
		}
	}

	return &b.Image{
		Width:  cropWidth,
		Height: cropHeight,
		Pixels: newPixels,
	}, nil
}

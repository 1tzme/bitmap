package transform

import (
	"fmt"
	"math"
	"os"

	bm "bitmap/internal/bmp"
)

func applyFilters(image *bm.Image, filters []string) {
	for _, filter := range filters {
		switch filter {
		case "blue":
			applyBlueFilter(image)
		case "red":
			applyRedFilter(image)
		case "green":
			applyGreenFilter(image)
		case "grayscale":
			applyGrayscaleFilter(image)
		case "negative":
			applyNegativeFilter(image)
		case "pixelate":
			applyPixelateFilter(image, 20)
		case "blur":
			applyBlurFilter(image, 20)
		}
	}
}

func applyBlueFilter(image *bm.Image) {
	for i := range image.Pixels {
		image.Pixels[i].R = 0
		image.Pixels[i].G = 0
	}
}

func applyRedFilter(image *bm.Image) {
	for i := range image.Pixels {
		image.Pixels[i].B = 0
		image.Pixels[i].G = 0
	}
}

func applyGreenFilter(image *bm.Image) {
	for i := range image.Pixels {
		image.Pixels[i].B = 0
		image.Pixels[i].R = 0
	}
}

func applyGrayscaleFilter(image *bm.Image) {
	for i := range image.Pixels {
		gray := uint8(0.299*float64(image.Pixels[i].R) + 0.587*float64(image.Pixels[i].G) + 0.114*float64(image.Pixels[i].B))
		image.Pixels[i].R = gray
		image.Pixels[i].G = gray
		image.Pixels[i].B = gray
	}
}

func applyNegativeFilter(image *bm.Image) {
	for i := range image.Pixels {
		image.Pixels[i].R = 255 - image.Pixels[i].R
		image.Pixels[i].G = 255 - image.Pixels[i].G
		image.Pixels[i].B = 255 - image.Pixels[i].B
	}
}

func applyPixelateFilter(image *bm.Image, blockSize int) {
	if blockSize < 0 {
		fmt.Fprintf(os.Stderr, "Block size cannot be negative: %v\n", blockSize)
		os.Exit(1)
	} else if blockSize == 0 {
		return
	}
	for y := 0; y < image.Height; y += blockSize {
		for x := 0; x < image.Width; x += blockSize {
			r, g, b, count := calculateBlockAverage(image, x, y, blockSize)
			if count > 0 {
				setBlockColor(image, x, y, blockSize, uint8(r/float64(count)), uint8(g/float64(count)), uint8(b/float64(count)))
			}
		}
	}
}

func calculateBlockAverage(image *bm.Image, x, y, blockSize int) (r, g, b float64, count int) {
	for dy := 0; dy < blockSize && y+dy < image.Height; dy++ {
		for dx := 0; dx < blockSize && x+dx < image.Width; dx++ {
			pixel := image.Pixels[(y+dy)*image.Width+(x+dx)]
			r += float64(pixel.R)
			g += float64(pixel.G)
			b += float64(pixel.B)
			count++
		}
	}
	return
}

func setBlockColor(image *bm.Image, x, y, blockSize int, avgR, avgG, avgB uint8) {
	for dy := 0; dy < blockSize && y+dy < image.Height; dy++ {
		for dx := 0; dx < blockSize && x+dx < image.Width; dx++ {
			image.Pixels[(y+dy)*image.Width+(x+dx)] = bm.Pixel{R: avgR, G: avgG, B: avgB}
		}
	}
}

func applyBlurFilter(image *bm.Image, kernelSize int) {
	if kernelSize < 0 {
		fmt.Fprintf(os.Stderr, "Kernel Size cannot be negative: %v\n", kernelSize)
		os.Exit(1)
	}
	if kernelSize%2 == 0 {
		kernelSize++
	}
	halfKernel := kernelSize / 2
	original := make([]bm.Pixel, len(image.Pixels))
	copy(original, image.Pixels)

	for y := 0; y < image.Height; y++ {
		for x := 0; x < image.Width; x++ {
			r, g, b, count := calculateBlurAverage(image, x, y, halfKernel, original)
			if count > 0 {
				image.Pixels[y*image.Width+x] = bm.Pixel{
					R: uint8(math.Round(r / float64(count))),
					G: uint8(math.Round(g / float64(count))),
					B: uint8(math.Round(b / float64(count))),
				}
			}
		}
	}
}

func calculateBlurAverage(image *bm.Image, x, y, halfKernel int, original []bm.Pixel) (r, g, b float64, count int) {
	for dy := -halfKernel; dy <= halfKernel; dy++ {
		for dx := -halfKernel; dx <= halfKernel; dx++ {
			ny, nx := y+dy, x+dx
			if ny >= 0 && ny < image.Height && nx >= 0 && nx < image.Width {
				pixel := original[ny*image.Width+nx]
				r += float64(pixel.R)
				g += float64(pixel.G)
				b += float64(pixel.B)
				count++
			}
		}
	}
	return
}

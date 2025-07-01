package transform

import (
	"errors"
	"strings"

	b "bitmap/internal/bmp"
)

func ApplyMirrors(img *b.Image, mirrorFlags []string) error {
	for _, m := range mirrorFlags {
		switch normalizeMirrorFlag(m) {
		case "horizontal":
			mirrorHorizontal(img)
		case "vertical":
			mirrorVertical(img)
		default:
			return errors.New("invalid mirror direction: " + m)
		}
	}
	return nil
}

func normalizeMirrorFlag(flag string) string {
	flag = strings.ToLower(flag)
	switch flag {
	case "h", "hor", "horizontal", "horizontally":
		return "horizontal"
	case "v", "ver", "vertical", "vertically":
		return "vertical"
	default:
		return ""
	}
}

func mirrorHorizontal(img *b.Image) {
	width := img.Width
	height := img.Height
	for y := 0; y < height; y++ {
		for x := 0; x < width/2; x++ {
			left := y*width + x
			right := y*width + (width - 1 - x)
			img.Pixels[left], img.Pixels[right] = img.Pixels[right], img.Pixels[left]
		}
	}
}

func mirrorVertical(img *b.Image) {
	width := img.Width
	height := img.Height
	for y := 0; y < height/2; y++ {
		for x := 0; x < width; x++ {
			top := y*width + x
			bottom := (height-1-y)*width + x
			img.Pixels[top], img.Pixels[bottom] = img.Pixels[bottom], img.Pixels[top]
		}
	}
}

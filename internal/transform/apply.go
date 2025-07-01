package transform

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	b "bitmap/internal/bmp"
	u "bitmap/internal/utils"
)

type flagSlice []string

func (f *flagSlice) String() string { return "" }
func (f *flagSlice) Set(value string) error {
	*f = append(*f, value)
	return nil
}

func HandleApplyCommand() error {
	fs := flag.NewFlagSet("apply", flag.ContinueOnError)

	var filters flagSlice
	var mirrors flagSlice
	var rotates flagSlice
	var crops flagSlice

	fs.Var(&filters, "filter", "Filter to apply (blue, red, green, grayscale, negative, pixelate, blur)")
	fs.Var(&mirrors, "mirror", "Mirror direction: horizontal, vertical")
	fs.Var(&rotates, "rotate", "Rotate direction: right, 90, 180, 270, left, -90, -180, -270")
	fs.Var(&crops, "crop", "Crop area: OffsetX-OffsetY[-Width-Height]")

	args := normalizeFlags(os.Args[2:])
	if err := fs.Parse(args); err != nil {
		u.PrintApplyUsage()
		return fmt.Errorf("error parsing flags: %v", err)
	}

	if len(fs.Args()) != 2 {
		u.PrintApplyUsage()
		return fmt.Errorf("expected source and output files")
	}

	if len(filters) == 0 && len(rotates) == 0 && len(mirrors) == 0 && len(crops) == 0 {
		u.PrintApplyUsage()
		return fmt.Errorf("no transformations specified")
	}

	sourceFile, outputFile := fs.Args()[0], fs.Args()[1]

	if _, err := os.Stat(sourceFile); os.IsNotExist(err) {
		return fmt.Errorf("source file %s does not exist", sourceFile)
	}
	if dir := filepath.Dir(outputFile); dir != "." && dir != "" {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			return fmt.Errorf("output directory %s does not exist", dir)
		}
	}

	bmp := b.ReadBMP(sourceFile)

	if len(crops) > 0 {
		var cropList []cropParams
		for _, c := range crops {
			param, err := parseCropParams(c)
			if err != nil {
				return fmt.Errorf("invalid crop parameters: %v", err)
			}
			cropList = append(cropList, param)
		}
		cropped, err := applyCrops(&bmp.Image, cropList)
		if err != nil {
			return fmt.Errorf("failed to apply crops: %v", err)
		}
		bmp.Image = *cropped
	}

	if len(mirrors) > 0 {
		if err := ApplyMirrors(&bmp.Image, mirrors); err != nil {
			return fmt.Errorf("failed to apply mirrors: %v", err)
		}
	}

	if len(rotates) > 0 {
		bmp.Image = *applyRotations(bmp.Image, rotates)
	}

	bmp.Header.WidthInPixels = int32(bmp.Image.Width)
	bmp.Header.HeightInPixels = int32(bmp.Image.Height)

	if len(filters) > 0 {
		applyFilters(&bmp.Image, filters)
	}

	b.WriteBMP(outputFile, bmp)
	return nil
}

func normalizeFlags(args []string) []string {
	var normalized []string
	for _, arg := range args {
		if strings.HasPrefix(arg, "---") {
			normalized = append(normalized, strings.Replace(arg, "---", "--", 1))
		} else {
			normalized = append(normalized, arg)
		}
	}
	return normalized
}

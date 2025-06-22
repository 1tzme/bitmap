package utils

import "fmt"

func PrintUsage() {
	fmt.Println(`$ ./bitmap
Usage:
  bitmap <command> [arguments]

The commands are:
  header    prints bitmap file header information
  apply     applies processing to the image and saves it to the file`)
}

func PrintHeaderUsage() {
	fmt.Println(`$ ./bitmap header --helps
Usage:
  bitmap header <source_file>

Description:
  Prints bitmap file header information`)
}

func PrintApplyUsage() {
	fmt.Println(`$ ./bitmap apply --help
Usage:
bitmap apply [options] <source_file> <output_file>

The options are:
-h, --help              prints program usage information
--filter=<filter_type>  applies a filter to the image (blue, red, green, grayscale, negative, pixelate, blur)`)
}

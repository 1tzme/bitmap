# Bitmap Image Processing Tool

A command-line tool for reading, writing, and manipulating BMP (bitmap) image files written in Go.

## Features

- **Header Analysis**: Read and display BMP file header information
- **Image Transformations**: Mirror (horizontal/vertical) and rotate images
- **Filters**: Apply various color filters (RGB channels, grayscale, negative, blur, pixelate)
- **Cropping**: Crop images with flexible parameters
- **Combine Operations**: Chain multiple operations together

## Installation

```bash
go build -o bitmap .
```

## Usage

### Display Header Information

```bash
./bitmap header <input_file.bmp>
```

Example:
```bash
./bitmap header sample.bmp
```

Output:
```
BMP Header:
- FileType BM
- FileSizeInBytes 518456
- HeaderSize 54
DIB Header:
- DibHeaderSize 40
- WidthInPixels 480
- HeightInPixels 360
- PixelSizeInBits 24
- ImageSizeInBytes 518402
```

### Apply Transformations

```bash
./bitmap apply [options] <input_file.bmp> <output_file.bmp>
```

#### Mirror Operations

```bash
# Mirror horizontally
./bitmap apply --mirror=horizontal input.bmp output.bmp

# Mirror vertically  
./bitmap apply --mirror=vertical input.bmp output.bmp

```

#### Rotation

```bash
# Rotate right (90° clockwise)
./bitmap apply --rotate=right input.bmp output.bmp

# Rotate 180°
./bitmap apply --rotate=180 input.bmp output.bmp

# Supported rotate options: right, left, 90, 180, 270, -90, -180, -270
```

#### Color Filters

```bash
# Keep only red channel
./bitmap apply --filter=red input.bmp output.bmp

# Keep only green channel
./bitmap apply --filter=green input.bmp output.bmp

# Keep only blue channel
./bitmap apply --filter=blue input.bmp output.bmp

# Convert to grayscale
./bitmap apply --filter=grayscale input.bmp output.bmp

# Apply negative filter
./bitmap apply --filter=negative input.bmp output.bmp

# Apply pixelation (20x20 blocks)
./bitmap apply --filter=pixelate input.bmp output.bmp

# Apply blur effect
./bitmap apply --filter=blur input.bmp output.bmp
```

#### Cropping

```bash
# Crop with offset and size: OffsetX-OffsetY-Width-Height
./bitmap apply --crop=20-20-100-100 input.bmp output.bmp

# Crop with offset only (uses remaining image size)
./bitmap apply --crop=400-300 input.bmp output.bmp
```

#### Combine Multiple Operations

Operations are applied sequentially in the order specified:

```bash
./bitmap apply --mirror=horizontal --rotate=right --filter=negative --crop=50-50-200-200 input.bmp output.bmp
```

### Help

```bash
./bitmap --help
```
### Contributors

- [armakhat](https://platform.alem.school/git/armakhat)  
- [zazholdas](https://platform.alem.school/git/zazholdas)  
- [zaaripzha](https://platform.alem.school/git/zaaripzha)  
- [datussupo](https://platform.alem.school/git/datussupo)  
- [nuarystan](https://platform.alem.school/git/nuarystan)  

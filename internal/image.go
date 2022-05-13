package internal

import (
	"image"
	"image/jpeg"
	"os"
)

const colorMax = 65535.0

type Options struct {
	Path    string
	Width   int
	Height  int
	Palette string
}

func loadImage(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return jpeg.Decode(f)
}

func (opt *Options) String() string {
	img, err := loadImage(opt.Path)
	if err != nil {
		panic(err)
	}

	return process(img, opt)
}

func process(img image.Image, options *Options) string {
	result := make([]byte, options.Width*options.Height+options.Height)

	chunkWidth, chunkHeight := chunkSizes(img, options)

	resIdx := 0
	for y := 0; y < options.Height; y++ {
		for x := 0; x < options.Width; x++ {
			intensity := chunkIntensity(img, x*chunkWidth, y*chunkHeight, chunkWidth, chunkHeight)
			charIdx := float32(intensity) / colorMax * float32(len(options.Palette))
			result[resIdx] = options.Palette[int(charIdx)]
			resIdx++
		}

		result[resIdx] = byte('\n')
		resIdx++
	}

	return string(result)
}

func chunkIntensity(img image.Image, x, y, chunkWidth, chunkHeight int) int {
	total := 0
	for yOffset := 0; yOffset < chunkHeight; yOffset++ {
		for xOffset := 0; xOffset < chunkWidth; xOffset++ {
			r, g, b, _ := img.At(x+xOffset, y+yOffset).RGBA()
			total += int(r) + int(g) + int(b)
		}
	}
	return total / 3 / (chunkWidth * chunkHeight)
}

func chunkSizes(img image.Image, options *Options) (int, int) {
	imgWidth := img.Bounds().Max.X
	imgHeight := img.Bounds().Max.Y
	newImgWidth := imgWidth - (imgWidth % options.Width)
	chunkSizeX := newImgWidth / options.Width
	newImgHeight := imgHeight - (imgHeight % options.Height)
	chunkSizeY := newImgHeight / options.Height

	return chunkSizeX, chunkSizeY
}

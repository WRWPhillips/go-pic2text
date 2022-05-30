package internal

import (
	"image"
	_ "image/jpeg"
	_ "image/png"

	"os"
)

const colorMax = 65535.0

type Options struct {
	Path    string
	Width   int
	Height  int
	Palette string
	Reverse bool
}

func loadImage(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}
	return img, err
}

func (opt *Options) String() string {
	img, err := loadImage(opt.Path)
	if err != nil {
		panic(err)
	}

	return process(img, opt)
}

func reverse(input string) string {
	n := 0
	rune := make([]rune, len(input))
	for _, r := range input {
		rune[n] = r
		n++
	}
	rune = rune[0:n]
	for i := 0; i < n/2; i++ {
		rune[i], rune[n-1-i] = rune[n-1-i], rune[i]
	}
	output := string(rune)
	return output
}

func process(img image.Image, options *Options) string {
	result := make([]byte, options.Width*options.Height+options.Height)
	if options.Reverse == true {
		options.Palette = reverse(options.Palette)
	}

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

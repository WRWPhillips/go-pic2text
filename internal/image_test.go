package internal

import (
	"image"
	"image/color"
	"testing"
)

type fakeImage struct {
	width  int
	height int
}

func (fakeImage) ColorModel() color.Model {
	return color.RGBAModel
}

func (img fakeImage) Bounds() image.Rectangle {
	return image.Rectangle{
		Min: image.Point{
			X: 0,
			Y: 0,
		},
		Max: image.Point{
			X: img.width,
			Y: img.height,
		},
	}
}

func (fakeImage) At(x, y int) color.Color {
	return color.RGBA{
		R: 0,
		G: 0,
		B: 0,
		A: 0,
	}
}

func TestChunkSize(t *testing.T) {
	options := &Options{
		Width:  80,
		Height: 25,
	}

	t.Run("100x100", func(t *testing.T) {
		img := fakeImage{
			width:  100,
			height: 100,
		}

		width, height := chunkSizes(img, options)
		if width != 1 {
			t.Fatalf("Expected %d, got %d", 1, width)
		}
		if height != 4 {
			t.Fatalf("Expected %d, got %d", 4, width)
		}
	})

	t.Run("1920x1080", func(t *testing.T) {
		img := fakeImage{
			width:  1920,
			height: 1080,
		}

		width, height := chunkSizes(img, options)
		if width != 24 {
			t.Fatalf("Expected %d, got %d", 24, width)
		}
		if height != 43 {
			t.Fatalf("Expected %d, got %d", 43, height)
		}
	})
}

func BenchmarkProcess(b *testing.B) {
	b.StopTimer()

	options := &Options{
		Path:    "/home/pi/src/go-pic2text/assets/input.jpeg",
		Width:   120,
		Height:  60,
		Palette: "'$@B%8&WM#*oahkbdpqwmZO0QLCJUYXzcvunxrjft/\\|()1{}[]?-_+~<>i!lI;:,\"^`.",
	}
	img, err := loadImage(options.Path)
	if err != nil {
		b.Fatal(err)
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		_ = process(img, options)
	}
}

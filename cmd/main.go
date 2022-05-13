package main

import (
	"flag"
	"fmt"

	"github.com/WRWPhillips/go-pic2text/internal"
)

func main() {
	path := flag.String("input", "assets/input.jpeg", "Input file")
	width := flag.Int("width", 80, "Output width")
	height := flag.Int("height", 25, "Output height")
	flag.Parse()

	options := internal.Options{
		Path:   *path,
		Width:  *width,
		Height: *height,
		//Palette: "0123456789012345678901234567890123456789012345678901234567890123456789",
		Palette: "'$@B%8&WM#*oahkbdpqwmZO0QLCJUYXzcvunxrjft/\\|()1{}[]?-_+~<>i!lI;:,\"^`.",
	}

	fmt.Println(options)
}

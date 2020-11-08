// Package image allows parsing an image which will later be used as a key
package image

import (
	gi "image"
	"image/color"
	"image/png"
	"math/rand"
	"os"
	"sync"
)

// Dimension contains infos about the image dimension
type Dimension struct {
	Width  int
	Height int
}

// Image contains the datastructure representing a greyscale image
type Image struct {
	sync.Mutex
	// Data contains the raw byte values
	Data []uint8 // addressing: Data[width + Dimension.Width * height]
	// Dimension contains the image dimensions
	Dimension Dimension
}

// Read supports loading a PNG file (greyscale)
func Read(f *os.File) (*Image, error) {
	img, err := png.Decode(f)
	if err != nil {
		return nil, err
	}
	b := img.Bounds()

	// Generate mock image
	i := &Image{
		Data:      make([]uint8, b.Dx()*b.Dy(), b.Dx()*b.Dy()),
		Dimension: Dimension{Width: b.Dx(), Height: b.Dy()},
	}

	// Fill mock image with mock data
	for h := 0; h < i.Dimension.Height; h++ {
		for w := 0; w < i.Dimension.Width; w++ {
			c := img.At(w, h)
			pixelValue, _, _, _ := c.RGBA()
			i.Data[w+i.Dimension.Width*h] = uint8(uint8(pixelValue))
		}
	}

	return i, nil
}

// Write supports writing a PNG file (greyscale)
func Write(f *os.File, i *Image) error {
	img := gi.NewGray(gi.Rectangle{
		Min: gi.Point{0, 0},
		Max: gi.Point{i.Dimension.Width, i.Dimension.Height},
	})

	for h := 0; h < i.Dimension.Height; h++ {
		for w := 0; w < i.Dimension.Width; w++ {
			c := color.Gray{i.Data[w+i.Dimension.Width*h]}
			img.Set(w, h, c)
		}
	}

	png.Encode(f, img)
	return nil
}

// Mock creates an image with 128x128 dimension for testing purposes
func Mock() *Image {
	// Generate mock image
	i := &Image{
		Data:      make([]uint8, 128*128, 128*128),
		Dimension: Dimension{128, 128},
	}

	// Fill mock image with mock data
	for h := 0; h < i.Dimension.Height; h++ {
		for w := 0; w < i.Dimension.Width; w++ {
			i.Data[w+i.Dimension.Width*h] = uint8(rand.Intn(128))
		}
	}

	return i
}

// CheckAccept check if we can map all characters on the image -> determine if it is suited or not.
func CheckAccept(i *Image) bool {
	acceptanceMap := make(map[uint8]bool)
	for _, b := range i.Data {
		acceptanceMap[b&0b01111111] = true
	}

	// Check if ASCII alphabet can be represented
	for c := 0; c < 128; c++ {
		if _, ok := acceptanceMap[uint8(c)]; !ok {
			return false
		}
	}

	return true
}

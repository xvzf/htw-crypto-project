// Package image allows parsing an image which will later be used as a key
package image

import (
	"math/rand"
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

// Mock creates an image with 1024x1024 dimension for testing purposes
func Mock() *Image {
	// Generate mock image
	i := &Image{
		Data:      make([]uint8, 1024*1024, 1024*1024),
		Dimension: Dimension{1024, 1024},
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

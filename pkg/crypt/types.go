package crypt

import (
	"github.com/xvzf/htw-crypto-project/pkg/image"
)

// PixelPosition represents a pixels position on the image
type PixelPosition struct {
	w int
	h int
}

// PixelGroups represents a grouping of pixels based on their value
type PixelGroups map[uint8][]PixelPosition

// Container contains neccessary infos for encrypting and decrypting
type Container struct {
	Image       *image.Image
	PixelGroups PixelGroups
}

// Encrypted contains a slice of PixelPositions
type Encrypted []PixelPosition

// Package crypt implements the symmetric encryption proposed in this paper:
// https://ieeexplore.ieee.org/document/7420966
package crypt

import (
	"errors"
	"math/rand"

	"github.com/xvzf/htw-crypto-project/pkg/image"
)

// ExtractGroups extract Pixel groups of an image
func ExtractGroups(i *image.Image) PixelGroups {

	p := make(PixelGroups)

	// Iterate through all pixels & update the grouping
	for h := 0; h < i.Dimension.Height; h++ {
		for w := 0; w < i.Dimension.Width; w++ {
			v := i.Data[w+i.Dimension.Width*h] & 0b01111111
			if _, ok := p[v]; !ok {
				p[v] = []PixelPosition{{w, h}}
			} else {
				tmp, _ := p[v]
				p[v] = append(tmp, PixelPosition{w, h})
			}
		}
	}

	return p
}

// New creates a new container allowing for endcryption&decryption
func New(i *image.Image) (*Container, error) {
	if !image.CheckAccept(i) {
		return nil, errors.New("Image not suiteable")
	}
	c := &Container{
		Image:       i,
		PixelGroups: ExtractGroups(i),
	}

	return c, nil
}

// Encrypt allows encryption of an arbitrary ASCII string
func (c *Container) Encrypt(s string) (Encrypted, error) {
	enc := make(Encrypted, len(s))

	// Iterate over the input string, determine (random) pixel position
	for i, b := range []uint8(s) {
		if pixelGroup, ok := c.PixelGroups[b]; ok {
			// Get the number of available options for the pixel value
			availOptions := len(pixelGroup)
			// Choose a random position out of the pixel group
			ppos := pixelGroup[rand.Intn(availOptions)]
			enc[i] = ppos
		}
	}

	return enc, nil
}

// Decrypt allows decryption of an arbitrary encrypted ASCII string
func (c *Container) Decrypt(enc Encrypted) (string, error) {
	dec := make([]byte, len(enc))

	for i, ec := range enc {
		// Check if in boundaries
		if ec.h < c.Image.Dimension.Height && ec.w < c.Image.Dimension.Width {
			// Calculate pixel position in the slice
			arrayPos := ec.w + c.Image.Dimension.Width*ec.h
			// Retrieve Byte
			dec[i] = byte(c.Image.Data[arrayPos] & 0b01111111)
		} else {
			// Invalid pixel position
			return "", errors.New("Invalid pixel position")
		}
	}

	// Converte encrypted bytes to string & return
	return string(dec), nil
}

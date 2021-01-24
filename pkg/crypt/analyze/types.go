package crypt

import "github.com/xvzf/htw-crypto-project/pkg/crypt"

type PixelPositionFrequency struct {
	PixelPosition crypt.PixelPosition
	Count         int
}

type PixelGroupFrequency struct {
	PixelPositions []crypt.PixelPosition
	Total          int
}

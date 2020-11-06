package crypt

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xvzf/htw-crypto-project/pkg/image"
)

func Test_ExtractGroups(t *testing.T) {

	// Generate valid test image
	i := image.Mock()
	for !image.CheckAccept(i) {
		i = image.Mock()
	}

	// Extract groups out of the image
	groups := ExtractGroups(i)

	// Check if a group exists for all ASCII characters
	for c := 0; c < 128; c++ {
		// retrieve PixelPosition group
		v, ok := groups[uint8(c)]
		assert.True(t, ok)
		// Check if there is more than one entry existing
		assert.Greater(t, len(v), 0)
	}
}

func TestContainer_Encrypt_Decrypt(t *testing.T) {
	// Generate valid test image
	i := image.Mock()
	for !image.CheckAccept(i) {
		i = image.Mock()
	}

	// Build new crypt engine
	cipher, err := New(i)
	assert.NoError(t, err)

	genAllASCII := func() string {
		buf := make([]byte, 128)

		for i := 0; i < 128; i++ {
			buf[i] = byte(i)
		}

		return string(buf)
	}

	for _, s := range []string{
		"hello",
		"world",
		"abcdefhijklmnopqrstuvwxyz",
		"this123isatestencryotun12523sa[@$!@{K352!21]5'|,125}]",
		genAllASCII(),
	} {
		// Encrypt
		enc, err := cipher.Encrypt(s)
		assert.NoError(t, err)
		// Decrypt
		dec, err := cipher.Decrypt(enc)
		assert.NoError(t, err)

		// assert decrypted is the same
		assert.Equal(t, s, dec)
	}
}

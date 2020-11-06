package crypt

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xvzf/htw-crypto-project/pkg/image"
)

var (
	cipher            *Container
	testData1MByte    string
	testData1MByteEnc []PixelPosition
)

func init() {

	var err error

	i := image.Mock()
	for !image.CheckAccept(i) {
		i = image.Mock()
	}

	// Build new crypt engine
	cipher, err = New(i)
	if err != nil {
		log.Fatal(err)
	}

	// initialize testdata
	// We want 1 MByte testdata for benchmarking
	buf := make([]byte, 1024*1024)
	for i, _ := range buf {
		buf[i] = byte(i % 128)
	}

	testData1MByte = string(buf)
	testData1MByteEnc, _ = cipher.Encrypt(testData1MByte)
}

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
	for _, s := range []string{
		"hello",
		"world",
		"abcdefhijklmnopqrstuvwxyz",
		"this123isatestencryotun12523sa[@$!@{K352!21]5'|,125}]",
		testData1MByte,
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

func BenchmarkContainer_Encrypt1MByte(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			cipher.Encrypt(testData1MByte)
		}
	})
}

func BenchmarkContainer_Encrypt1024(b *testing.B) {
	for n := 0; n < b.N; n++ {
		cipher.Encrypt(testData1MByte)
	}
}

func BenchmarkContainer_Decrypt1024(b *testing.B) {
	for n := 0; n < b.N; n++ {
		cipher.Decrypt(testData1MByteEnc)
	}
}

func BenchmarkContainer_Decrypt1MByte(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			cipher.Decrypt(testData1MByteEnc)
		}
	})
}

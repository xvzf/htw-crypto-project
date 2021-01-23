# Implementation And Analysis Of A Polyalphabetic Substitution Cipher Including Attacks


## Details on the cipher.

The cipher presented in [*A New Symmetric Key Encryption Algorithm Using Images as Secret Keys*](https://ieeexplore.ieee.org/document/7420966, University Project HTW Saar) encrypts arbitrary ASCII data and uses a given image as a shared secret.
It substitutes the ASCII to a ciphertext by mapping characters to pixel positions, thus generating a monoalphabetic substitution cipher.

## Implementation

During this project, the described cipher has been implemented using the [Go programming language](https://golang.org/) aiming at a high encryption and decryption performance.

### Key Verification

As described in the underlying paper, the substitution alphabet is derived from a given image source.
The `pkg/image` package implements an interface to read and write arbitrary images used as keys and also provides an acceptance test checking the validity using an image as a key for the cipher.

The cipher requires a map of the (7bits) ASCII alphabet to a group of pixel positions representing each character. The mapping shall be one-way inversible meaning `a <=> any of { 1, 2, 3, 4, 5 }`.
In order for an image being accepted as a valid key, all ASCII characters must have at least one pixel representing the character.
In order to allow the one-way inversion of character to pixel mapping, the pixel value is ANDed with `0b01111111` thus elimiting the most significant bit.
Once an image fulfills the requirement of mapping all ASCII characters to one or many pixels, it is accepted.

```go
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
```

For testing purposes, the function `Mock()` generates an image with 128 pixels width and 128 height containing random pixel values. This iamge can be verified to be a valid key for the cipher:
```go
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
```

### Encryption

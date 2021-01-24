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
For every character, the encryption chooses a random pixel representing the character.

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
Before the encryption starts, the key is converted into an optimized datastructure implemented by a hash-map with the character to be encrypted as key and an array of pixel positions representing this character.
```go
// PixelPosition represents a pixels position on the image
type PixelPosition struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

// PixelGroups represents a grouping of pixels based on their value
type PixelGroups map[uint8][]PixelPosition
```
The key, an abitrary image, is loaded using the aforementioned `pkg/image` package and checked for validity.
After it passes the test, a trivial iteration across all pixels helps generating the `PixelGroups` datastructure.
```go
// Encrypt allows encryption of an arbitrary ASCII string
func (c *Container) Encrypt(s string) (Encrypted, error) {
	enc := make(Encrypted, len(s))
	rnd := make([]byte, 4)

	// Iterate over the input string, determine (random) pixel position
	for i, b := range []uint8(s) {
		if pixelGroup, ok := c.PixelGroups[b]; ok {
			// Get the number of available options for the pixel value
			availOptions := len(pixelGroup)
			// Choose a random position out of the pixel group
			d := int(binary.BigEndian.Uint32(rnd)) % availOptions

			ppos := pixelGroup[d]
			enc[i] = ppos
		}
	}

	return enc, nil
}
```

During the encryption process, a loop iterates across the plaintext and determines the substitution for every character by using a random entry out of the array mapped for the character.
```go
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
```

### Decryption
The Decryption builds on a lookup function.
A loop iterates across the ciphertext and retrieves the pixel value at the position the ciphertext references. This value is ANDed by `0b01111111` for retrieving the plaintext character.

```go
// Decrypt allows decryption of an arbitrary encrypted ASCII string
func (c *Container) Decrypt(enc Encrypted) (string, error) {
	dec := make([]byte, len(enc))

	for i, ec := range enc {
		// Check if in boundaries
		if ec.Height < c.Image.Dimension.Height && ec.Width < c.Image.Dimension.Width {
			// Calculate pixel position in the slice
			arrayPos := ec.Width + c.Image.Dimension.Width*ec.Height
			// Retrieve Byte
			dec[i] = byte(c.Image.Data[arrayPos] & 0b01111111)
		} else {
			// Invalid pixel position
			return "", errors.New("Invalid pixel position")
		}
	}

	// Convert encrypted bytes to string & return
	return string(dec), nil
}
```

### Input/Output operations
For the ease of analysing the ciphertext with different tools, a JSON format has been chosen which encodes an array of `PixelPosition`.
```go
// Encrypted contains a slice of PixelPositions
type Encrypted []PixelPosition

// Reads a ciphertext into an internal datastructure
func Read(f *os.File) ([]PixelPosition, error) {
	var out []PixelPosition
	dec := json.NewDecoder(f)
	// ...
	return out, nil
}

// Encodes ciphertext as JSON
func Write(f *os.File, in []PixelPosition) error {
	enc := json.NewEncoder(f)
	if err := enc.Encode(in); err != nil {
		return err
	}
	return nil
}
```

### Testing The Implementation

To ensure the implementation is working and can encrypt & afterwards decrypt ciphertext with a shared key, test cases have been implemented with the Go integrated unit-test framework.
As an initial step, an initialization function generates a mocked image which passes the acceptance criteria. Afterwards, a random plaintext is generated with the size of 1 Mbyte.

Several test cases check the implementation for the desired functionality. the `Test_ExtractGroups` verifies the key extraction for an image is successful when the image is valid.
The `TestContainer_Encrypt_Decrypt` case verifies encrypted plaintext can be reverted to the original plaintext by using the same key.


### Performance Testing The Implementation

The performance of the encryption & decryption operations is implemented using the Go internal benchmarking framework.
Different testcases for parallel as well as synchronous chunk decryption and encryption are used to compute the implementation performance.

Benchmark implementation:
```go
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
```

Benchmark results (*Intel i7-9750H*):
```bash
xvzf in crypto/pkg/crypt on ? master ❯ go test -bench .
goos: darwin
goarch: amd64
pkg: github.com/xvzf/htw-crypto-project/pkg/crypt
BenchmarkContainer_Encrypt1Mbyte               	      22	  50381141 ns/op	16777216 B/op	       1 allocs/op
BenchmarkContainer_Encrypt1Mbyte-6             	      21	  52016303 ns/op	16777230 B/op	       1 allocs/op
BenchmarkContainer_Encrypt1Mbyte-12            	      21	  51928491 ns/op	16777225 B/op	       1 allocs/op
BenchmarkContainer_Decrypt1Mbyte               	     559	   2031963 ns/op	 2097152 B/op	       2 allocs/op
BenchmarkContainer_Decrypt1Mbyte-6             	     560	   2022927 ns/op	 2097152 B/op	       2 allocs/op
BenchmarkContainer_Decrypt1Mbyte-12            	     578	   2020687 ns/op	 2097152 B/op	       2 allocs/op
BenchmarkContainer_Encrypt1MByteParallel       	      20	  50529087 ns/op	16777224 B/op	       1 allocs/op
BenchmarkContainer_Encrypt1MByteParallel-6     	      21	 114059236 ns/op	16777276 B/op	       1 allocs/op
BenchmarkContainer_Encrypt1MByteParallel-12    	      21	 134402985 ns/op	16777444 B/op	       2 allocs/op
BenchmarkContainer_Decrypt1MByteParallel       	     584	   2014307 ns/op	 2097152 B/op	       2 allocs/op
BenchmarkContainer_Decrypt1MByteParallel-6     	    2061	    527559 ns/op	 2097154 B/op	       2 allocs/op
BenchmarkContainer_Decrypt1MByteParallel-12    	    2544	    456374 ns/op	 2097156 B/op	       2 allocs/op
PASS
ok  	github.com/xvzf/htw-crypto-project/pkg/crypt	20.776s
```

Extrapolating the benchmark results of `BenchmarkContainer_Encrypt1Mbyte`, on an average of 21 runs, the encryption duration of 1 MB data is `50381141 ns` which represents a throughput of *19.85 MB/s*.
As seen in the results of the `BenchmarkContainer_Encrypt1MByteParallel-*` scenarios, running the encryption with multiple cores hurts the performance. This is likely due to the implementation and usage of the integrated random number generator.

The encryption however benefits from running it in parallell and decrypts at `492 MB/s` with a single CPU core and `1.9 GB/s` with 6 cores.
Adding the hyper-threading of the Intel processor results in a further performance boost to `2.19 GB/s`.


## Security Analysis
As described in the paper, the cipher maps one characters to possibly multiple pixels in an image.
While this seems to be secure it is under some circumstances vulnerable to a frequency analysis.
Based on the plaintext length, the image dimenson, pixel distribution mapping to characters in the image, and the character frequency of the plaintext it is possible to build groups of pixels based on their occurence assuming the random generator is uniform.


The security of the cipher therefore drastically depends on the plaintext lenght and the image size.
Given a frequency vulnerable plaintext with a lenght of *10.000* characters containing *50* unique characters, the image used as key decides whether the cipher can be cracked or not:
- The image contains distinguished pixels for all characters, thus the ciphertext does not contain double values; **OTP like security**
- The image does not contain distinguished pixels for all characters, resulting ciphertext values cannot be grouped based on their frequency; **Vulnerable to additional attacks, discussid in TODO**
- The image does not contain distinguished pixels for all characters, resulting ciphertext values can be grouped based on their frequency; **Allows further analysis on simple substitution cipher**

### Grouping Pixels
Given an image that allows grouping of pixels, the following algorithm

### Advanced Attack: 1 plaintext - n ciphertext
Since the encryption is not injecting a nonce, always generating a unique plaintext, this attack scheme aims at analysing different ciphertext based on one plaintext.
Especially on large plain/ciphertexts this allows an easy grouping of pixels unlinked to their occurence. A common battern is e.g. an e-mail header.

This breaks down the cipher to a simple substitution cipher which can be cracked e.g. using a trivial frequency analysis


### Advanced Attack: Adding known-plaintext attack
The aforementioned attack allows us to group pixels based on many ciphertexts.
While a simple freqency analysis can be effective solving the cipher, a known plaintext-ciphertext pair helps cracking a plaintext which is not vulnerable to frequency analysis.


### Extracting the key image

It is not possible to extract the original image as the pixel value space is reduced from 8 to 7 bits.
However, mapping the cracked key to an image may help to identify structures of the original image while generating a valid key.

package crypt

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrite(t *testing.T) {
	f, err := ioutil.TempFile("/tmp", "crypt_io_write")
	assert.NoError(t, err)
	defer f.Close()

	err = Write(f, testData1MByteEnc)
	assert.NoError(t, err)
}

func TestRead(t *testing.T) {
	f, err := ioutil.TempFile("/tmp", "crypt_io_write")
	assert.NoError(t, err)

	err = Write(f, testData1MByteEnc)
	assert.NoError(t, err)
	f.Close()

	f, err = os.Open(f.Name())
	assert.NoError(t, err)

	v, err := Read(f)
	assert.Equal(t, testData1MByteEnc, v)
}

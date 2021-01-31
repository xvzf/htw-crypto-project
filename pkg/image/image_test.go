package image

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMockAcceptance(t *testing.T) {
	i := Mock()
	assert.True(t, CheckAccept(i))
}

// Test Image Writing capabilities
func TestWrite(t *testing.T) {
	i := Mock()
	f, err := ioutil.TempFile("/tmp", "imgwrite")
	assert.NoError(t, err)
	defer f.Close()

	Write(f, i)

}

// Test image read capabilities
func TestRead(t *testing.T) {
	i := Mock()
	f, err := ioutil.TempFile("/tmp", "imgread")
	assert.NoError(t, err)
	Write(f, i)
	f.Close()

	f, err = os.Open(f.Name())
	assert.NoError(t, err)
	defer f.Close()

	n, err := Read(f)
	assert.NoError(t, err)
	assert.Equal(t, i, n)
}

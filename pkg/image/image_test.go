package image

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMockAcceptance(t *testing.T) {
	i := Mock()
	assert.True(t, CheckAccept(i))
}

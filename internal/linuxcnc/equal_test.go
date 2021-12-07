package linuxcnc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEqual(t *testing.T) {
	s1 := Status{}
	s2 := Status{}

	assert.True(t, s1.Equal(&s2))
}

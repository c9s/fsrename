package fsrename

import "testing"
import "github.com/stretchr/testify/assert"

func TestSequence(t *testing.T) {
	seq := NewSequence(0, 1)
	var id int32
	id = seq.Next()
	assert.Equal(t, id, int32(1))

	id = seq.Next()
	assert.Equal(t, id, int32(2))

	id = seq.Next()
	assert.Equal(t, id, int32(3))

	seq.Reset()
	id = seq.Next()
	assert.Equal(t, id, int32(1))

	id = seq.Next()
	assert.Equal(t, id, int32(2))
}

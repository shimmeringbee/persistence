package converter

import (
	"github.com/shimmeringbee/persistence/impl/memory"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRetrieve(t *testing.T) {
	t.Run("default value provided is returned if not found and default provided", func(t *testing.T) {
		s := memory.New()

		expected := time.Duration(1)

		actual, found := Retrieve(s, Key, DurationDecoder, expected)
		assert.False(t, found)
		assert.Equal(t, expected, actual)
	})

	t.Run("zero value provided is returned if not found and no default", func(t *testing.T) {
		s := memory.New()

		expected := time.Duration(0)

		actual, found := Retrieve(s, Key, DurationDecoder)
		assert.False(t, found)
		assert.Equal(t, expected, actual)
	})
}

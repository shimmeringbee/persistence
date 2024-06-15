package converter

import (
	"github.com/shimmeringbee/persistence/impl/memory"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const Key = "key"

func TestTime(t *testing.T) {
	t.Run("time is stored and retrieved to the millisecond level", func(t *testing.T) {
		s := memory.New()

		expected := time.UnixMilli(time.Now().UnixMilli())

		Store(s, Key, expected, TimeEncoder)

		actual, found := Retrieve(s, Key, TimeDecoder)
		assert.True(t, found)
		assert.Equal(t, expected, actual)
	})
}

func TestDuration(t *testing.T) {
	t.Run("duration is stored and retrieved to the millisecond level", func(t *testing.T) {
		s := memory.New()

		expected := time.Duration(1234) * time.Millisecond

		Store(s, Key, expected, DurationEncoder)

		actual, found := Retrieve(s, Key, DurationDecoder)
		assert.True(t, found)
		assert.Equal(t, expected, actual)
	})
}

package memory

import (
	"github.com/shimmeringbee/persistence/impl/test"
	"testing"
)

func TestMemory(t *testing.T) {
	test.Impl{New: New}.Test(t)
}

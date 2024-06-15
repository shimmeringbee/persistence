package memory

import (
	"github.com/shimmeringbee/persistence/impl/test"
	"testing"
)

func TestMemory(t *testing.T) {
	test.Impl{
		New:    New,
		Done:   test.EmptyDone,
		Switch: test.EmptySwitch,
	}.Test(t)
}

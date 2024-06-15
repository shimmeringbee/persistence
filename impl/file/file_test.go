package file

import (
	"github.com/shimmeringbee/persistence"
	"github.com/shimmeringbee/persistence/impl/test"
	"os"
	"sync"
	"testing"
)

type tracker struct {
	m  *sync.Mutex
	db map[persistence.Section]string
}

func (t *tracker) New() persistence.Section {
	t.m.Lock()
	defer t.m.Unlock()

	dir, err := os.MkdirTemp("", "*")
	if err != nil {
		panic(err)
	}

	return t.new(dir)
}

func (t *tracker) new(dir string) persistence.Section {
	p := New(dir)
	t.db[p] = dir

	return p
}

func (t *tracker) Switch(p persistence.Section) persistence.Section {
	t.m.Lock()
	defer t.m.Unlock()

	dir, ok := t.db[p]
	if !ok {
		panic("switch called on non existent persistence")
	}

	if f, ok := p.(*file); ok {
		f.Sync()
	}

	return t.new(dir)
}

func (t *tracker) Done(p persistence.Section) {
	t.m.Lock()
	defer t.m.Unlock()

	dir, ok := t.db[p]
	if !ok {
		panic("switch called on non existent persistence")
	}

	delete(t.db, p)

	if err := os.RemoveAll(dir); err != nil {
		panic(err)
	}
}

func TestFile(t *testing.T) {
	tr := tracker{m: &sync.Mutex{}, db: make(map[persistence.Section]string)}
	test.Impl{New: tr.New, Switch: tr.Switch, Done: tr.Done}.Test(t)
}

package file

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/shimmeringbee/persistence"
	"github.com/shimmeringbee/persistence/impl/memory"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

func New(dir string) persistence.Section {
	f := &file{m: &sync.RWMutex{}, cache: memory.New(), sections: make(map[string]*file)}

	dirWithoutPathSep, _ := strings.CutSuffix(dir, string(os.PathSeparator))
	f.dir = fmt.Sprintf("%s%c", dirWithoutPathSep, os.PathSeparator)
	_ = os.MkdirAll(f.dir, 600)

	f.load()

	return f
}

type file struct {
	dir   string
	cache persistence.Section

	m        *sync.RWMutex
	sections map[string]*file

	dirtyTimer *time.Timer
}

var _ persistence.Section = (*file)(nil)

func (f *file) Section(key ...string) persistence.Section {
	f.m.Lock()
	defer f.m.Unlock()

	s, ok := f.sections[key[0]]

	if !ok {
		dir := fmt.Sprintf("%s%s", f.dir, key[0])

		s, _ = New(dir).(*file)
		f.sections[key[0]] = s
	}

	if len(key) > 1 {
		return s.Section(key[1:]...)
	} else {
		return s
	}
}

func (f *file) SectionKeys() []string {
	f.m.RLock()
	defer f.m.RUnlock()

	var keys = make([]string, 0, len(f.sections))

	for k := range f.sections {
		keys = append(keys, k)
	}

	return keys
}

func (f *file) SectionExists(key string) bool {
	f.m.RLock()
	defer f.m.RUnlock()

	_, found := f.sections[key]

	return found
}

func (f *file) SectionDelete(key string) bool {
	f.m.Lock()
	defer f.m.Unlock()

	if s, ok := f.sections[key]; ok {
		for _, k := range s.SectionKeys() {
			s.SectionDelete(k)
		}

		s.sectionDeleteSelf()
		delete(f.sections, key)
		return true
	} else {
		return false
	}
}

func (f *file) sectionDeleteSelf() {
	_ = os.RemoveAll(f.dir)
}

func (f *file) Keys() []string {
	return f.cache.Keys()
}

func (f *file) Exists(key string) bool {
	return f.cache.Exists(key)
}

func (f *file) Type(key string) persistence.ValueType {
	return f.cache.Type(key)
}

func (f *file) Int(key string, defValue ...int64) (int64, bool) {
	return f.cache.Int(key, defValue...)
}

func (f *file) UInt(key string, defValue ...uint64) (uint64, bool) {
	return f.cache.UInt(key, defValue...)
}

func (f *file) String(key string, defValue ...string) (string, bool) {
	return f.cache.String(key, defValue...)
}

func (f *file) Bool(key string, defValue ...bool) (bool, bool) {
	return f.cache.Bool(key, defValue...)
}

func (f *file) Float(key string, defValue ...float64) (float64, bool) {
	return f.cache.Float(key, defValue...)
}

func (f *file) Bytes(key string, defValue ...[]byte) ([]byte, bool) {
	return f.cache.Bytes(key, defValue...)
}

func (f *file) Set(key string, value interface{}) {
	f.cache.Set(key, value)
	f.dirty()
}

func (f *file) Delete(key string) bool {
	ok := f.cache.Delete(key)
	f.dirty()
	return ok
}

const dataFile = "data.json"

type Value struct {
	Value any
	Type  persistence.ValueType
}

func (f *file) load() {
	var sections []string
	dataPresent := false

	entries, err := os.ReadDir(f.dir)
	if err != nil {
		panic(err)
	}

	for _, ent := range entries {
		if ent.IsDir() {
			sections = append(sections, ent.Name())
		} else if ent.Name() == dataFile {
			dataPresent = true
		}
	}

	if dataPresent {
		var d map[string]Value

		b, err := os.ReadFile(fmt.Sprintf("%s%s", f.dir, dataFile))
		if err != nil {
			panic(err)
		}

		dec := json.NewDecoder(bytes.NewReader(b))
		dec.UseNumber()

		if err := dec.Decode(&d); err != nil {
			panic(err)
		}

		for k, v := range d {
			f.loadValue(k, v)
		}
	}

	for _, subsection := range sections {
		s := New(fmt.Sprintf("%s%s", f.dir, subsection)).(*file)
		f.sections[subsection] = s
	}
}

const dirtyDelay = 500 * time.Millisecond

func (f *file) dirty() {
	f.m.Lock()
	defer f.m.Unlock()

	if f.dirtyTimer != nil {
		f.dirtyTimer.Stop()
	}

	f.dirtyTimer = time.AfterFunc(dirtyDelay, f.dirtySync)
}

func (f *file) dirtySync() {
	f.m.Lock()

	if f.dirtyTimer != nil {
		f.dirtyTimer.Stop()
		f.dirtyTimer = nil
	}

	f.m.Unlock()

	f.sync(false)
}

func (f *file) loadValue(k string, v Value) {
	switch v.Type {
	case persistence.Int:
		if jn, ok := v.Value.(json.Number); ok {
			if n, err := strconv.Atoi(string(jn)); err == nil {
				f.cache.Set(k, int64(n))
			}
		}
	case persistence.UnsignedInt:
		if jn, ok := v.Value.(json.Number); ok {
			if n, err := strconv.Atoi(string(jn)); err == nil {
				f.cache.Set(k, uint64(n))
			}
		}
	case persistence.String:
		if s, ok := v.Value.(string); ok {
			f.cache.Set(k, s)
		}
	case persistence.Bool:
		if b, ok := v.Value.(bool); ok {
			f.cache.Set(k, b)
		}
	case persistence.Float:
		if jn, ok := v.Value.(json.Number); ok {
			if n, err := strconv.ParseFloat(string(jn), 64); err == nil {
				f.cache.Set(k, n)
			}
		}
	case persistence.Bytes:
		if ba, ok := v.Value.(string); ok {
			var data []byte

			for ; len(ba) > 0; ba = ba[2:] {
				if b, err := strconv.ParseInt(ba[:2], 16, 8); err == nil {
					data = append(data, byte(b))
				} else {
					return
				}
			}

			f.cache.Set(k, data)
		}
	}
}

func (f *file) Sync() {
	f.sync(true)
}

func (f *file) sync(recursive bool) {
	f.m.RLock()
	defer f.m.RUnlock()

	data := make(map[string]Value)

	for _, k := range f.cache.Keys() {
		var v any
		t := f.cache.Type(k)

		switch t {
		case persistence.Int:
			v, _ = f.cache.Int(k)
		case persistence.UnsignedInt:
			v, _ = f.cache.UInt(k)
		case persistence.String:
			v, _ = f.cache.String(k)
		case persistence.Bool:
			v, _ = f.cache.Bool(k)
		case persistence.Float:
			v, _ = f.cache.Float(k)
		case persistence.Bytes:
			bs, _ := f.cache.Bytes(k)

			var bytesOut []string

			for _, b := range bs {
				bytesOut = append(bytesOut, fmt.Sprintf("%02x", b))
			}

			v = strings.Join(bytesOut, "")
		}

		data[k] = Value{
			Value: v,
			Type:  t,
		}
	}

	r, err := os.Create(fmt.Sprintf("%s%s", f.dir, dataFile))
	if err != nil {
		panic(err)
	}
	defer r.Close()

	enc := json.NewEncoder(r)
	enc.SetIndent("", "  ")

	if err := enc.Encode(data); err != nil {
		panic(err)
	}

	if recursive {
		for _, v := range f.sections {
			v.sync(recursive)
		}
	}
}

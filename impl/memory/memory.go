package memory

import (
	"fmt"
	"github.com/shimmeringbee/persistence"
	"sync"
)

func New() persistence.Section {
	return &memory{m: &sync.RWMutex{}, kv: make(map[string]interface{}), sections: make(map[string]persistence.Section)}
}

type memory struct {
	m        *sync.RWMutex
	kv       map[string]interface{}
	sections map[string]persistence.Section
}

func (m *memory) SectionExists(key string) bool {
	m.m.RLock()
	defer m.m.RUnlock()

	_, found := m.sections[key]

	return found
}

func (m *memory) Section(key ...string) persistence.Section {
	m.m.RLock()
	s, ok := m.sections[key[0]]
	m.m.RUnlock()

	if !ok {
		s = New()
		m.m.Lock()
		m.sections[key[0]] = s
		m.m.Unlock()
	}

	if len(key) > 1 {
		return s.Section(key[1:]...)
	} else {
		return s
	}
}

func (m *memory) SectionKeys() []string {
	m.m.RLock()
	defer m.m.RUnlock()

	var keys = make([]string, 0, len(m.sections))

	for k := range m.sections {
		keys = append(keys, k)
	}

	return keys
}

func (m *memory) SectionDelete(key string) bool {
	m.m.Lock()
	defer m.m.Unlock()

	_, found := m.sections[key]

	if found {
		delete(m.sections, key)
	}

	return found
}

func (m *memory) Exists(key string) bool {
	m.m.RLock()
	defer m.m.RUnlock()

	_, found := m.kv[key]

	return found
}

func (m *memory) Keys() []string {
	m.m.RLock()
	defer m.m.RUnlock()

	var keys = make([]string, 0, len(m.kv))

	for k := range m.kv {
		keys = append(keys, k)
	}

	return keys
}

func genericRetrieve[T any](m *memory, key string, defValue ...T) (T, bool) {
	m.m.RLock()
	defer m.m.RUnlock()

	v, ok := m.kv[key]

	if ok {
		if iV, cok := v.(T); cok {
			return iV, true
		}
	}

	if len(defValue) > 0 {
		return defValue[0], false
	} else {
		zero := *new(T)
		return zero, false
	}
}

func (m *memory) Int(key string, defValue ...int64) (int64, bool) {
	return genericRetrieve(m, key, defValue...)
}

func (m *memory) UInt(key string, defValue ...uint64) (uint64, bool) {
	return genericRetrieve(m, key, defValue...)
}

func (m *memory) String(key string, defValue ...string) (string, bool) {
	return genericRetrieve(m, key, defValue...)
}

func (m *memory) Bool(key string, defValue ...bool) (bool, bool) {
	return genericRetrieve(m, key, defValue...)
}

func (m *memory) Float(key string, defValue ...float64) (float64, bool) {
	return genericRetrieve(m, key, defValue...)
}

func (m *memory) Bytes(key string, defValue ...[]byte) ([]byte, bool) {
	return genericRetrieve(m, key, defValue...)
}

func (m *memory) Set(key string, value interface{}) error {
	var sV interface{}

	switch v := value.(type) {
	case string:
		sV = v
	case int:
		sV = int64(v)
	case int8:
		sV = int64(v)
	case int16:
		sV = int64(v)
	case int32:
		sV = int64(v)
	case int64:
		sV = v
	case uint:
		sV = uint64(v)
	case uint8:
		sV = uint64(v)
	case uint16:
		sV = uint64(v)
	case uint32:
		sV = uint64(v)
	case uint64:
		sV = v
	case float32:
		sV = float64(v)
	case float64:
		sV = v
	case bool:
		sV = v
	case []byte:
		sV = v
	default:
		return fmt.Errorf("section set: unknown type: %T", v)
	}

	m.m.Lock()
	defer m.m.Unlock()

	m.kv[key] = sV

	return nil
}

func (m *memory) Delete(key string) bool {
	m.m.Lock()
	defer m.m.Unlock()

	_, found := m.kv[key]

	if found {
		delete(m.kv, key)
	}

	return found
}

var _ persistence.Section = (*memory)(nil)

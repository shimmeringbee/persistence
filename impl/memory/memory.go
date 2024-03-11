package memory

import (
	"fmt"
	"github.com/shimmeringbee/persistence"
	"sync"
)

func New() persistence.Section {
	return &memory{m: &sync.RWMutex{}, kv: make(map[string]interface{})}
}

type memory struct {
	m  *sync.RWMutex
	kv map[string]interface{}
}

func (m *memory) Exists(key string) bool {
	m.m.RLock()
	defer m.m.RUnlock()

	_, ok := m.kv[key]

	return ok
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

func (m *memory) Section(key ...string) persistence.Section {
	m.m.RLock()
	v, ok := m.kv[key[0]]
	m.m.RUnlock()

	var s persistence.Section

	if ok {
		if cs, cok := v.(persistence.Section); cok {
			s = cs
		} else {
			ok = false
		}
	}

	if !ok {
		s = New()
		m.m.Lock()
		m.kv[key[0]] = s
		m.m.Unlock()
	}

	if len(key) > 1 {
		return s.Section(key[1:]...)
	} else {
		return s
	}
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

func (m *memory) Int(key string, defValue ...int) (int, bool) {
	return genericRetrieve(m, key, defValue...)
}

func (m *memory) UInt(key string, defValue ...uint) (uint, bool) {
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
		sV = v
	case int8:
		sV = int(v)
	case int16:
		sV = int(v)
	case int32:
		sV = int(v)
	case int64:
		sV = int(v)
	case uint:
		sV = v
	case uint8:
		sV = uint(v)
	case uint16:
		sV = uint(v)
	case uint32:
		sV = uint(v)
	case uint64:
		sV = uint(v)
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

func (m *memory) Delete(key string) {
	m.m.Lock()
	defer m.m.Unlock()
	delete(m.kv, key)
}

var _ persistence.Section = (*memory)(nil)

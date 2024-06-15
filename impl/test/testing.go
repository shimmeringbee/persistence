package test

import (
	"github.com/shimmeringbee/persistence"
	"github.com/stretchr/testify/assert"
	"testing"
)

func EmptySwitch(p persistence.Section) persistence.Section {
	return p
}

func EmptyDone(_ persistence.Section) {}

type Impl struct {
	New    func() persistence.Section
	Switch func(persistence.Section) persistence.Section
	Done   func(persistence.Section)
}

func (tt Impl) Test(t *testing.T) {
	for name, test := range map[string]func(*testing.T){
		"Keys":               tt.Keys,
		"Type":               tt.Type,
		"Delete":             tt.Delete,
		"Bool":               tt.Bool,
		"Bytes":              tt.Bytes,
		"String":             tt.String,
		"Float":              tt.Float,
		"Int":                tt.Int,
		"UInt":               tt.UInt,
		"Section":            tt.Section,
		"SectionKeys":        tt.SectionKeys,
		"SectionExists":      tt.SectionExists,
		"SectionDelete":      tt.SectionDelete,
		"Exists":             tt.Exists,
		"SectionKeyNotClash": tt.SectionKeyNotClash,
	} {
		t.Run(name, func(t *testing.T) {
			test(t)
		})
	}
}

func (tt Impl) Keys(t *testing.T) {
	t.Run("added keys are returned", func(t *testing.T) {
		s := tt.New()
		defer tt.Done(s)

		s.Set("a", "one")
		s.Set("b", "two")

		keys := s.Keys()
		assert.Len(t, keys, 2)
		assert.Contains(t, keys, "a")
		assert.Contains(t, keys, "b")
	})
}

func (tt Impl) Type(t *testing.T) {
	t.Run("types return correct values", func(t *testing.T) {
		s := tt.New()
		defer tt.Done(s)

		s.Set("int", 1)
		s.Set("uint", uint(1))
		s.Set("string", "Hello World")
		s.Set("float", 1.0)
		s.Set("bool", true)
		s.Set("bytes", []byte("data"))

		s2 := tt.Switch(s)
		defer tt.Done(s2)

		assert.Equal(t, persistence.Int, s2.Type("int"))
		assert.Equal(t, persistence.UnsignedInt, s2.Type("uint"))
		assert.Equal(t, persistence.String, s2.Type("string"))
		assert.Equal(t, persistence.Float, s2.Type("float"))
		assert.Equal(t, persistence.Bool, s2.Type("bool"))
		assert.Equal(t, persistence.Bytes, s2.Type("bytes"))
		assert.Equal(t, persistence.None, s2.Type("missing"))
	})
}

func (tt Impl) Delete(t *testing.T) {
	t.Run("deleting a key removes it", func(t *testing.T) {
		s := tt.New()
		defer tt.Done(s)

		s.Set("a", "one")

		assert.Contains(t, s.Keys(), "a")

		assert.True(t, s.Delete("a"))

		assert.NotContains(t, s.Keys(), "a")
	})

	t.Run("returns false if key not present", func(t *testing.T) {
		s := tt.New()
		defer tt.Done(s)

		assert.False(t, s.Delete("a"))
	})
}

func (tt Impl) Bool(t *testing.T) {
	t.Run("can be set and retrieved", func(t *testing.T) {
		s := tt.New()
		defer tt.Done(s)

		val, found := s.Bool("boolKey")
		assert.False(t, val)
		assert.False(t, found)

		val, found = s.Bool("boolKey", true)
		assert.True(t, val)
		assert.False(t, found)

		s.Set("boolKey", true)

		s2 := tt.Switch(s)
		defer tt.Done(s2)

		val, found = s2.Bool("boolKey", true)
		assert.True(t, val)
		assert.True(t, found)
	})
}

func (tt Impl) Bytes(t *testing.T) {
	t.Run("can be set and retrieved", func(t *testing.T) {
		s := tt.New()
		defer tt.Done(s)

		val, found := s.Bytes("bytesKey")
		assert.Nil(t, val)
		assert.False(t, found)

		val, found = s.Bytes("bytesKey", []byte{})
		assert.Equal(t, []byte{}, val)
		assert.False(t, found)

		s.Set("bytesKey", []byte{0x01})

		s2 := tt.Switch(s)
		defer tt.Done(s2)

		val, found = s2.Bytes("bytesKey", nil)
		assert.Equal(t, []byte{0x01}, val)
		assert.True(t, found)
	})
}

func (tt Impl) String(t *testing.T) {
	t.Run("can be set and retrieved", func(t *testing.T) {
		s := tt.New()
		defer tt.Done(s)

		val, found := s.String("stringKey")
		assert.Equal(t, "", val)
		assert.False(t, found)

		val, found = s.String("stringKey", "none")
		assert.Equal(t, "none", val)
		assert.False(t, found)

		s.Set("stringKey", "test")

		s2 := tt.Switch(s)
		defer tt.Done(s2)

		val, found = s2.String("stringKey", "other")
		assert.Equal(t, "test", val)
		assert.True(t, found)
	})
}

func (tt Impl) Float(t *testing.T) {
	t.Run("can be set and retrieved", func(t *testing.T) {
		s := tt.New()
		defer tt.Done(s)

		val, found := s.Float("float64Key")
		assert.Equal(t, 0.0, val)
		assert.False(t, found)

		val, found = s.Float("float64Key", 0.1)
		assert.Equal(t, 0.1, val)
		assert.False(t, found)

		s.Set("float64Key", 0.2)

		val, found = s.Float("float64Key", 0.1)
		assert.Equal(t, 0.2, val)
		assert.True(t, found)

		s.Set("float32Key", float32(0.2))

		s2 := tt.Switch(s)
		defer tt.Done(s2)

		val, found = s2.Float("float32Key", 0.1)
		assert.InDelta(t, 0.2, val, 0.0001)
		assert.True(t, found)
	})
}

func (tt Impl) Int(t *testing.T) {
	t.Run("can be set and retrieved", func(t *testing.T) {
		s := tt.New()
		defer tt.Done(s)

		val, found := s.Int("intKey")
		assert.Equal(t, int64(0), val)
		assert.False(t, found)

		val, found = s.Int("intKey", 1)
		assert.Equal(t, int64(1), val)
		assert.False(t, found)

		s.Set("intKey", 2)

		val, found = s.Int("intKey", 1)
		assert.Equal(t, int64(2), val)
		assert.True(t, found)

		s.Set("int8Key", int8(2))
		s.Set("int16Key", int16(2))
		s.Set("int32Key", int32(2))
		s.Set("int64Key", int64(2))

		s2 := tt.Switch(s)
		defer tt.Done(s2)

		val, found = s2.Int("int8Key", 1)
		assert.Equal(t, int64(2), val)
		assert.True(t, found)

		val, found = s2.Int("int16Key", 1)
		assert.Equal(t, int64(2), val)
		assert.True(t, found)

		val, found = s2.Int("int32Key", 1)
		assert.Equal(t, int64(2), val)
		assert.True(t, found)

		val, found = s2.Int("int64Key", 1)
		assert.Equal(t, int64(2), val)
		assert.True(t, found)
	})
}

func (tt Impl) UInt(t *testing.T) {
	t.Run("can be set and retrieved", func(t *testing.T) {
		s := tt.New()
		defer tt.Done(s)

		val, found := s.UInt("intKey")
		assert.Equal(t, uint64(0), val)
		assert.False(t, found)

		val, found = s.UInt("intKey", 1)
		assert.Equal(t, uint64(1), val)
		assert.False(t, found)

		s.Set("intKey", uint(2))

		val, found = s.UInt("intKey", 1)
		assert.Equal(t, uint64(2), val)
		assert.True(t, found)

		s.Set("int8Key", uint8(2))
		s.Set("int16Key", uint16(2))
		s.Set("int32Key", uint32(2))
		s.Set("int64Key", uint64(2))

		s2 := tt.Switch(s)
		defer tt.Done(s2)

		val, found = s2.UInt("int8Key", 1)
		assert.Equal(t, uint64(2), val)
		assert.True(t, found)

		val, found = s2.UInt("int16Key", 1)
		assert.Equal(t, uint64(2), val)
		assert.True(t, found)

		val, found = s2.UInt("int32Key", 1)
		assert.Equal(t, uint64(2), val)
		assert.True(t, found)

		val, found = s2.UInt("int64Key", 1)
		assert.Equal(t, uint64(2), val)
		assert.True(t, found)
	})
}

func (tt Impl) Section(t *testing.T) {
	t.Run("a chained section can be created and persists upon retrieval", func(t *testing.T) {
		s := tt.New()
		defer tt.Done(s)

		cs := s.Section("tier1", "tier2")
		cs.Set("key", "value")

		s2 := tt.Switch(s)
		defer tt.Done(s2)

		assert.Contains(t, s2.SectionKeys(), "tier1")

		t1 := s2.Section("tier1")

		assert.Contains(t, t1.SectionKeys(), "tier2")

		t2 := t1.Section("tier2")

		v, _ := t2.String("key")
		assert.Equal(t, "value", v)
	})
}

func (tt Impl) SectionKeys(t *testing.T) {
	t.Run("seconds can be listed", func(t *testing.T) {
		s := tt.New()
		defer tt.Done(s)

		s.Section("one")
		s.Section("two")

		s2 := tt.Switch(s)
		defer tt.Done(s2)

		assert.Contains(t, s2.SectionKeys(), "one")
		assert.Contains(t, s2.SectionKeys(), "two")
	})
}

func (tt Impl) SectionDelete(t *testing.T) {
	t.Run("seconds can be deleted", func(t *testing.T) {
		s := tt.New()
		defer tt.Done(s)

		s.Section("one")
		assert.Contains(t, s.SectionKeys(), "one")

		s.SectionDelete("one")

		s2 := tt.Switch(s)
		defer tt.Done(s2)

		assert.NotContains(t, s2.SectionKeys(), "one")
	})
}

func (tt Impl) Exists(t *testing.T) {
	t.Run("returns if a key exists", func(t *testing.T) {
		s := tt.New()
		defer tt.Done(s)

		s.Set("key", "value")

		s2 := tt.Switch(s)
		defer tt.Done(s2)

		assert.True(t, s2.Exists("key"))
		assert.False(t, s2.Exists("otherKey"))
	})
}

func (tt Impl) SectionExists(t *testing.T) {
	t.Run("returns if a section exists", func(t *testing.T) {
		s := tt.New()
		defer tt.Done(s)

		_ = s.Section("key")

		s2 := tt.Switch(s)
		defer tt.Done(s2)

		assert.True(t, s2.SectionExists("key"))
		assert.False(t, s2.SectionExists("otherKey"))
	})
}

func (tt Impl) SectionKeyNotClash(t *testing.T) {
	t.Run("ensure that keys and sections dont shared the same name space", func(t *testing.T) {
		s := tt.New()
		defer tt.Done(s)

		s.Section("key")
		s.Section("key2")
		s.Set("key", 42)
		s.Set("key3", 42)

		actualKeyInt, _ := s.Int("key")
		assert.Equal(t, int64(42), actualKeyInt)

		s2 := tt.Switch(s)
		defer tt.Done(s2)

		assert.Contains(t, s2.Keys(), "key")
		assert.NotContains(t, s2.Keys(), "key2")
		assert.Contains(t, s2.Keys(), "key3")

		assert.Contains(t, s2.SectionKeys(), "key")
		assert.Contains(t, s2.SectionKeys(), "key2")
		assert.NotContains(t, s2.SectionKeys(), "key3")
	})
}

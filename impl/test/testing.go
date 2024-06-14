package test

import (
	"github.com/shimmeringbee/persistence"
	"github.com/stretchr/testify/assert"
	"testing"
)

type Impl struct {
	New func() persistence.Section
}

func (tt Impl) Test(t *testing.T) {
	for name, test := range map[string]func(*testing.T){
		"Keys":               tt.Keys,
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

		_ = s.Set("a", "one")
		_ = s.Set("b", "two")

		keys := s.Keys()
		assert.Len(t, keys, 2)
		assert.Contains(t, keys, "a")
		assert.Contains(t, keys, "b")
	})
}

func (tt Impl) Delete(t *testing.T) {
	t.Run("deleting a key removes it", func(t *testing.T) {
		s := tt.New()

		_ = s.Set("a", "one")

		assert.Contains(t, s.Keys(), "a")

		assert.True(t, s.Delete("a"))

		assert.NotContains(t, s.Keys(), "a")
	})

	t.Run("returns false if key not present", func(t *testing.T) {
		s := tt.New()

		assert.False(t, s.Delete("a"))
	})
}

func (tt Impl) Bool(t *testing.T) {
	t.Run("can be set and retrieved", func(t *testing.T) {
		s := tt.New()

		val, found := s.Bool("boolKey")
		assert.False(t, val)
		assert.False(t, found)

		val, found = s.Bool("boolKey", true)
		assert.True(t, val)
		assert.False(t, found)

		assert.NoError(t, s.Set("boolKey", true))

		val, found = s.Bool("boolKey", true)
		assert.True(t, val)
		assert.True(t, found)
	})
}

func (tt Impl) Bytes(t *testing.T) {
	t.Run("can be set and retrieved", func(t *testing.T) {
		s := tt.New()

		val, found := s.Bytes("bytesKey")
		assert.Nil(t, val)
		assert.False(t, found)

		val, found = s.Bytes("bytesKey", []byte{})
		assert.Equal(t, []byte{}, val)
		assert.False(t, found)

		assert.NoError(t, s.Set("bytesKey", []byte{0x01}))

		val, found = s.Bytes("bytesKey", nil)
		assert.Equal(t, []byte{0x01}, val)
		assert.True(t, found)
	})
}

func (tt Impl) String(t *testing.T) {
	t.Run("can be set and retrieved", func(t *testing.T) {
		s := tt.New()

		val, found := s.String("stringKey")
		assert.Equal(t, "", val)
		assert.False(t, found)

		val, found = s.String("stringKey", "none")
		assert.Equal(t, "none", val)
		assert.False(t, found)

		assert.NoError(t, s.Set("stringKey", "test"))

		val, found = s.String("stringKey", "other")
		assert.Equal(t, "test", val)
		assert.True(t, found)
	})
}

func (tt Impl) Float(t *testing.T) {
	t.Run("can be set and retrieved", func(t *testing.T) {
		s := tt.New()

		val, found := s.Float("float64Key")
		assert.Equal(t, 0.0, val)
		assert.False(t, found)

		val, found = s.Float("float64Key", 0.1)
		assert.Equal(t, 0.1, val)
		assert.False(t, found)

		assert.NoError(t, s.Set("float64Key", 0.2))

		val, found = s.Float("float64Key", 0.1)
		assert.Equal(t, 0.2, val)
		assert.True(t, found)

		assert.NoError(t, s.Set("float32Key", float32(0.2)))

		val, found = s.Float("float32Key", 0.1)
		assert.InDelta(t, 0.2, val, 0.0001)
		assert.True(t, found)
	})
}

func (tt Impl) Int(t *testing.T) {
	t.Run("can be set and retrieved", func(t *testing.T) {
		s := tt.New()

		val, found := s.Int("intKey")
		assert.Equal(t, int64(0), val)
		assert.False(t, found)

		val, found = s.Int("intKey", 1)
		assert.Equal(t, int64(1), val)
		assert.False(t, found)

		assert.NoError(t, s.Set("intKey", 2))

		val, found = s.Int("intKey", 1)
		assert.Equal(t, int64(2), val)
		assert.True(t, found)

		assert.NoError(t, s.Set("int8Key", int8(2)))

		val, found = s.Int("int8Key", 1)
		assert.Equal(t, int64(2), val)
		assert.True(t, found)

		assert.NoError(t, s.Set("int16Key", int16(2)))

		val, found = s.Int("int16Key", 1)
		assert.Equal(t, int64(2), val)
		assert.True(t, found)

		assert.NoError(t, s.Set("int32Key", int32(2)))

		val, found = s.Int("int32Key", 1)
		assert.Equal(t, int64(2), val)
		assert.True(t, found)

		assert.NoError(t, s.Set("int64Key", int64(2)))

		val, found = s.Int("int64Key", 1)
		assert.Equal(t, int64(2), val)
		assert.True(t, found)
	})
}

func (tt Impl) UInt(t *testing.T) {
	t.Run("can be set and retrieved", func(t *testing.T) {
		s := tt.New()

		val, found := s.UInt("intKey")
		assert.Equal(t, uint64(0), val)
		assert.False(t, found)

		val, found = s.UInt("intKey", 1)
		assert.Equal(t, uint64(1), val)
		assert.False(t, found)

		assert.NoError(t, s.Set("intKey", uint(2)))

		val, found = s.UInt("intKey", 1)
		assert.Equal(t, uint64(2), val)
		assert.True(t, found)

		assert.NoError(t, s.Set("int8Key", uint8(2)))

		val, found = s.UInt("int8Key", 1)
		assert.Equal(t, uint64(2), val)
		assert.True(t, found)

		assert.NoError(t, s.Set("int16Key", uint16(2)))

		val, found = s.UInt("int16Key", 1)
		assert.Equal(t, uint64(2), val)
		assert.True(t, found)

		assert.NoError(t, s.Set("int32Key", uint32(2)))

		val, found = s.UInt("int32Key", 1)
		assert.Equal(t, uint64(2), val)
		assert.True(t, found)

		assert.NoError(t, s.Set("int64Key", uint64(2)))

		val, found = s.UInt("int64Key", 1)
		assert.Equal(t, uint64(2), val)
		assert.True(t, found)
	})
}

func (tt Impl) Section(t *testing.T) {
	t.Run("a chained section can be created and persists upon retrieval", func(t *testing.T) {
		s := tt.New()

		cs := s.Section("tier1", "tier2")
		_ = cs.Set("key", "value")

		assert.Contains(t, s.SectionKeys(), "tier1")

		t1 := s.Section("tier1")

		assert.Contains(t, t1.SectionKeys(), "tier2")

		t2 := t1.Section("tier2")

		v, _ := t2.String("key")
		assert.Equal(t, "value", v)
	})
}

func (tt Impl) SectionKeys(t *testing.T) {
	t.Run("seconds can be listed", func(t *testing.T) {
		s := tt.New()

		s.Section("one")
		s.Section("two")

		assert.Contains(t, s.SectionKeys(), "one")
		assert.Contains(t, s.SectionKeys(), "two")
	})
}

func (tt Impl) SectionDelete(t *testing.T) {
	t.Run("seconds can be deleted", func(t *testing.T) {
		s := tt.New()

		s.Section("one")
		assert.Contains(t, s.SectionKeys(), "one")

		s.SectionDelete("one")
		assert.NotContains(t, s.SectionKeys(), "one")
	})
}

func (tt Impl) Exists(t *testing.T) {
	t.Run("returns if a key exists", func(t *testing.T) {
		s := tt.New()

		_ = s.Set("key", "value")
		assert.True(t, s.Exists("key"))
		assert.False(t, s.Exists("otherKey"))
	})
}

func (tt Impl) SectionExists(t *testing.T) {
	t.Run("returns if a section exists", func(t *testing.T) {
		s := tt.New()

		_ = s.Section("key")
		assert.True(t, s.SectionExists("key"))
		assert.False(t, s.SectionExists("otherKey"))
	})
}

func (tt Impl) SectionKeyNotClash(t *testing.T) {
	t.Run("ensure that keys and sections dont shared the same name space", func(t *testing.T) {
		s := tt.New()

		s.Section("key")
		s.Section("key2")
		s.Set("key", 42)
		s.Set("key3", 42)

		actualKeyInt, _ := s.Int("key")
		assert.Equal(t, int64(42), actualKeyInt)

		assert.Contains(t, s.Keys(), "key")
		assert.NotContains(t, s.Keys(), "key2")
		assert.Contains(t, s.Keys(), "key3")

		assert.Contains(t, s.SectionKeys(), "key")
		assert.Contains(t, s.SectionKeys(), "key2")
		assert.NotContains(t, s.SectionKeys(), "key3")
	})
}
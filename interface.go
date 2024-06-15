package persistence

type Section interface {
	Section(key ...string) Section
	SectionKeys() []string
	SectionExists(key string) bool
	SectionDelete(key string) bool

	Keys() []string
	Exists(key string) bool
	Type(key string) ValueType

	Int(key string, defValue ...int64) (int64, bool)
	UInt(key string, defValue ...uint64) (uint64, bool)
	String(key string, defValue ...string) (string, bool)
	Bool(key string, defValue ...bool) (bool, bool)
	Float(key string, defValue ...float64) (float64, bool)
	Bytes(key string, defValue ...[]byte) ([]byte, bool)

	Set(key string, value interface{})

	Delete(key string) bool
}

type ValueType uint8

const (
	Int         ValueType = 0
	UnsignedInt ValueType = 1
	String      ValueType = 2
	Bool        ValueType = 3
	Float       ValueType = 4
	Bytes       ValueType = 5
	None        ValueType = 255
)

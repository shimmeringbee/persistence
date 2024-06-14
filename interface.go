package persistence

type Section interface {
	Section(key ...string) Section
	SectionKeys() []string
	DeleteSection(key string) bool

	Keys() []string
	Exists(key string) bool

	Int(key string, defValue ...int64) (int64, bool)
	UInt(key string, defValue ...uint64) (uint64, bool)
	String(key string, defValue ...string) (string, bool)
	Bool(key string, defValue ...bool) (bool, bool)
	Float(key string, defValue ...float64) (float64, bool)
	Bytes(key string, defValue ...[]byte) ([]byte, bool)

	Set(key string, value interface{}) error

	Delete(key string) bool
}

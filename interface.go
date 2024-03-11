package persistence

type Section interface {
	Keys() []string
	Exists(key string) bool

	Section(key ...string) Section

	Int(key string, defValue ...int) (int, bool)
	UInt(key string, defValue ...uint) (uint, bool)
	String(key string, defValue ...string) (string, bool)
	Bool(key string, defValue ...bool) (bool, bool)
	Float(key string, defValue ...float64) (float64, bool)
	Bytes(key string, defValue ...[]byte) ([]byte, bool)

	Set(key string, value interface{}) error

	Delete(key string)
}

func StoreComplex[T any](section Section, key string, val T, enc func(Section, string, T) error) error {
	return enc(section, key, val)
}

func RetrieveComplex[T any](section Section, key string, dec func(Section, string) (T, bool), defValue ...T) (T, bool) {
	if v, ok := dec(section, key); ok {
		return v, ok
	} else {
		if len(defValue) > 0 {
			return defValue[0], false
		} else {
			v = *new(T)
			return v, false
		}
	}
}

/*

	IEEEAddress(key string, defValue ...zigbee.IEEEAddress) (zigbee.IEEEAddress, bool)
	ClusterID(key string, defValue ...zigbee.ClusterID) (zigbee.ClusterID, bool)
	Endpoint(key string, defValue ...zigbee.Endpoint) (zigbee.Endpoint, bool)
	AttributeID(key string, defValue ...zcl.AttributeID) (zcl.AttributeID, bool)

	As(key string, destValue any) bool
*/

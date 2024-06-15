package converter

import "github.com/shimmeringbee/persistence"

func Store[T any](section persistence.Section, key string, val T, enc func(persistence.Section, string, T)) {
	enc(section, key, val)
}

func Retrieve[T any](section persistence.Section, key string, dec func(persistence.Section, string) (T, bool), defValue ...T) (T, bool) {
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

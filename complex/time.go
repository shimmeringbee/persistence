package complex

import (
	"github.com/shimmeringbee/persistence"
	"time"
)

func TimeEncoder(s persistence.Section, k string, v time.Time) error {
	return s.Set(k, v.UnixMilli())
}

func TimeDecoder(s persistence.Section, k string) (time.Time, bool) {
	if ev, found := s.Int(k); found {
		return time.UnixMilli(ev), true
	} else {
		return time.Time{}, false
	}
}

func DurationEncoder(s persistence.Section, k string, v time.Duration) error {
	return s.Set(k, v.Milliseconds())
}

func DurationDecoder(s persistence.Section, k string) (time.Duration, bool) {
	if ev, found := s.Int(k); found {
		return time.Duration(ev) * time.Millisecond, true
	} else {
		return time.Duration(0), false
	}
}

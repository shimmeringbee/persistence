package converter

import (
	"github.com/shimmeringbee/persistence"
	"github.com/shimmeringbee/zcl"
	"github.com/shimmeringbee/zigbee"
	"strconv"
)

func AttributeIDEncoder(s persistence.Section, k string, v zcl.AttributeID) error {
	return s.Set(k, int64(v))
}

func AttributeIDDecoder(s persistence.Section, k string) (zcl.AttributeID, bool) {
	if ev, found := s.Int(k); found {
		return zcl.AttributeID(ev), true
	} else {
		return 0, false
	}
}

func AttributeDataTypeEncoder(s persistence.Section, k string, v zcl.AttributeDataType) error {
	return s.Set(k, int64(v))
}

func AttributeDataTypeDecoder(s persistence.Section, k string) (zcl.AttributeDataType, bool) {
	if ev, found := s.Int(k); found {
		return zcl.AttributeDataType(ev), true
	} else {
		return 0, false
	}
}

func IEEEEncoder(s persistence.Section, k string, v zigbee.IEEEAddress) error {
	return s.Set(k, v.String())
}

func IEEEDecoder(s persistence.Section, k string) (zigbee.IEEEAddress, bool) {
	if ev, found := s.String(k); found {
		if value, err := strconv.ParseUint(ev, 16, 64); err != nil {
			return zigbee.EmptyIEEEAddress, false
		} else {
			return zigbee.IEEEAddress(value), true
		}
	} else {
		return zigbee.EmptyIEEEAddress, false
	}
}

func ClusterIDEncoder(s persistence.Section, k string, v zigbee.ClusterID) error {
	return s.Set(k, int64(v))
}

func ClusterIDDecoder(s persistence.Section, k string) (zigbee.ClusterID, bool) {
	if ev, found := s.Int(k); found {
		return zigbee.ClusterID(ev), true
	} else {
		return 0, false
	}
}

func EndpointEncoder(s persistence.Section, k string, v zigbee.Endpoint) error {
	return s.Set(k, int64(v))
}

func EndpointDecoder(s persistence.Section, k string) (zigbee.Endpoint, bool) {
	if ev, found := s.Int(k); found {
		return zigbee.Endpoint(ev), true
	} else {
		return 0, false
	}
}

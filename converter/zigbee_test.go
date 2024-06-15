package converter

import (
	"github.com/shimmeringbee/persistence/impl/memory"
	"github.com/shimmeringbee/zcl"
	"github.com/shimmeringbee/zigbee"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestClusterID(t *testing.T) {
	t.Run("stored and retrieved", func(t *testing.T) {
		s := memory.New()

		expected := zigbee.ClusterID(1)

		Store(s, Key, expected, ClusterIDEncoder)

		actual, found := Retrieve(s, Key, ClusterIDDecoder)
		assert.True(t, found)
		assert.Equal(t, expected, actual)
	})
}

func TestEndpoint(t *testing.T) {
	t.Run("stored and retrieved", func(t *testing.T) {
		s := memory.New()

		expected := zigbee.Endpoint(1)

		Store(s, Key, expected, EndpointEncoder)

		actual, found := Retrieve(s, Key, EndpointDecoder)
		assert.True(t, found)
		assert.Equal(t, expected, actual)
	})
}

func TestIEEEAddress(t *testing.T) {
	t.Run("stored and retrieved", func(t *testing.T) {
		s := memory.New()

		expected := zigbee.GenerateLocalAdministeredIEEEAddress()

		Store(s, Key, expected, IEEEEncoder)

		actual, found := Retrieve(s, Key, IEEEDecoder)
		assert.True(t, found)
		assert.Equal(t, expected, actual)
	})
}

func TestNetworkAddress(t *testing.T) {
	t.Run("stored and retrieved", func(t *testing.T) {
		s := memory.New()

		expected := zigbee.NetworkAddress(0x1122)

		Store(s, Key, expected, NetworkAddressEncoder)

		actual, found := Retrieve(s, Key, NetworkAddressDecoder)
		assert.True(t, found)
		assert.Equal(t, expected, actual)
	})
}

func TestLogicalType(t *testing.T) {
	t.Run("stored and retrieved", func(t *testing.T) {
		s := memory.New()

		expected := zigbee.Router

		Store(s, Key, expected, LogicalTypeEncoder)

		actual, found := Retrieve(s, Key, LogicalTypeDecoder)
		assert.True(t, found)
		assert.Equal(t, expected, actual)
	})
}

func TestAttributeID(t *testing.T) {
	t.Run("stored and retrieved", func(t *testing.T) {
		s := memory.New()

		expected := zcl.AttributeID(1)

		Store(s, Key, expected, AttributeIDEncoder)

		actual, found := Retrieve(s, Key, AttributeIDDecoder)
		assert.True(t, found)
		assert.Equal(t, expected, actual)
	})
}

func TestAttributeDataType(t *testing.T) {
	t.Run("stored and retrieved", func(t *testing.T) {
		s := memory.New()

		expected := zcl.AttributeDataType(1)

		Store(s, Key, expected, AttributeDataTypeEncoder)

		actual, found := Retrieve(s, Key, AttributeDataTypeDecoder)
		assert.True(t, found)
		assert.Equal(t, expected, actual)
	})
}

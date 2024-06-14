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

		err := Store(s, Key, expected, ClusterIDEncoder)
		assert.NoError(t, err)

		actual, found := Retrieve(s, Key, ClusterIDDecoder)
		assert.True(t, found)
		assert.Equal(t, expected, actual)
	})
}

func TestEndpoint(t *testing.T) {
	t.Run("stored and retrieved", func(t *testing.T) {
		s := memory.New()

		expected := zigbee.Endpoint(1)

		err := Store(s, Key, expected, EndpointEncoder)
		assert.NoError(t, err)

		actual, found := Retrieve(s, Key, EndpointDecoder)
		assert.True(t, found)
		assert.Equal(t, expected, actual)
	})
}

func TestIEEEAddress(t *testing.T) {
	t.Run("stored and retrieved", func(t *testing.T) {
		s := memory.New()

		expected := zigbee.GenerateLocalAdministeredIEEEAddress()

		err := Store(s, Key, expected, IEEEEncoder)
		assert.NoError(t, err)

		actual, found := Retrieve(s, Key, IEEEDecoder)
		assert.True(t, found)
		assert.Equal(t, expected, actual)
	})
}

func TestNetworkAddress(t *testing.T) {
	t.Run("stored and retrieved", func(t *testing.T) {
		s := memory.New()

		expected := zigbee.NetworkAddress(0x1122)

		err := Store(s, Key, expected, NetworkAddressEncoder)
		assert.NoError(t, err)

		actual, found := Retrieve(s, Key, NetworkAddressDecoder)
		assert.True(t, found)
		assert.Equal(t, expected, actual)
	})
}

func TestLogicalType(t *testing.T) {
	t.Run("stored and retrieved", func(t *testing.T) {
		s := memory.New()

		expected := zigbee.Router

		err := Store(s, Key, expected, LogicalTypeEncoder)
		assert.NoError(t, err)

		actual, found := Retrieve(s, Key, LogicalTypeDecoder)
		assert.True(t, found)
		assert.Equal(t, expected, actual)
	})
}

func TestAttributeID(t *testing.T) {
	t.Run("stored and retrieved", func(t *testing.T) {
		s := memory.New()

		expected := zcl.AttributeID(1)

		err := Store(s, Key, expected, AttributeIDEncoder)
		assert.NoError(t, err)

		actual, found := Retrieve(s, Key, AttributeIDDecoder)
		assert.True(t, found)
		assert.Equal(t, expected, actual)
	})
}

func TestAttributeDataType(t *testing.T) {
	t.Run("stored and retrieved", func(t *testing.T) {
		s := memory.New()

		expected := zcl.AttributeDataType(1)

		err := Store(s, Key, expected, AttributeDataTypeEncoder)
		assert.NoError(t, err)

		actual, found := Retrieve(s, Key, AttributeDataTypeDecoder)
		assert.True(t, found)
		assert.Equal(t, expected, actual)
	})
}

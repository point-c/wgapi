package parser

import (
	"encoding/hex"
	"github.com/point-c/wgapi/internal"
	"github.com/point-c/wgapi/internal/value"
	"github.com/point-c/wgapi/internal/value/wgkey"
	"github.com/stretchr/testify/require"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	"net"
	"testing"
)

func TestParseIPNet(t *testing.T) {
	input := []byte("192.168.1.0/24")
	_, expected, _ := net.ParseCIDR("192.168.1.0/24")

	parsed, err := ParseIPNet[internal.TestKey](input)
	require.NoError(t, err)
	require.Equal(t, value.IPNet[internal.TestKey](*expected), parsed)
	_, err = ParseIPNet[internal.TestKey](nil)
	require.Error(t, err)
}

func TestParseOne(t *testing.T) {
	input := []byte("1")

	parsed, err := ParseOne[internal.TestKey](input)
	require.NoError(t, err)
	require.Equal(t, value.One[internal.TestKey]{}, parsed)
	_, err = ParseOne[internal.TestKey](nil)
	require.Error(t, err)
}

func TestParseTrue(t *testing.T) {
	input := []byte("true")

	parsed, err := ParseTrue[internal.TestKey](input)
	require.NoError(t, err)
	require.Equal(t, value.True[internal.TestKey]{}, parsed)
	_, err = ParseTrue[internal.TestKey](nil)
	require.Error(t, err)
}

func TestParseUint16(t *testing.T) {
	input := []byte("123")
	expected := uint16(123)

	parsed, err := ParseUint16[internal.TestKey](input)
	require.NoError(t, err)
	require.Equal(t, value.Uint16[internal.TestKey](expected), parsed)
}

func TestParseUint32(t *testing.T) {
	input := []byte("12345")
	expected := uint32(12345)

	parsed, err := ParseUint32[internal.TestKey](input)
	require.NoError(t, err)
	require.Equal(t, value.Uint32[internal.TestKey](expected), parsed)
}

func TestParseUint64(t *testing.T) {
	input := []byte("1234567890")
	expected := uint64(1234567890)

	parsed, err := ParseUint64[internal.TestKey](input)
	require.NoError(t, err)
	require.Equal(t, value.Uint64[internal.TestKey](expected), parsed)
}

func TestParseInt64(t *testing.T) {
	input := []byte("-12345")
	expected := int64(-12345)

	parsed, err := ParseInt64[internal.TestKey](input)
	require.NoError(t, err)
	require.Equal(t, value.Int64[internal.TestKey](expected), parsed)
}

func TestParseKey(t *testing.T) {
	privateKey, _ := wgtypes.GeneratePrivateKey()
	key := value.Key[internal.TestKey, wgkey.Private](privateKey)

	parsed, err := ParseKey[internal.TestKey, wgkey.Private]([]byte(key.String()))
	require.NoError(t, err)
	require.Equal(t, key, parsed)
	_, err = ParseKey[internal.TestKey, wgkey.Private]([]byte("test"))
	require.ErrorAs(t, err, new(hex.InvalidByteError))
	_, err = ParseKey[internal.TestKey, wgkey.Private]([]byte("beef"))
	require.ErrorContains(t, err, "wrong key size")
}

func TestParseUDPAddr(t *testing.T) {
	input := []byte("192.168.1.1:8080")
	expected, _ := net.ResolveUDPAddr("udp", "192.168.1.1:8080")

	parsed, err := ParseUDPAddr[internal.TestKey](input)
	require.NoError(t, err)
	require.Equal(t, value.UDPAddr[internal.TestKey](*expected), parsed)
	_, err = ParseUDPAddr[internal.TestKey]([]byte(".."))
	require.Error(t, err)
}

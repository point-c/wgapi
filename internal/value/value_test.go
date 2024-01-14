package value

import (
	"encoding/hex"
	"github.com/point-c/wgapi/internal"
	"github.com/point-c/wgapi/internal/value/wgkey"
	"github.com/stretchr/testify/require"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	"net"
	"strconv"
	"testing"
)

func TestUint32(t *testing.T) {
	var v Uint32[internal.TestKey] = 123
	require.Equal(t, internal.TestKeyString, v.Key())
	require.Equal(t, strconv.Itoa(int(v)), v.String())
}

func TestUint16(t *testing.T) {
	var v Uint16[internal.TestKey] = 123
	require.Equal(t, internal.TestKeyString, v.Key())
	require.Equal(t, strconv.Itoa(int(v)), v.String())
}

func TestUint64(t *testing.T) {
	var v Uint64[internal.TestKey] = 123
	require.Equal(t, internal.TestKeyString, v.Key())
	require.Equal(t, strconv.Itoa(int(v)), v.String())
}

func TestInt64(t *testing.T) {
	var v Int64[internal.TestKey] = 123
	require.Equal(t, internal.TestKeyString, v.Key())
	require.Equal(t, strconv.Itoa(int(v)), v.String())
}

func TestTrue(t *testing.T) {
	var v True[internal.TestKey]
	require.Equal(t, internal.TestKeyString, v.Key())
	require.Equal(t, "true", v.String())
}

func TestOne(t *testing.T) {
	var v One[internal.TestKey]
	require.Equal(t, internal.TestKeyString, v.Key())
	require.Equal(t, "1", v.String())
}

func TestUDPAddr(t *testing.T) {
	addr := UDPAddr[internal.TestKey]{IP: net.IPv4(192, 168, 0, 1), Port: 8080}
	require.Equal(t, internal.TestKeyString, addr.Key())
	require.Equal(t, "192.168.0.1:8080", addr.String())
}

func TestIPNet(t *testing.T) {
	_, network, _ := net.ParseCIDR("192.168.0.0/24")
	ipnet := IPNet[internal.TestKey](*network)
	require.Equal(t, internal.TestKeyString, ipnet.Key())
	require.Equal(t, "192.168.0.0/24", ipnet.String())
}

func TestKey(t *testing.T) {
	privateKey, _ := wgtypes.GeneratePrivateKey()
	key := Key[internal.TestKey, wgkey.Private](privateKey)

	require.Equal(t, internal.TestKeyString, key.Key())
	require.Equal(t, hex.EncodeToString(privateKey[:]), key.String())

	public, err := key.Public()
	require.NoError(t, err)
	pub := privateKey.PublicKey()
	require.Equal(t, hex.EncodeToString(pub[:]), public.String())
	oldPublic := public
	public, err = public.Public()
	require.NoError(t, err)
	require.Equal(t, oldPublic, public)

	marshaled, err := key.MarshalText()
	require.NoError(t, err)
	require.Equal(t, privateKey.String(), string(marshaled))

	var newKey Key[internal.TestKey, wgkey.Private]
	require.NoError(t, newKey.UnmarshalText(marshaled))
	require.Equal(t, key, newKey)
	require.Error(t, newKey.UnmarshalText([]byte{}))

	var pre Key[internal.TestKey, wgkey.PreShared]
	_, err = pre.Public()
	require.Error(t, err)
}

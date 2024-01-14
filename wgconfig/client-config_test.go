package wgconfig

import (
	"github.com/point-c/wgapi"
	"github.com/stretchr/testify/require"
	"io"
	"net"
	"testing"
)

func TestClientWGConfig(t *testing.T) {
	private, public, err := wgapi.NewPrivatePublic()
	require.NoError(t, err)
	preshared, err := wgapi.NewPreshared()
	require.NoError(t, err)
	client := &Client{
		Private:             private,
		Public:              public,
		PreShared:           preshared,
		Endpoint:            net.UDPAddr{IP: net.IP{192, 168, 1, 1}, Port: 51820},
		PersistentKeepalive: new(uint16),
		AllowedIPs:          []net.IPNet{{IP: net.IP{10, 0, 0, 1}, Mask: net.CIDRMask(32, 32)}},
	}

	*client.PersistentKeepalive = 25

	confReader := client.WGConfig()
	require.NotNil(t, confReader)

	b, err := io.ReadAll(confReader)
	require.NoError(t, err)
	require.NotEmpty(t, b)
}

func TestClientAllowAllIPs(t *testing.T) {
	client := &Client{}

	client.AllowAllIPs()
	require.Len(t, client.AllowedIPs, 1)
	require.Equal(t, net.IPNet(wgapi.EmptySubnet), client.AllowedIPs[0])
}

func TestClientDefaultPersistentKeepAlive(t *testing.T) {
	client := &Client{}

	client.DefaultPersistentKeepAlive()
	require.NotNil(t, client.PersistentKeepalive)
	require.Equal(t, uint16(wgapi.DefaultPersistentKeepalive), *client.PersistentKeepalive)
}

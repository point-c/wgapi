package wgconfig

import (
	"github.com/point-c/wgapi"
	"github.com/stretchr/testify/require"
	"io"
	"net"
	"testing"
)

func TestServerWGConfig(t *testing.T) {
	private, public, err := wgapi.NewPrivatePublic()
	require.NoError(t, err)
	preshared, err := wgapi.NewPreshared()
	require.NoError(t, err)
	server := &Server{
		Private:    private,
		ListenPort: 51820,
		Peers:      []*Peer{},
	}

	server.AddPeer(public, preshared, net.IP{192, 168, 1, 2})
	server.AddPeer(public, preshared, net.IP{192, 168, 1, 3})

	confReader := server.WGConfig()
	require.NotNil(t, confReader)

	b, err := io.ReadAll(confReader)
	require.NoError(t, err)
	require.NotEmpty(t, b)
}

func TestServerDefaultListenPort(t *testing.T) {
	server := &Server{}

	server.DefaultListenPort()
	require.Equal(t, uint16(wgapi.DefaultListenPort), server.ListenPort)
}

func TestServerAddPeer(t *testing.T) {
	_, public, err := wgapi.NewPrivatePublic()
	require.NoError(t, err)
	preshared, err := wgapi.NewPreshared()
	require.NoError(t, err)
	server := &Server{Peers: []*Peer{}}
	ip := net.IP{192, 168, 1, 2}

	server.AddPeer(public, preshared, ip)

	require.Len(t, server.Peers, 1)
	require.Equal(t, public, server.Peers[0].Public)
	require.Equal(t, preshared, server.Peers[0].PreShared)
	require.Equal(t, []net.IPNet{net.IPNet(wgapi.IdentitySubnet(ip))}, server.Peers[0].AllowedIPs)
}

func TestPeerWGConfig(t *testing.T) {
	_, public, err := wgapi.NewPrivatePublic()
	require.NoError(t, err)
	preshared, err := wgapi.NewPreshared()
	require.NoError(t, err)
	peer := &Peer{
		Public:     public,
		PreShared:  preshared,
		AllowedIPs: []net.IPNet{net.IPNet(wgapi.IdentitySubnet(net.IP{192, 168, 1, 2}))},
	}

	confReader := peer.WGConfig()
	require.NotNil(t, confReader)

	b, err := io.ReadAll(confReader)
	require.NoError(t, err)
	require.NotEmpty(t, b)
}

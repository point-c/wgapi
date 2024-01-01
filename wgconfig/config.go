// Package wgconfig has examples for simple client and server configuration.
package wgconfig

import (
	"github.com/point-c/wgapi"
	"golang.zx2c4.com/wireguard/ipc"
	"net"
	"sync"
)

// DefaultPersistentKeepalive is the default persistent keepalive interval.
const DefaultPersistentKeepalive wgapi.PersistentKeepalive = 25

// DefaultListenPort is the default wireguard server port.
const DefaultListenPort wgapi.ListenPort = 51820

// EmptySubnet is the 0.0.0.0/0 subnet
var EmptySubnet = sync.OnceValue(func() wgapi.AllowedIP {
	const allSubnets = "0.0.0.0/0"
	_, ip, _ := net.ParseCIDR(allSubnets)
	return wgapi.AllowedIP(*ip)
})()

const (
	// ErrnoNone means no error occurred
	ErrnoNone = wgapi.Errno(0)
	// ErrnoIO references [ipc.IpcErrorIO].
	ErrnoIO = wgapi.Errno(ipc.IpcErrorIO)
	// ErrnoProtocol references [ipc.IpcErrorProtocol].
	ErrnoProtocol = wgapi.Errno(ipc.IpcErrorProtocol)
	// ErrnoInvalid references [ipc.IpcErrorInvalid].
	ErrnoInvalid = wgapi.Errno(ipc.IpcErrorInvalid)
	// ErrnoPortInUse references [ipc.IpcErrorPortInUse].
	ErrnoPortInUse = wgapi.Errno(ipc.IpcErrorPortInUse)
	// ErrnoUnknown references [ipc.IpcErrorUnknown].
	ErrnoUnknown = wgapi.Errno(int64(ipc.IpcErrorUnknown))
)

// GenerateConfigPair generates a client and server config that are set up to connect to each other.
func GenerateConfigPair(endpoint *net.UDPAddr, clientIP net.IP) (client *Client, server *Server, _ error) {
	clientPrivate, clientPublic, err := wgapi.NewPrivatePublic()
	if err != nil {
		return nil, nil, err
	}

	serverPrivate, serverPublic, err := wgapi.NewPrivatePublic()
	if err != nil {
		return nil, nil, err
	}

	sharedKey, err := wgapi.NewPreshared()
	if err != nil {
		return nil, nil, err
	}

	client = &Client{
		Private:   clientPrivate,
		Public:    serverPublic,
		PreShared: sharedKey,
		Endpoint:  *endpoint,
	}
	client.AllowAllIPs()
	client.DefaultPersistentKeepAlive()

	server = &Server{
		Private:    serverPrivate,
		ListenPort: uint16(endpoint.Port),
	}
	server.AddPeer(clientPublic, sharedKey, clientIP)
	return
}

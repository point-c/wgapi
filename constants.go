package wgapi

import (
	"golang.zx2c4.com/wireguard/ipc"
	"net"
	"sync"
)

// IdentitySubnet converts an IP (v4 or v6) to ipv6/128.
func IdentitySubnet(ip net.IP) AllowedIP {
	return AllowedIP{
		IP:   ip.To16(),
		Mask: net.CIDRMask(128, 128),
	}
}

// DefaultPersistentKeepalive is the default persistent keepalive interval.
const DefaultPersistentKeepalive PersistentKeepalive = 25

// DefaultListenPort is the default wireguard server port.
const DefaultListenPort ListenPort = 51820

// EmptySubnet is the 0.0.0.0/0 subnet
var EmptySubnet = sync.OnceValue(func() AllowedIP {
	const allSubnets = "0.0.0.0/0"
	_, ip, _ := net.ParseCIDR(allSubnets)
	return AllowedIP(*ip)
})()

const (
	// ErrnoNone means no error occurred
	ErrnoNone = Errno(0)
	// ErrnoIO references [ipc.IpcErrorIO].
	ErrnoIO = Errno(ipc.IpcErrorIO)
	// ErrnoProtocol references [ipc.IpcErrorProtocol].
	ErrnoProtocol = Errno(ipc.IpcErrorProtocol)
	// ErrnoInvalid references [ipc.IpcErrorInvalid].
	ErrnoInvalid = Errno(ipc.IpcErrorInvalid)
	// ErrnoPortInUse references [ipc.IpcErrorPortInUse].
	ErrnoPortInUse = Errno(ipc.IpcErrorPortInUse)
	// ErrnoUnknown references [ipc.IpcErrorUnknown].
	ErrnoUnknown = Errno(int64(ipc.IpcErrorUnknown))
)

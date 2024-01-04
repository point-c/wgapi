package wgapi

import (
	"golang.zx2c4.com/wireguard/ipc"
	"net"
)

// IdentitySubnet takes an IP address (either IPv4 or IPv6) and returns it as an IPv6 address with a
// subnet mask of /128. This effectively identifies a single address, as a /128 mask specifies all
// 128 bits of the IPv6 address, leaving no room for a range of addresses.
func IdentitySubnet(ip net.IP) AllowedIP {
	return AllowedIP{
		IP:   ip.To16(),
		Mask: net.CIDRMask(128, 128),
	}
}

const DefaultPersistentKeepalive PersistentKeepalive = 25

const DefaultListenPort ListenPort = 51820

// EmptySubnet represents a subnet configuration that permits all IP addresses.
// It's defined by setting the allowed IP range to '0.0.0.0/0', which effectively means
// there are no restrictions on the IP addresses allowed through this subnet. It's typically
// used in VPN configurations to indicate that all traffic should be routed through the VPN.
var EmptySubnet = parseAllowedIP("0.0.0.0/0")

func parseAllowedIP(s string) AllowedIP {
	_, ip, _ := net.ParseCIDR(s)
	return AllowedIP(*ip)
}

const (
	ErrnoNone      = Errno(0)
	ErrnoIO        = Errno(ipc.IpcErrorIO)
	ErrnoProtocol  = Errno(ipc.IpcErrorProtocol)
	ErrnoInvalid   = Errno(ipc.IpcErrorInvalid)
	ErrnoPortInUse = Errno(ipc.IpcErrorPortInUse)
	ErrnoUnknown   = Errno(int64(ipc.IpcErrorUnknown))
)

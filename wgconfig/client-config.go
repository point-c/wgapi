package wgconfig

import (
	"github.com/point-c/wgapi"
	"io"
	"net"
)

// Client represents a WireGuard peer with necessary configuration parameters.
// It contains cryptographic keys, network settings, and other options to set up a WireGuard connection.
type Client struct {
	Private             wgapi.PrivateKey   // Private is the private key of the peer. (required)
	Public              wgapi.PublicKey    // Public is the public key of the server that the peer intends to connect to. (required)
	PreShared           wgapi.PresharedKey // PreShared is a key shared between the peer and server to further secure the connection. (required)
	Endpoint            net.UDPAddr        // Endpoint specifies the server's address and port as a UDP address. (required)
	PersistentKeepalive *uint16            // PersistentKeepalive is an optional field that specifies the frequency in seconds of keepalive messages sent to maintain the connection. (0 to disable, nil is ignored)
	AllowedIPs          []net.IPNet        // AllowedIPs specifies IP ranges that are allowed to be routed through this WireGuard tunnel.
}

// WGConfig assembles and returns a WireGuard configuration for the client.
func (cfg *Client) WGConfig() io.Reader {
	conf := wgapi.IPC{
		cfg.Private,
		wgapi.ReplacePeers{},
		cfg.Public,
		wgapi.Endpoint(cfg.Endpoint),
		cfg.PreShared,
	}

	if cfg.PersistentKeepalive != nil {
		conf = append(conf, wgapi.PersistentKeepalive(*cfg.PersistentKeepalive))
	}

	for _, allowed := range cfg.AllowedIPs {
		conf = append(conf, wgapi.AllowedIP(allowed))
	}
	return conf.WGConfig()
}

// AllowAllIPs clears [Client].AllowedIPs and sets it to [wgapi.EmptySubnet].
func (cfg *Client) AllowAllIPs() { cfg.AllowedIPs = []net.IPNet{net.IPNet(wgapi.EmptySubnet)} }

// DefaultPersistentKeepAlive sets [Client].PersistentKeepalive to [wgapi.DefaultPersistentKeepalive].
func (cfg *Client) DefaultPersistentKeepAlive() {
	cfg.PersistentKeepalive = new(uint16)
	*cfg.PersistentKeepalive = uint16(wgapi.DefaultPersistentKeepalive)
}

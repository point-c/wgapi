package wgconfig

import (
	"github.com/point-c/wgapi"
	"io"
	"net"
)

// Client is a basic peer configuration.
type Client struct {
	Private             wgapi.PrivateKey   // Private is the peer's private key
	Public              wgapi.PublicKey    // Public is the server's public key
	PreShared           wgapi.PresharedKey // PreShared is the key shared between the peer and server (required)
	Endpoint            net.UDPAddr        // Endpoint is the address and port of the wireguard server
	PersistentKeepalive *uint16            // PersistentKeepalive is the interval (0 to disable, nil is ignored)
	AllowedIPs          []net.IPNet        // AllowedIPs are the addresses allowed to communicate in the tunnel.
}

func (cfg *Client) WGConfig() io.Reader {
	conf := wgapi.IPC{
		cfg.Private,
		wgapi.ReplacePeers{},
		cfg.Public,
		wgapi.Endpoint(cfg.Endpoint),
		cfg.PreShared,
	}.WGConfig()

	if cfg.PersistentKeepalive != nil {
		conf = io.MultiReader(conf, wgapi.IPC{wgapi.PersistentKeepalive(*cfg.PersistentKeepalive)}.WGConfig())
	}

	for _, allowed := range cfg.AllowedIPs {
		conf = io.MultiReader(conf, wgapi.IPC{
			wgapi.AllowedIP(allowed),
		}.WGConfig())
	}
	return conf
}

// AllowAllIPs clears [Client.AllowedIPs] and sets it to [EmptySubnet].
func (cfg *Client) AllowAllIPs() { cfg.AllowedIPs = []net.IPNet{net.IPNet(EmptySubnet)} }

// DefaultPersistentKeepAlive sets [Client.PersistentKeepalive] to [DefaultPersistentKeepalive].
func (cfg *Client) DefaultPersistentKeepAlive() {
	cfg.PersistentKeepalive = new(uint16)
	*cfg.PersistentKeepalive = uint16(DefaultPersistentKeepalive)
}

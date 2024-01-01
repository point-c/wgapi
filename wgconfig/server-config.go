package wgconfig

import (
	"github.com/point-c/wgapi"
	"io"
	"net"
)

type (
	// Server is a basic server configuration.
	Server struct {
		Private    wgapi.PrivateKey // Private is the server's private key.
		ListenPort uint16           // ListenPort is the system port wireguard will listen on.
		Peers      []*Peer          // Peers are the peers allowed to connect to this server.
	}
	// Peer is a client that can connect to the server.
	Peer struct {
		Public     wgapi.PublicKey    // Public is the public key of this peer.
		PreShared  wgapi.PresharedKey // PreShared is the preshared key for this peer (required).
		AllowedIPs []net.IPNet        // AllowedIPs are the addresses allowed to communicate in the tunnel.
	}
)

func (cfg *Server) WGConfig() io.Reader {
	conf := wgapi.IPC{
		cfg.Private,
		wgapi.ListenPort(cfg.ListenPort),
	}.WGConfig()

	for _, peer := range cfg.Peers {
		conf = io.MultiReader(conf, peer.WGConfig())
	}

	return conf
}

// DefaultListenPort sets [Server.DefaultListenPort] to [DefaultListenPort].
func (cfg *Server) DefaultListenPort() {
	cfg.ListenPort = uint16(wgapi.DefaultListenPort)
}

// AddPeer adds a peer with the given public and preshared keys. AllowedIPs is set to the [IdentitySubnet] of the given ip.
func (cfg *Server) AddPeer(publicKey wgapi.PublicKey, preShared wgapi.PresharedKey, ip net.IP) {
	cfg.Peers = append(cfg.Peers, &Peer{
		Public:     publicKey,
		PreShared:  preShared,
		AllowedIPs: []net.IPNet{net.IPNet(wgapi.IdentitySubnet(ip))},
	})
}

func (cfg *Peer) WGConfig() io.Reader {
	return append(wgapi.IPC{
		cfg.Public,
		cfg.PreShared,
	}, func() (ipc wgapi.IPC) {
		for _, ip := range cfg.AllowedIPs {
			ipc = append(ipc, wgapi.AllowedIP(ip))
		}
		return
	}()...).WGConfig()
}

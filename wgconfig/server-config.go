package wgconfig

import (
	"github.com/point-c/wgapi"
	"io"
	"net"
)

type (
	// Server represents a WireGuard server configuration.
	// It holds details necessary to set up and manage a WireGuard server including its private key, the port it listens on, and the set of peers that can connect to it.
	Server struct {
		Private    wgapi.PrivateKey // Private is the server's private key used for establishing secure connections.
		ListenPort uint16           // ListenPort specifies the port number on which the WireGuard server listens for incoming connections.
		Peers      []*Peer          // Peers is a list of clients (peers) that are allowed to establish a connection with this server.
	}
	// Peer represents a client configuration that can establish a connection with a WireGuard server.
	// It includes the client's public key, a pre-shared key for additional security, and a list of allowed IP addresses for routing traffic through the tunnel.
	Peer struct {
		Public     wgapi.PublicKey    // Public is the peer's public key used to identify and authenticate it in the network. (required)
		PreShared  wgapi.PresharedKey // PreShared is an optional additional preshared key for enhancing the security of the peer connection. (required)
		AllowedIPs []net.IPNet        // AllowedIPs are the IP addresses that this peer is allowed to send and receive traffic from in the VPN tunnel. (required)
	}
)

// WGConfig generates and returns a WireGuard configuration for the server.
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

// DefaultListenPort sets [Server].ListenPort to [wgapi.DefaultListenPort].
func (cfg *Server) DefaultListenPort() {
	cfg.ListenPort = uint16(wgapi.DefaultListenPort)
}

// AddPeer creates and adds a new peer to the server's peer list.
// It takes the peer's public and preshared keys, along with an IP address to define its identity subnet for allowed IPs.
// This method is used to dynamically expand the server's network with additional clients.
func (cfg *Server) AddPeer(publicKey wgapi.PublicKey, preShared wgapi.PresharedKey, ip net.IP) {
	cfg.Peers = append(cfg.Peers, &Peer{
		Public:     publicKey,
		PreShared:  preShared,
		AllowedIPs: []net.IPNet{net.IPNet(wgapi.IdentitySubnet(ip))},
	})
}

// WGConfig generates and returns a WireGuard configuration for the peer.
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

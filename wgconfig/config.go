// Package wgconfig has examples for simple client and server configuration.
package wgconfig

import (
	"github.com/point-c/wgapi"
	"net"
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

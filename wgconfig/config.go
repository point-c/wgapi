// Package wgconfig provides examples of simple client and server configurations using the WireGuard API.
// It includes functionality to generate and manage basic WireGuard VPN configurations.
package wgconfig

import (
	"github.com/point-c/wgapi"
	"net"
)

// GenerateConfigPair creates a paired client and server configuration that are
// pre-configured to connect to each other. It generates private/public key pairs
// for both client and server, a shared pre-shared key, and initializes the client
// and server structs with these and additional passed parameters.
//
// The function requires an endpoint (UDP address of the server) and clientIP (IP address
// assigned to the client for VPN use). It returns either a configured client and server struct
// or the error encountered during the configuration process. The client and server
// are set to communicate with each other, with the client configured to allow all IPs and
// have a default persistent keepalive interval.
//
// Parameters:
// - endpoint: A [net.UDPAddr] representing the server's address and port.
// - clientIP: A [net.IP] representing the IP address to be assigned to the client in the VPN.
//
// Returns:
// - client: A pointer to an initialized [Client].
// - server: A pointer to an initialized [Server].
// - error: An error if the generation of keys fail, otherwise nil.
func GenerateConfigPair(endpoint *net.UDPAddr, clientIP net.IP) (client *Client, server *Server, _ error) {
	// Generate private/public key pair for the client.
	clientPrivate, clientPublic, err := wgapi.NewPrivatePublic()
	if err != nil {
		return nil, nil, err
	}

	// Generate private/public key pair for the server.
	serverPrivate, serverPublic, err := wgapi.NewPrivatePublic()
	if err != nil {
		return nil, nil, err
	}

	// Generate a new pre-shared key for additional security between client and server.
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

	// Initialize server configuration with the listening port derived from the endpoint.
	server = &Server{
		Private:    serverPrivate,
		ListenPort: uint16(endpoint.Port),
	}
	server.AddPeer(clientPublic, sharedKey, clientIP)

	return
}

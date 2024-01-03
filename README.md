# WGAPI - Wireguard API Helper

## Introduction
WGAPI is a Go package designed to assist in communicating with the Wireguard userspace module. It focuses on programmatically creating, reading, and handling Wireguard configuration in a text-based format, making it easier to manage Wireguard instances programmatically.

The `wgconfig` package also provides convenient structures and functions for setting up Wireguard server and client configurations. It simplifies the process of generating configurations, managing keys, and setting up network parameters.

## Features
- Parsing and handling Wireguard configuration data.
- Generating new private, public, and pre-shared keys.
- Efficient and structured API communication with Wireguard.

## Usage

### Generating Keys

- **Generate a new private and public key pair:**
  - Errors from `New<key>` functions can only happen if `rand.Reader` causes an error.

  ```go
  private, public, err := wgapi.NewPrivatePublic()
  if err != nil {
      panic(err)
  }
  fmt.Printf("Private Key: %s\nPublic Key: %s\n", private, public)
  ```

- **Generate private key:**
  ```go
  private, err := wgapi.NewPrivate()
  if err != nil {
      panic(err)
  }
  ```

- **Generate preshared key:**
  ```go
  preshared, err := wgapi.NewPreshared()
  if err != nil {
      panic(err)
  }
  ```
  
- **Convert private key to public key:**
  - An error will only happen if you try to get the public key from a preshared key

  ```go
  public, err := private.Public()
  if err != nil {
      panic(err)
  }
  ```
  
### Parsing Wireguard Configuration

```go
// Assuming configData is the Wireguard configuration data you have
configData := "PrivateKey=aSdFgH123456...\nListenPort=51820\n"

// Create a new IPCGet instance and write config data to it
var get wgapi.IPCGet
get.Write([]byte(configData))

// Parse and retrieve the configuration as key-value pairs
ipc, err := get.Value()
if err != nil {
    panic(err)
}

// Handle the parsed IPC data
var tx, rx uint64
for _, v := range ipc {
	// Type assertions can be used to find specific keys in the response
	tx, ok := v.(wgapi.TXBytes)
	if ok {
        tx += uint64(tx)
    }
}
```

### Creating a Wireguard Configuration

#### Freeform

`wgapi.IPC` is a slice type. It can be configured programmatically as you see fit.

```go
r := wgapi.IPC{
	// config values...
}.WGConfig()
// r will be an io.Reader that contains a syntactically valid wireguard configuration
```

#### Config Helpers

##### Server

```go
// Server is a basic server configuration.
type Server struct {
    Private    wgapi.PrivateKey // Private is the server's private key.
    ListenPort uint16           // ListenPort is the system port wireguard will listen on.
    Peers      []*Peer          // Peers are the peers allowed to connect to this server.
}
```

Wireguard's default listening port is `51820`. This can be set with:

```go
// DefaultListenPort sets [Server.DefaultListenPort] to [DefaultListenPort].
func (cfg *Server) DefaultListenPort()
```

Peers can be easily added:

```go
peerIP := net.IPv4(192, 168, 99, 2)

preShared, err := wgapi.NewPreshared()
if err != nil {
    panic(err)
}

privateKey, err := wgapi.NewPrivate()
if err != nil {
    panic(err)
}
// Use private key in client config
publicKey, err := privateKey.Public()

serverCfg.AddPeer(publicKey, preShared, peerIP)
```

##### Client

```go
// Client is a basic peer configuration.
type Client struct {
	Private             wgapi.PrivateKey   // Private is the peer's private key
	Public              wgapi.PublicKey    // Public is the server's public key
	PreShared           wgapi.PresharedKey // PreShared is the key shared between the peer and server (required)
	Endpoint            net.UDPAddr        // Endpoint is the address and port of the wireguard server
	PersistentKeepalive *uint16            // PersistentKeepalive is the interval (0 to disable, nil is ignored)
	AllowedIPs          []net.IPNet        // AllowedIPs are the addresses allowed to communicate in the tunnel.
}
```

Clients will likely want to allow all IPs. This can easily be done with the provided method:

```go
// AllowAllIPs clears [Client.AllowedIPs] and sets it to [EmptySubnet].
func (cfg *Client) AllowAllIPs()
```

The default persistent keepalive can be set with:

```go
// DefaultPersistentKeepAlive sets [Client.PersistentKeepalive] to [DefaultPersistentKeepalive].
func (cfg *Client) DefaultPersistentKeepAlive()
```

##### Pair Generator

A generator exists to create configs quickly.

```go
endpoint, err := net.ResolveUDPAddr("udp", "endpoint address")
if err != nil {
    panic(err)
}

clientIP := net.IPv4(192.168.99.2)

client, server, err := GenerateConfigPair(endpoint, clientIP)
if err != nil {
    panic(err)
}
// `client` and `server` contain valid config to connect to one another.
// Matching keys are generated for each and defaults are filled in.
```
# wgapi

[![Go Reference](https://img.shields.io/badge/godoc-reference-%23007d9c.svg)](https://point-c.github.io/wgapi)

## Introduction
WGAPI is a Go package designed to assist in communicating with the Wireguard userspace module. It focuses on programmatically creating, reading, and handling Wireguard configuration in a text-based format, making it easier to manage Wireguard instances programmatically.

The `wgconfig` package also provides convenient structures and functions for setting up Wireguard server and client configurations. It simplifies the process of generating configurations, managing keys, and setting up network parameters.

## Features
- Parsing and handling Wireguard configuration data.
- Generating new private, public, and pre-shared keys.
- Efficient and structured API communication with Wireguard.
  
## Installation

To use wgevents in your Go project, install it using `go get`:

```bash
go get github.com/point-c/wgapi
```

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
type Server struct {
    Private    wgapi.PrivateKey
    ListenPort uint16      
    Peers      []*Peer     
}
```

Wireguard's default listening port is `51820`. This can be set with:

```go
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
type Client struct {
	Private             wgapi.PrivateKey  
	Public              wgapi.PublicKey   
	PreShared           wgapi.PresharedKey
	Endpoint            net.UDPAddr       
	PersistentKeepalive *uint16          
	AllowedIPs          []net.IPNet      
}
```

Clients will likely want to allow all IPs. This can easily be done with the provided method:

```go
func (cfg *Client) AllowAllIPs()
```

The default persistent keepalive can be set with:

```go
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

## Testing

The package includes tests that demonstrate its functionality. Use Go's testing tools to run the tests:

```bash
go test ./...
```

## Godocs

To regenerate godocs:

```bash
go generate -tags docs ./...
```


## Code Generation

To regenerate generated packages:

```bash
go generate ./...
```

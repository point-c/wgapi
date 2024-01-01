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

```go
r := wgapi.IPC{
	// config values...
}.WGConfig()
// r will be an io.Reader that contains a syntactically valid wireguard configuration
```
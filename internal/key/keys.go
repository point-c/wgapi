package key

type PrivateKey struct{}

func (PrivateKey) String() string { return "private_key" }

type PublicKey struct{}

func (PublicKey) String() string { return "public_key" }

type PresharedKey struct{}

func (PresharedKey) String() string { return "preshared_key" }

type PersistentKeepalive struct{}

func (PersistentKeepalive) String() string { return "persistent_keepalive_interval" }

type ReplacePeers struct{}

func (ReplacePeers) String() string { return "replace_peers" }

type Remove struct{}

func (Remove) String() string { return "remove" }

type UpdateOnly struct{}

func (UpdateOnly) String() string { return "update_only" }

type FWMark struct{}

func (FWMark) String() string { return "fwmark" }

type ProtocolVersion struct{}

func (ProtocolVersion) String() string { return "protocol_version" }

type ReplaceAllowedIPs struct{}

func (ReplaceAllowedIPs) String() string { return "replace_allowed_ips" }

type ListenPort struct{}

func (ListenPort) String() string { return "listen_port" }

type Errno struct{}

func (Errno) String() string { return "errno" }

type LastHandshakeTimeSec struct{}

func (LastHandshakeTimeSec) String() string { return "last_handshake_time_sec" }

type LastHandshakeTimeNSec struct{}

func (LastHandshakeTimeNSec) String() string { return "last_handshake_time_nsec" }

type RXBytes struct{}

func (RXBytes) String() string { return "rx_bytes" }

type TXBytes struct{}

func (TXBytes) String() string { return "tx_bytes" }

type Endpoint struct{}

func (Endpoint) String() string { return "endpoint" }

type AllowedIP struct{}

func (AllowedIP) String() string { return "allowed_ip" }

type Get struct{}

func (Get) String() string { return "get" }

type Set struct{}

func (Set) String() string { return "set" }

package wgapi

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/point-c/wgapi/internal/key"
	"github.com/point-c/wgapi/internal/parser"
	"github.com/point-c/wgapi/internal/value/wgkey"
)

// IPCGet is used to help get information from a wireguard userspace configuration as documented in [wireguard cross-platform documentation].
//
// [wireguard cross-platform documentation]: https://www.wireguard.com/xplatform/
type IPCGet struct{ buf bytes.Buffer }

// Write is used to write 'key=value\n' lines from the wireguard IPC.
// This is used by the wireguard-go's IPC library to write values to this parser.
func (get *IPCGet) Write(b []byte) (int, error) { return get.buf.Write(b) }

// Reset allows this to be reused for another operation.
// Without calling this [IPCGet.Value] will only return the data from the first time this was used.
func (get *IPCGet) Reset() { get.buf.Reset() }

var kvCutChar = []byte{'='}

// Value converts the text data into an [IPC] containing typed values.
func (get *IPCGet) Value() (IPC, error) {
	sc := bufio.NewScanner(&get.buf)
	sc.Split(parser.ScanLines)
	var ipc IPC
	for i := 0; sc.Scan(); i++ {
		line := sc.Bytes()

		k, v, ok := bytes.Cut(line, kvCutChar)
		if !ok {
			return nil, fmt.Errorf("line %d malformed, expected key=value, got %q", i, line)
		}

		p, ok := parsers[string(k)]
		if !ok {
			return nil, fmt.Errorf("line %d malformed, key %q is not valid", i, k)
		}

		kv, err := p(v)
		if err != nil {
			return nil, fmt.Errorf("failed to parse %q on line %d, %w", k, i, err)
		}
		ipc = append(ipc, kv)
	}
	return ipc, sc.Err()
}

var parsers = map[string]func([]byte) (IPCKeyValue, error){
	key.Endpoint{}.String():              parser.ParseUDPAddr[key.Endpoint],
	key.AllowedIP{}.String():             parser.ParseIPNet[key.AllowedIP],
	key.ReplacePeers{}.String():          parser.ParseTrue[key.ReplacePeers],
	key.Remove{}.String():                parser.ParseTrue[key.Remove],
	key.UpdateOnly{}.String():            parser.ParseTrue[key.UpdateOnly],
	key.ReplaceAllowedIPs{}.String():     parser.ParseTrue[key.ReplaceAllowedIPs],
	key.Get{}.String():                   parser.ParseOne[key.Get],
	key.Set{}.String():                   parser.ParseOne[key.Set],
	key.ProtocolVersion{}.String():       parser.ParseOne[key.ProtocolVersion],
	key.RXBytes{}.String():               parser.ParseUint64[key.RXBytes],
	key.TXBytes{}.String():               parser.ParseUint64[key.TXBytes],
	key.FWMark{}.String():                parser.ParseUint32[key.FWMark],
	key.ListenPort{}.String():            parser.ParseUint16[key.ListenPort],
	key.PersistentKeepalive{}.String():   parser.ParseUint16[key.PersistentKeepalive],
	key.Errno{}.String():                 parser.ParseInt64[key.Errno],
	key.LastHandshakeTimeSec{}.String():  parser.ParseInt64[key.LastHandshakeTimeSec],
	key.LastHandshakeTimeNSec{}.String(): parser.ParseInt64[key.LastHandshakeTimeNSec],
	key.PrivateKey{}.String():            parser.ParseKey[key.PrivateKey, wgkey.Private],
	key.PublicKey{}.String():             parser.ParseKey[key.PublicKey, wgkey.Public],
	key.PresharedKey{}.String():          parser.ParseKey[key.PresharedKey, wgkey.PreShared],
}

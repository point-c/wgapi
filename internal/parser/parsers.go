// Package parser contains functions that assist with parsing a wireguard-go ipc response.
package parser

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/point-c/wgapi/internal"
	"github.com/point-c/wgapi/internal/value"
	"github.com/point-c/wgapi/internal/value/wgkey"
	"golang.org/x/exp/constraints"
	"net"
	"strconv"
)

type IPCKeyValue = internal.KeyValue

func ParseIPNet[K fmt.Stringer](b []byte) (IPCKeyValue, error) {
	_, addr, err := net.ParseCIDR(string(b))
	if err != nil {
		return nil, err
	}
	return value.IPNet[K](*addr), nil
}

func ParseOne[K fmt.Stringer](b []byte) (IPCKeyValue, error) {
	return ParseConstant("1", b, value.One[K]{})
}

func ParseTrue[K fmt.Stringer](b []byte) (IPCKeyValue, error) {
	return ParseConstant("true", b, value.True[K]{})
}

func ParseConstant(str string, b []byte, v IPCKeyValue) (IPCKeyValue, error) {
	if string(b) == str {
		return v, nil
	}
	return nil, fmt.Errorf("value %q needs to be exactly %q", b, str)
}

func ParseUint16[K fmt.Stringer](b []byte) (IPCKeyValue, error) {
	return ParseUint[value.Uint16[K]](b)
}

func ParseUint32[K fmt.Stringer](b []byte) (IPCKeyValue, error) {
	return ParseUint[value.Uint32[K]](b)
}

func ParseUint64[K fmt.Stringer](b []byte) (IPCKeyValue, error) {
	return ParseUint[value.Uint64[K]](b)
}

func ParseUint[N constraints.Unsigned](b []byte) (N, error) {
	v, err := strconv.ParseUint(string(b), 10, 8*binary.Size(*new(N)))
	return N(v), err
}

func ParseInt64[K fmt.Stringer](b []byte) (IPCKeyValue, error) {
	n, err := ParseInt[value.Int64[K]](b)
	return n, err
}

func ParseInt[N constraints.Signed](b []byte) (N, error) {
	v, err := strconv.ParseInt(string(b), 10, 8*binary.Size(*new(N)))
	return N(v), err
}

// ParseKey parses a key in hex string format.
func ParseKey[K fmt.Stringer, T wgkey.Type](b []byte) (IPCKeyValue, error) {
	var k wgkey.Key[T]
	n, err := hex.Decode(k[:], b)
	if err != nil {
		return nil, err
	} else if n != len(k) {
		return nil, fmt.Errorf("wrong key size, got %d bytes expected %d", n, len(k))
	}
	return value.Key[K, T](k), nil
}

// ParseUDPAddr parses an address using [net.ResolveUDPAddr].
func ParseUDPAddr[K fmt.Stringer](b []byte) (IPCKeyValue, error) {
	addr, err := net.ResolveUDPAddr("udp", string(b))
	if err != nil {
		return nil, err
	}
	return value.UDPAddr[K](*addr), nil
}

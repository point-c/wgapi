package wgapi_test

import (
	"bufio"
	"encoding/hex"
	"errors"
	"github.com/point-c/wgapi"
	"github.com/point-c/wgapi/internal/parser"
	"io"
	"net"
	"reflect"
	"strings"
	"testing"
)

const exampleSet = `set=1
private_key=e84b5a6d2717c1003a13b431570353dbaca9146cf150c5f8575680feba52027a
fwmark=0
listen_port=12912
replace_peers=true
public_key=b85996fecc9c7f1fc6d2572a76eda11d59bcd20be8e543b15ce4bd85a8e75a33
preshared_key=188515093e952f5f22e865cef3012e72f8b5f0b598ac0309d5dacce3b70fcf52
replace_allowed_ips=true
allowed_ip=192.168.4.4/32
endpoint=[abcd:23::33%2]:51820
public_key=58402e695ba1772b1cc9309755f043251ea77fdcf10fbe63989ceb7e19321376
replace_allowed_ips=true
allowed_ip=192.168.4.6/32
persistent_keepalive_interval=111
endpoint=182.122.22.19:3233
public_key=662e14fd594556f522604703340351258903b64f35553763f19426ab2a515c58
endpoint=5.152.198.39:51820
replace_allowed_ips=true
allowed_ip=192.168.4.10/32
allowed_ip=192.168.4.11/32
public_key=e818b58db5274087fcc1be5dc728cf53d3b5726b4cef6b9bab8f8f8c2452c25c
remove=true
`

// TestIpcRequest_Set validates a set request against the example at https://www.wireguard.com/xplatform/
func TestIpcRequest_Set(t *testing.T) {
	r := wgapi.IPC{
		wgapi.Set{},
		mustParseKey[wgapi.PrivateKey](t, "e84b5a6d2717c1003a13b431570353dbaca9146cf150c5f8575680feba52027a"),
		wgapi.FWMark(0),
		wgapi.ListenPort(12912),
		wgapi.ReplacePeers{},
		mustParseKey[wgapi.PublicKey](t, "b85996fecc9c7f1fc6d2572a76eda11d59bcd20be8e543b15ce4bd85a8e75a33"),
		mustParseKey[wgapi.PresharedKey](t, "188515093e952f5f22e865cef3012e72f8b5f0b598ac0309d5dacce3b70fcf52"),
		wgapi.ReplaceAllowedIPs{},
		wgapi.IdentitySubnet(net.IPv4(192, 168, 4, 4)),
		wgapi.Endpoint{IP: net.ParseIP("abcd:23::33"), Port: 51820, Zone: "2"},
		mustParseKey[wgapi.PublicKey](t, "58402e695ba1772b1cc9309755f043251ea77fdcf10fbe63989ceb7e19321376"),
		wgapi.ReplaceAllowedIPs{},
		wgapi.IdentitySubnet(net.IPv4(192, 168, 4, 6)),
		wgapi.PersistentKeepalive(111),
		wgapi.Endpoint{IP: net.IPv4(182, 122, 22, 19), Port: 3233},
		mustParseKey[wgapi.PublicKey](t, "662e14fd594556f522604703340351258903b64f35553763f19426ab2a515c58"),
		wgapi.Endpoint{IP: net.IPv4(5, 152, 198, 39), Port: 51820},
		wgapi.ReplaceAllowedIPs{},
		wgapi.IdentitySubnet(net.IPv4(192, 168, 4, 10)),
		wgapi.IdentitySubnet(net.IPv4(192, 168, 4, 11)),
		mustParseKey[wgapi.PublicKey](t, "e818b58db5274087fcc1be5dc728cf53d3b5726b4cef6b9bab8f8f8c2452c25c"),
		wgapi.Remove{},
	}.WGConfig()

	sc1 := bufio.NewScanner(strings.NewReader(exampleSet))
	sc1.Split(parser.ScanLines)
	sc2 := bufio.NewScanner(r)
	sc2.Split(parser.ScanLines)
	for i := 0; sc1.Scan() && sc2.Scan(); i++ {
		if sc1.Text() != sc2.Text() {
			t.Fatalf("line %d invalid, %q != %q", i, sc1.Text(), sc2.Text())
		}
	}

	if err := errors.Join(sc1.Err(), sc2.Err()); err != nil {
		t.Fatal(err)
	}
}

func mustParseKey[K wgapi.PublicKey | wgapi.PresharedKey | wgapi.PrivateKey](t *testing.T, key string) K {
	t.Helper()
	var k K
	n, err := hex.Decode(k[:], []byte(key))
	if n != len(k) {
		t.Fatalf("invalid key length %d", n)
	} else if err != nil {
		t.Fatal(err)
	}
	return k
}

const exampleGet = `private_key=e84b5a6d2717c1003a13b431570353dbaca9146cf150c5f8575680feba52027a
listen_port=12912
public_key=b85996fecc9c7f1fc6d2572a76eda11d59bcd20be8e543b15ce4bd85a8e75a33
preshared_key=188515093e952f5f22e865cef3012e72f8b5f0b598ac0309d5dacce3b70fcf52
allowed_ip=192.168.4.4/32
endpoint=[abcd:23::33%2]:51820
public_key=58402e695ba1772b1cc9309755f043251ea77fdcf10fbe63989ceb7e19321376
tx_bytes=38333
rx_bytes=2224
allowed_ip=192.168.4.6/32
persistent_keepalive_interval=111
endpoint=182.122.22.19:3233
public_key=662e14fd594556f522604703340351258903b64f35553763f19426ab2a515c58
endpoint=5.152.198.39:51820
allowed_ip=192.168.4.10/32
allowed_ip=192.168.4.11/32
tx_bytes=1212111
rx_bytes=1929999999
protocol_version=1
errno=0
`

func TestIpcRequest_Get(t *testing.T) {
	var p wgapi.IPCGet
	defer p.Reset()
	if _, err := io.WriteString(&p, exampleGet); err != nil {
		t.Fatal(err)
	}

	v, err := p.Value()
	if err != nil {
		t.Fatal(err)
	}

	comp := wgapi.IPC{
		mustParseKey[wgapi.PrivateKey](t, "e84b5a6d2717c1003a13b431570353dbaca9146cf150c5f8575680feba52027a"),
		wgapi.ListenPort(12912),
		mustParseKey[wgapi.PublicKey](t, "b85996fecc9c7f1fc6d2572a76eda11d59bcd20be8e543b15ce4bd85a8e75a33"),
		mustParseKey[wgapi.PresharedKey](t, "188515093e952f5f22e865cef3012e72f8b5f0b598ac0309d5dacce3b70fcf52"),
		wgapi.IdentitySubnet(net.IPv4(192, 168, 4, 4)),
		wgapi.Endpoint{IP: net.ParseIP("abcd:23::33"), Port: 51820, Zone: "2"},
		mustParseKey[wgapi.PublicKey](t, "58402e695ba1772b1cc9309755f043251ea77fdcf10fbe63989ceb7e19321376"),
		wgapi.TXBytes(38333),
		wgapi.RXBytes(2224),
		wgapi.IdentitySubnet(net.IPv4(192, 168, 4, 6)),
		wgapi.PersistentKeepalive(111),
		wgapi.Endpoint{IP: net.IPv4(182, 122, 22, 19), Port: 3233},
		mustParseKey[wgapi.PublicKey](t, "662e14fd594556f522604703340351258903b64f35553763f19426ab2a515c58"),
		wgapi.Endpoint{IP: net.IPv4(5, 152, 198, 39), Port: 51820},
		wgapi.IdentitySubnet(net.IPv4(192, 168, 4, 10)),
		wgapi.IdentitySubnet(net.IPv4(192, 168, 4, 11)),
		wgapi.TXBytes(1212111),
		wgapi.RXBytes(1929999999),
		wgapi.ProtocolVersion{},
		wgapi.ErrnoNone,
	}

	if len(v) != len(comp) {
		t.Fatalf("expected size %d got %d", len(comp), len(v))
	}

	for i, e := range comp {
		v := v[i]
		if e.Key() != v.Key() {
			t.Fatalf("line %d invalid, got key %q, expected %q", i, v.Key(), e.Key())
		} else if v.String() != v.String() {
			t.Fatalf("line %d invalid, got value %q, expected %q", i, v.String(), e.String())
		}

		if !reflect.TypeOf(e).AssignableTo(reflect.TypeOf(v)) {
			t.Fatalf("line %d invalid, got type %T, expected %T", i, v, e)
		}
	}
}

func TestIpcResquest_NoNewline(t *testing.T) {
	var p wgapi.IPCGet
	defer p.Reset()
	if _, err := io.WriteString(&p, `protocol_version=1`); err != nil {
		t.Fatal(err)
	}

	_, err := p.Value()
	if err == nil {
		t.Fail()
	}
}

func TestAssignability(t *testing.T) {
	tx := wgapi.TXBytes(16)
	var val wgapi.IPCKeyValue = tx
	rx, ok := val.(wgapi.RXBytes)
	if ok {
		t.Fail()
	} else if int(rx) == int(tx) {
		t.Fail()
	}
}

package wgapi

import (
	"github.com/stretchr/testify/require"
	"io"
	"testing"
)

func TestIPCGet_Value(t *testing.T) {
	for _, tt := range []struct {
		name  string
		input string
	}{
		{
			name:  "invalid key",
			input: "foo=bar\n",
		},
		{
			name:  "invalid line",
			input: "foo\n",
		},
		{
			name:  "invalid value",
			input: "fwmark=aa\n",
		},
		{
			name:  "invalid key",
			input: "public_key=zz\n",
		},
		{
			name:  "invalid key length",
			input: "public_key=ab\n",
		},
		{
			name:  "invalid true",
			input: "remove=false\n",
		},
		{
			name:  "invalid ipnet",
			input: "allowed_ip=..\n",
		},
		{
			name:  "invalid udp addr",
			input: "endpoint=..\n",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var g IPCGet
			_, err := io.WriteString(&g, tt.input)
			require.NoError(t, err)
			_, err = g.Value()
			require.Error(t, err)
		})
	}
}

package wgconfig

import (
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"io"
	"net"
	"testing"
)

func TestGenerateConfigPair(t *testing.T) {
	randR := rand.Reader
	defer func() { rand.Reader = randR }()
	errExp := errors.New(uuid.NewString())

	for i := range make([]struct{}, 4) {
		t.Run(fmt.Sprintf("err after %d key generates", i), func(t *testing.T) {
			r := newTestCountdownErrReader(t, randR, i, errExp)
			rand.Reader = r
			client, server, err := GenerateConfigPair(&net.UDPAddr{}, net.IPv4allrouter)
			switch i {
			case 0, 1, 2:
				require.ErrorContains(t, err, errExp.Error())
				require.Nil(t, client)
				require.Nil(t, server)
			case 3:
				require.NoError(t, err)
				require.NotNil(t, client)
				require.NotNil(t, server)
			default:
				t.Fatalf("invalid i %d", i)
			}
		})
	}
}

type testCountdownErrReader struct {
	t         testing.TB
	base      io.Reader
	countdown int
	err       error
}

func newTestCountdownErrReader(t testing.TB, base io.Reader, countdown int, err error) *testCountdownErrReader {
	t.Helper()
	return &testCountdownErrReader{
		t:         t,
		base:      base,
		countdown: countdown,
		err:       err,
	}
}

func (t *testCountdownErrReader) Read(b []byte) (int, error) {
	if t.countdown <= 0 {
		return 0, t.err
	}
	t.countdown--
	return t.base.Read(b)
}
